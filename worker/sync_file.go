package worker

import (
	"encoding/json"
	"fmt"
	"jtrans/constants"
	"jtrans/db"
	dbmodels "jtrans/db/models"
	"jtrans/encrypt"
	"jtrans/jbox"
	"jtrans/tbox"
	"jtrans/tbox/models"
	"jtrans/utils"
	"strings"
	"unicode/utf8"

	"github.com/schollz/progressbar/v3"
	"github.com/spf13/cast"
)

type FileSyncWorker struct {
	retryTimes int
	path       string
	jboxhash   string
	size       int64
	curChunk   *TBoxUploadPartEntry
	chunkCount int64
	succChunk  int64
	state      SyncTaskState
	dbModel    *dbmodels.FileSyncTask
	uploader   *TBoxUploadWorker
	downloader *JBoxDownloadWorker
	crc64      *encrypt.CRC64
	md5        *encrypt.MD5

	bar *progressbar.ProgressBar
}

func newFileSyncWorker(jcli jbox.IClient, path, hash string, size int64, bar *progressbar.ProgressBar) *FileSyncWorker {
	return &FileSyncWorker{
		path:       path,
		jboxhash:   hash,
		size:       size,
		chunkCount: utils.ComputeChunkCount(size, constants.ChunkSize),
		succChunk:  0,
		state:      SyncWait,
		downloader: NewJBoxDownloadWorker(jcli, path, size, func(downloaded int64, total int64) {
			_ = bar.Set64(downloaded)
		}),
		retryTimes: 3,
		bar:        bar,
	}
}

func NewFileSyncWorkerFromDBModel(jcli jbox.IClient, tcli tbox.IClient, model *dbmodels.FileSyncTask, bar *progressbar.ProgressBar) *FileSyncWorker {
	w := newFileSyncWorker(jcli, model.FilePath, model.MD5Ori, model.Size, bar)
	w.dbModel = model
	if model.ConfirmKey == "" {
		w.crc64 = encrypt.NewCRC64()
		w.uploader = NewTBoxUploadWorker(tcli, w.path, w.size, func(uploaded int64, total int64) {
			_ = w.bar.Set64(uploaded)
		})
		w.md5 = encrypt.NewMD5()
	} else {
		var remains []int64
		w.crc64 = encrypt.NewCRC64FromValue(uint64(model.CRC64Part))
		_ = json.Unmarshal([]byte(model.RemainParts), &remains)
		w.succChunk = w.chunkCount - cast.ToInt64(len(remains))
		w.uploader = NewTBoxUploadWorkerForRenewing(tcli, w.path, w.size, model.ConfirmKey, remains, func(uploaded int64, total int64) {
			_ = w.bar.Set64(uploaded)
		})
		storage := encrypt.MD5StateStorage{}
		_ = json.Unmarshal([]byte(model.MD5Part), &storage)
		w.md5 = encrypt.NewMD5FromStorage(&storage)
	}
	return w
}

func (w *FileSyncWorker) GetName() string {
	parts := strings.Split(w.path, "/")
	name := parts[len(parts)-1]
	return name
}

func (w *FileSyncWorker) GetFormatedFileName() string {
	name := w.GetName()
	if utf8.RuneCountInString(name) > 10 {
		return utils.Utf8Substr(name, 0, 10) + "..."
	}
	return name
}

func (w *FileSyncWorker) GetPath() string {
	return w.path
}

func (w *FileSyncWorker) GetParentPath() string {
	parts := strings.Split(w.path, "/")
	parent := strings.Join(parts[:len(parts)-1], "/")
	return parent
}

func (w *FileSyncWorker) Start() error {
	if err := w.internalStart(); err != nil {
		fmt.Println(err.Error())
		return err
	}
	return nil
}

func (w *FileSyncWorker) refreshBar(max int64, desc string) {
	w.bar.Reset()
	w.bar.ChangeMax64(max)
	w.bar.Describe(desc)
}

func (w *FileSyncWorker) handleError() {
	w.state = SyncError
	w.dbModel.State = dbmodels.Error
	_ = db.Update(w.dbModel)
}

func (w *FileSyncWorker) handleComplete() {
	w.state = SyncComplete
	w.dbModel.State = dbmodels.Done
	_ = db.Update(w.dbModel)
}

