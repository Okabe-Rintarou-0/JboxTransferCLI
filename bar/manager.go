package bar

import (
	"jtrans/db/models"
)

type IManager interface {
	Set64(path string, val int64, total int64)
	Prepare(dbModel *models.FileSyncTask)
	Error(dbModel *models.FileSyncTask)
	PrepareDownloadChunk(curChunkSize, chunkNo, chunkCount int64, dbModel *models.FileSyncTask)
	PrepareUploadChunk(curChunkSize, chunkNo, chunkCount int64, dbModel *models.FileSyncTask)
	Finish(dbModel *models.FileSyncTask)
	Wait()
}
