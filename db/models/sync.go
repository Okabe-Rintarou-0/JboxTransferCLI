package models

import (
	"encoding/json"
	"jtrans/constants"
	"jtrans/jbox/models"
	"jtrans/utils"
)

type TaskState = int
type FileType = int

const (
	Idle TaskState = iota
	Busy
	Error
	Done
	Cancel
)

const (
	File TaskState = iota
	Directory
)

type FileSyncTask struct {
	ID    uint `gorm:"column:Id;primaryKey;autoIncrement:true"`
	Order int  `gorm:"column:Order;index"`
	/// <summary>
	/// 0:file
	/// 1:folder
	/// </summary>
	Type       int    `gorm:"column:Type"`
	FileName   string `gorm:"column:FileName"`
	FilePath   string `gorm:"column:FilePath"`
	Size       int64  `gorm:"column:Size"`
	ConfirmKey string `gorm:"column:ConfirmKey"`
	/// <summary>
	/// 0:Idle
	/// 1:Busy
	/// 2:Error
	/// 3:Done
	/// 4:Cancel
	/// </summary>
	State       int    `gorm:"column:State"`
	MD5Part     string `gorm:"column:MD5_Part"`
	MD5Ori      string `gorm:"column:MD5_Ori"`
	CRC64Part   int64  `gorm:"column:CRC64_Part"`
	RemainParts string `gorm:"column:RemainParts"`
}

func (t *FileSyncTask) IsDir() bool {
	return t.Type == Directory
}

func (t *FileSyncTask) IsFile() bool {
	return t.Type == File
}

func (t *FileSyncTask) GetCompletedSize() int64 {
	var parts []int64
	err := json.Unmarshal([]byte(t.RemainParts), &parts)
	if err != nil {
		return 0
	}

	var succCount int64 = 0
	if len(parts) > 0 {
		succCount = parts[0]
	}
	return succCount * constants.ChunkSize
}

func NewFileSyncTask(fileType int, path string, size int64, order int, hash string) *FileSyncTask {
	return &FileSyncTask{
		Order:    order,
		Type:     fileType,
		FileName: utils.GetFileName(path),
		FilePath: path,
		Size:     size,
		State:    Idle,
		MD5Ori:   hash,
	}
}

func FromJBoxFileInfo(file *models.FileInfo, order int) *FileSyncTask {
	if file.IsDir {
		return NewFileSyncTask(Directory, file.Path, 0, order, "")
	}
	return NewFileSyncTask(File, file.Path, file.Bytes, order, file.Hash)
}

func (FileSyncTask) TableName() string {
	return "SyncTaskDbModel"
}
