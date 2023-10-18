package worker

import (
	"jtrans/constants"
	"jtrans/jbox"
	"jtrans/jbox/models"
	"jtrans/utils"
)

type JBoxDownloadWorker struct {
	path          string
	size          int64
	chunkCount    int64
	chunkProgress int64
	cli           jbox.IClient
	onProgress    models.DownloadProgressHandler
}

func NewJBoxDownloadWorker(cli jbox.IClient, path string, size int64, onProgress models.DownloadProgressHandler) *JBoxDownloadWorker {
	return &JBoxDownloadWorker{
		path:          path,
		size:          size,
		chunkCount:    utils.ComputeChunkCount(size, constants.ChunkSize),
		chunkProgress: 0,
		cli:           cli,
		onProgress:    onProgress,
	}
}

func (w *JBoxDownloadWorker) GetChunk(chunkNo int64) ([]byte, error) {
	curChunkSize := utils.ComputeCurrentChunkSize(chunkNo, w.chunkCount, w.size)
	return w.cli.GetChunk(w.path, chunkNo, curChunkSize, func(downloaded int64, total int64) {
		w.chunkProgress = downloaded
		if w.onProgress != nil {
			w.onProgress(downloaded, total)
		}
	})
}

func (w *JBoxDownloadWorker) ClearProgress() {
	w.chunkProgress = 0
}
