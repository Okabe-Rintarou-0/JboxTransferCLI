package models

import (
	"strings"
)

type UploadProgressHandler = func(uploaded int64, total int64)

type DirectoryInfo struct {
	Path          []string        `json:"path"`
	SubDirCount   int64           `json:"subDirCount"`
	FileCount     int64           `json:"fileCount"`
	TotalNum      int64           `json:"totalNum"`
	ETag          string          `json:"eTag"`
	Contents      []*FileInfo     `json:"contents"`
	LocalSync     interface{}     `json:"localSync"`
	AuthorityList map[string]bool `json:"authorityList"`
}

type FileInfo struct {
	Path                     []string        `json:"path"`
	Name                     string          `json:"name"`
	Type                     string          `json:"type"`
	UserID                   string          `json:"userId"`
	CreationTime             string          `json:"creationTime"`
	ModificationTime         string          `json:"modificationTime"`
	VersionID                interface{}     `json:"versionId"`
	VirusAuditStatus         int64           `json:"virusAuditStatus"`
	SensitiveWordAuditStatus int64           `json:"sensitiveWordAuditStatus"`
	ContentType              string          `json:"contentType"`
	Size                     string          `json:"size"`
	ETag                     string          `json:"eTag"`
	Crc64                    string          `json:"crc64"`
	MetaData                 MetaData        `json:"metaData"`
	AuthorityList            map[string]bool `json:"authorityList"`
	FileType                 string          `json:"fileType"`
	PreviewByDoc             bool            `json:"previewByDoc"`
	PreviewByCI              bool            `json:"previewByCI"`
	PreviewAsIcon            bool            `json:"previewAsIcon"`
}

func (d *DirectoryInfo) FullPath() string {
	return "/" + strings.Join(d.Path, "/")
}

func (f *FileInfo) FullPath() string {
	return "/" + strings.Join(f.Path, "/")
}

func (f *FileInfo) IsDir() bool {
	return f.Type == "dir"
}

type MetaData map[string]interface{}

type StartChunkUploadResult struct {
	ConfirmKey string                          `json:"confirmKey"`
	Domain     string                          `json:"domain"`
	Path       string                          `json:"path"`
	UploadID   string                          `json:"uploadId"`
	Parts      map[string]StartChunkUploadPart `json:"parts"`
	Expiration string                          `json:"expiration"`
}

type StartChunkUploadPart struct {
	Headers Headers `json:"headers"`
}

type Headers struct {
	XAmzDate          string `json:"x-amz-date"`
	XAmzContentSha256 string `json:"x-amz-content-sha256"`
	Authorization     string `json:"authorization"`
}

type ConfirmChunkUploadResult struct {
	Path                     []string `json:"path"`
	Name                     string   `json:"name"`
	Type                     string   `json:"type"`
	CreationTime             string   `json:"creationTime"`
	ModificationTime         string   `json:"modificationTime"`
	ContentType              string   `json:"contentType"`
	Size                     string   `json:"size"`
	ETag                     string   `json:"eTag"`
	Crc64                    string   `json:"crc64"`
	MetaData                 MetaData `json:"metaData"`
	IsOverwrittened          bool     `json:"isOverwrittened"`
	VirusAuditStatus         int64    `json:"virusAuditStatus"`
	SensitiveWordAuditStatus int64    `json:"sensitiveWordAuditStatus"`
	PreviewByDoc             bool     `json:"previewByDoc"`
	PreviewByCI              bool     `json:"previewByCI"`
	PreviewAsIcon            bool     `json:"previewAsIcon"`
	FileType                 string   `json:"fileType"`
}

type ChunkUploadInfo struct {
	Confirmed                  bool              `json:"confirmed"`
	Path                       []string          `json:"path"`
	Type                       string            `json:"type"`
	CreationTime               string            `json:"creationTime"`
	ConflictResolutionStrategy string            `json:"conflictResolutionStrategy"`
	Force                      bool              `json:"force"`
	UploadID                   string            `json:"uploadId"`
	Parts                      []ChunkUploadPart `json:"parts"`
}

type ChunkUploadPart struct {
	PartNumber   int64  `json:"PartNumber"`
	Size         int64  `json:"Size"`
	ETag         string `json:"ETag"`
	LastModified string `json:"LastModified"`
}

type BatchMoveData struct {
	From                       string `json:"from"`
	To                         string `json:"to"`
	Type                       string `json:"type"`
	ConflictResolutionStrategy string `json:"conflictResolutionStrategy"`
	MoveAuthority              bool   `json:"moveAuthority"`
}

type BatchMoveResult struct {
	Result []BatchMoveResultEntry `json:"result"`
}

type BatchMoveResultEntry struct {
	To            []string `json:"to"`
	From          []string `json:"from"`
	Path          []string `json:"path"`
	Status        int64    `json:"status"`
	MoveAuthority bool     `json:"moveAuthority"`
}