func (w *FileSyncWorker) updateRemainParts() {
	if w.uploader == nil || w.dbModel == nil {
		return
	}
	remain := make([]int64, len(w.uploader.remainParts))
	for i, part := range w.uploader.remainParts {
		remain[i] = part.PartNumber
	}
	marshalled, _ := json.Marshal(remain)
	w.dbModel.RemainParts = string(marshalled)
}

func (w *FileSyncWorker) internalStart() error {
	var (
		err           error
		chunkData     []byte
		confirmResult *models.ConfirmChunkUploadResult
	)

	w.state = SyncRunning
	fmt.Printf("同步文件 \"%s\"，共有%d个分块，剩余%d个分块\n", w.GetName(), w.chunkCount, w.chunkCount-w.succChunk)
	fmt.Printf("正在准备上传...")
	err = w.uploader.PrepareForUpload()
	if err != nil {
		w.handleError()
		return err
	}
	fmt.Println("完毕")

	w.curChunk, err = w.uploader.GetNextPart()
	if err != nil {
		w.state = SyncError
		w.handleError()
		return err
	}

	w.dbModel.ConfirmKey = w.uploader.confirmKey
	w.updateRemainParts()
	w.dbModel.CRC64Part = int64(w.crc64.GetValue())
	md5Part, _ := json.Marshal(w.md5.GetValue())
	w.dbModel.MD5Part = string(md5Part)
	err = db.Update(w.dbModel)
	if err != nil {
		return err
	}

	for w.curChunk != nil {
		t := w.retryTimes
		w.downloader.ClearProgress()
		w.uploader.ClearProgress()

		chunkNo := w.curChunk.PartNumber
		curChunkSize := utils.ComputeCurrentChunkSize(chunkNo, w.chunkCount, w.size)
		for t > 0 {
			t -= 1

			w.refreshBar(curChunkSize, fmt.Sprintf("[cyan][%d/%d][reset] [red]%s[reset] 下载分块...", chunkNo, w.chunkCount, w.GetFormatedFileName()))
			err = w.uploader.EnsureNoExpire(chunkNo)
			if err != nil {
				continue
			}
			chunkData, err = w.downloader.GetChunk(chunkNo)
			if err != nil {
				err = fmt.Errorf("下载块 %d 发生错误：%s", chunkNo, err.Error())
				continue
			}

			w.refreshBar(curChunkSize, fmt.Sprintf("[cyan][%d/%d][reset] [red]%s[reset] 上传分块...", chunkNo, w.chunkCount, w.GetFormatedFileName()))
			err = w.uploader.Upload(chunkData, chunkNo)
			if err != nil {
				err = fmt.Errorf("上传块 %d 发生错误：%s", chunkNo, err.Error())
				continue
			}
			break
		}
		if t <= 0 {
			w.uploader.ResetPartNumber(w.curChunk)
			w.handleError()
			return err
		}

		if len(chunkData) > 0 {
			sha256 := encrypt.SHA256Hash(chunkData)
			if w.curChunk.PartNumber != 1 {
				sha256 = "," + sha256
			}
			encrypt.MD5HashProc(w.md5, []byte(sha256))
		}
		w.crc64.TransformBlock(chunkData, 0, len(chunkData))
		w.uploader.CompletePart(w.curChunk)
		w.succChunk += 1

		w.dbModel.CRC64Part = int64(w.crc64.GetValue())
		marshalled, _ := json.Marshal(w.md5.GetValue())
		w.dbModel.MD5Part = string(marshalled)
		w.updateRemainParts()
		err = db.Update(w.dbModel)
		if err != nil {
			return err
		}

		w.curChunk, err = w.uploader.GetNextPart()
		if err != nil {
			w.state = SyncError
			return fmt.Errorf("获取下一块发生错误：%s", err.Error())
		}
	}
	fmt.Println("完成！")
	fmt.Printf("正在计算文件校验和...")
	actualHash := encrypt.MD5HashProcFinish(w.md5)
	actualCRC64 := w.crc64.TransformFinalBlock()
	if w.succChunk == w.chunkCount && w.curChunk == nil {
		confirmResult, err = w.uploader.Confirm()
		if err != nil {
			w.handleError()
			return err
		}
		if cast.ToString(actualCRC64) != confirmResult.Crc64 || actualHash != w.jboxhash {
			w.handleError()
			return fmt.Errorf("校验和不匹配！%s %s", actualHash, w.jboxhash)
		}
		w.handleComplete()
		fmt.Println("完毕")
		fmt.Println("同步成功！")
	} else {
		fmt.Println("同步失败！")
	}
	return nil
}
