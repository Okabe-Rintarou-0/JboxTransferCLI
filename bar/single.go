package bar

import (
	"fmt"
	"jtrans/constants"
	"jtrans/db/models"
	"jtrans/utils"
	"unicode/utf8"

	"github.com/schollz/progressbar/v3"
	"github.com/spf13/cast"
)

type SingleBarManager struct {
	bar      *progressbar.ProgressBar
	name     string
	fileType string
}

func (s *SingleBarManager) Wait() {

}

func (s *SingleBarManager) Error(dbModel *models.FileSyncTask) {
	c := "green"
	if dbModel.IsDir() {
		c = "red"
	}
	s.bar.Describe(fmt.Sprintf("[%s]%s[reset] [yellow]下载失败！[reset]", c, s.name))
}

func (s *SingleBarManager) Finish(dbModel *models.FileSyncTask) {
	_ = s.bar.Finish()
	fmt.Println(" 完毕！")
}

func syncFileSingleBar() *progressbar.ProgressBar {
	return progressbar.NewOptions(
		cast.ToInt(constants.ChunkSize),
		progressbar.OptionEnableColorCodes(true),
		progressbar.OptionShowBytes(true),
		progressbar.OptionSetWidth(10),
	)
}
func syncDirSingleBar() *progressbar.ProgressBar {
	return progressbar.NewOptions(
		cast.ToInt(constants.ChunkSize),
		progressbar.OptionEnableColorCodes(true),
		progressbar.OptionShowBytes(false),
		progressbar.OptionSetWidth(10),
	)
}

func (s *SingleBarManager) refreshBar(max int64, desc string) {
	s.bar.Reset()
	s.bar.ChangeMax64(max)
	s.bar.Describe(desc)
}

func (s *SingleBarManager) PrepareDownloadChunk(curChunkSize, chunkNo, chunkCount int64, dbModel *models.FileSyncTask) {
	s.refreshBar(curChunkSize, fmt.Sprintf("[cyan][%d/%d][reset] [green]%s[reset] 下载分块...", chunkNo, chunkCount, s.name))
}

func (s *SingleBarManager) PrepareUploadChunk(curChunkSize, chunkNo, chunkCount int64, dbModel *models.FileSyncTask) {
	s.refreshBar(curChunkSize, fmt.Sprintf("[cyan][%d/%d][reset] [green]%s[reset] 上传分块...", chunkNo, chunkCount, s.name))
}

func (s *SingleBarManager) Set64(path string, val int64, total int64) {
	_ = s.bar.Set64(val)
}

func (s *SingleBarManager) Prepare(dbModel *models.FileSyncTask) {
	s.name = getFormatedFileName(dbModel.FileName)
	if dbModel.IsDir() {
		s.bar = syncDirSingleBar()
		s.fileType = "目录"
		s.refreshBar(1, fmt.Sprintf("[red]%s[reset] 同步中...", s.name))
	} else {
		s.bar = syncFileSingleBar()
		s.fileType = "文件"
	}
}

func getFormatedFileName(name string) string {
	if utf8.RuneCountInString(name) > 13 {
		return utils.Utf8Substr(name, 0, 10) + "..."
	}
	return name
}

func NewSingleBarManager() IManager {
	return &SingleBarManager{
		bar:      syncFileSingleBar(),
		fileType: "文件",
	}
}
