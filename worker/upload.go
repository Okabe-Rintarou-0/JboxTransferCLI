package worker

import (
	"fmt"
	"jtrans/constants"
	"jtrans/tbox"
	tmodels "jtrans/tbox/models"
	"jtrans/utils"
	"time"

	"github.com/spf13/cast"
)

type TBoxUploadState int

const (
	TBoxUploadNotInit TBoxUploadState = iota
	TBoxUploadConfirmKeyInit
	TBoxUploadReady
	TBoxUploading
	TBoxUploadDone
	TBoxUploadError
)

type TBoxUploadPartEntry struct {
	PartNumber int64
	Uploading  bool
}

type TBoxUploadWorker struct {
	ctx           *tmodels.StartChunkUploadResult
	chunkProgress int64
	chunkCount    int64
	size          int64
	path          string
	confirmKey    string
	state         TBoxUploadState
	remainParts   []*TBoxUploadPartEntry
	cli           tbox.IClient
	onProgress    tmodels.UploadProgressHandler
}

func NewTBoxUploadWorker(cli tbox.IClient, path string, size int64, onProgress tmodels.UploadProgressHandler) *TBoxUploadWorker {
	return &TBoxUploadWorker{
		ctx:           nil,
		chunkProgress: 0,
		chunkCount:    utils.ComputeChunkCount(size, constants.ChunkSize),
		size:          size,
		path:          path,
		confirmKey:    "",
		state:         TBoxUploadNotInit,
		remainParts:   nil,
		cli:           cli,
		onProgress:    onProgress,
	}
}

func NewTBoxUploadWorkerForRenewing(cli tbox.IClient, path string, size int64, confirmKey string, remainParts []int64, onProgress tmodels.UploadProgressHandler) *TBoxUploadWorker {
	w := NewTBoxUploadWorker(cli, path, size, onProgress)
	w.confirmKey = confirmKey
	w.state = TBoxUploadConfirmKeyInit
	w.remainParts = make([]*TBoxUploadPartEntry, len(remainParts))
	for i, part := range remainParts {
		w.remainParts[i] = &TBoxUploadPartEntry{
			PartNumber: part,
			Uploading:  false,
		}
	}
	return w
}

func (w *TBoxUploadWorker) ClearProgress() {
	w.chunkProgress = 0
}

func (w *TBoxUploadWorker) CompletePart(part *TBoxUploadPartEntry) {
	for i, entry := range w.remainParts {
		if entry == part {
			w.remainParts = append(w.remainParts[:i], w.remainParts[i+1:]...)
		}
	}
}

func (w *TBoxUploadWorker) GetRefreshPartNumberList() []int64 {
	count := utils.Min(len(w.remainParts), 50)
	entries := w.remainParts[:count]
	var ret []int64
	for _, entry := range entries {
		ret = append(ret, entry.PartNumber)
	}
	return ret
}

func (w *TBoxUploadWorker) PrepareForUpload() (err error) {
	if w.state == TBoxUploadReady || w.state == TBoxUploading || w.state == TBoxUploadDone {
		return nil
	}
	if w.state == TBoxUploadNotInit || w.state == TBoxUploadError {
		w.ctx, err = w.cli.StartChunkUpload(w.path, w.chunkCount)
		if err != nil {
			return err
		}

		w.confirmKey = w.ctx.ConfirmKey
		var i int64
		for i = 1; i <= w.chunkCount; i++ {
			w.remainParts = append(w.remainParts, &TBoxUploadPartEntry{
				PartNumber: i,
				Uploading:  false,
			})
		}
		w.state = TBoxUploadReady
	} else {
		w.ctx, err = w.cli.RenewChunkUpload(w.confirmKey, w.GetRefreshPartNumberList())
		if err != nil {
			return err
		}
		w.state = TBoxUploadReady
	}
	return nil
}

func (w *TBoxUploadWorker) GetNextPart() (*TBoxUploadPartEntry, error) {
	if w.state != TBoxUploadReady && w.state != TBoxUploading {
		return nil, fmt.Errorf("非法状态")
	}
	if len(w.remainParts) == 0 {
		return nil, nil
	}
	for _, part := range w.remainParts {
		if !part.Uploading {
			return part, nil
		}
	}

	return nil, nil
}

func (w *TBoxUploadWorker) EnsureNoExpire(partNumber int64) error {
	if w.ctx.Expiration == "" {
		return nil
	}

	exp, err := time.Parse(time.RFC3339, w.ctx.Expiration)
	if err != nil {
		return err
	}

	if exp.Second()-time.Now().Second() < 30 {
		w.ctx, err = w.cli.RenewChunkUpload(w.confirmKey, w.GetRefreshPartNumberList())
		if err != nil {
			return fmt.Errorf("刷新分块凭据出错: %s", err.Error())
		}
	}
	partKey := cast.ToString(partNumber)
	if _, ok := w.ctx.Parts[partKey]; !ok {
		w.ctx, err = w.cli.RenewChunkUpload(w.confirmKey, w.GetRefreshPartNumberList())
		if err != nil {
			return fmt.Errorf("刷新分块凭据出错: %s", err.Error())
		}
	}
	if _, ok := w.ctx.Parts[partKey]; !ok {
		return fmt.Errorf("已刷新上传凭据，但是未找到块 %d 的信息", partNumber)
	}
	return nil
}

func (w *TBoxUploadWorker) Upload(data []byte, partNumber int64) error {
	return w.cli.Upload(w.ctx, data, partNumber, func(uploaded int64, total int64) {
		w.chunkProgress = uploaded
		if w.onProgress != nil {
			w.onProgress(uploaded, total)
		}
	})
}

func (w *TBoxUploadWorker) Confirm() (*tmodels.ConfirmChunkUploadResult, error) {
	return w.cli.ConfirmChunkUpload(w.confirmKey)
}

func (w *TBoxUploadWorker) ResetPartNumber(part *TBoxUploadPartEntry) {
	part.Uploading = false
}
