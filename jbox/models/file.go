package models

type DownloadProgressHandler = func(downloaded int64, total int64)

type FileInfo struct {
	PrefixNeid      string `json:"prefix_neid"`
	PID             string `json:"pid"`
	LocalModifyTime string `json:"localModifyTime"`
	Approvable      bool   `json:"approveable"`
	FromName        string `json:"from_name"`
	IsTeam          bool   `json:"is_team"`
	Result          string `json:"result"`
	Path            string `json:"path"`
	Nsid            int64  `json:"nsid"`
	IsDeleted       bool   `json:"is_deleted"`
	ContentSize     int64  `json:"content_size"`
	IsDir           bool   `json:"is_dir"`
	IsShared        bool   `json:"is_shared"`
	Modified        string `json:"modified"`
	CreatorUid      int64  `json:"creator_uid"`
	From            string `json:"from"`
	IsBookmark      bool   `json:"is_bookmark"`
	Neid            string `json:"neid"`
	Creator         string `json:"creator"`
	Offset          int64  `json:"offset"`
	Authable        bool   `json:"authable"`
	SupportPreview  string `json:"support_preview"`
	PathType        string `json:"path_type"`
	AccessMode      int64  `json:"access_mode"`
	IsGroup         bool   `json:"is_group"`
	DirType         int64  `json:"dir_type"`
	DeliveryCode    string `json:"delivery_code"`
	Size            string `json:"size"`
	UpdatorUid      int64  `json:"updator_uid"`
	Bytes           int64  `json:"bytes"`
	Updator         string `json:"updator"`
	ShareToPersonal bool   `json:"share_to_personal"`
	Hash            string `json:"hash"`
	Desc            string `json:"desc"`
	IsDisplay       bool   `json:"is_display"`
	Message         string `json:"message"`
	Type            string `json:"type"`
}

type DirectoryInfo struct {
	FileInfo
	Content []*DirectoryInfo `json:"content,omitempty"`
}

type BatchMoveData struct {
	To        ToData     `json:"to"`
	From      []FromData `json:"from"`
	OtherData string     `json:"other_data"`
}

type FromData struct {
	Root     string `json:"root"`
	Rev      string `json:"rev"`
	Path     string `json:"path"`
	PathType string `json:"path_type"`
	From     string `json:"from"`
	Nsid     int64  `json:"nsid"`
	Neid     string `json:"neid"`
}

type ToData struct {
	Root     string `json:"root"`
	Path     string `json:"path"`
	PathType string `json:"path_type"`
	From     string `json:"from"`
	Neid     string `json:"neid"`
	Nsid     int64  `json:"nsid"`
}

type BatchMoveResult struct {
	Result   string               `json:"result"`
	TaskID   string               `json:"task_id"`
	TaskType int64                `json:"task_type"`
	Success  []BatchMoveFileEntry `json:"success"`
	Failed   []BatchMoveFileEntry `json:"failed"`
}

type BatchMoveFileEntry struct {
	Result  string                 `json:"result"`
	Neid    float64                `json:"neid"`
	Path    string                 `json:"path"`
	Nsid    int64                  `json:"nsid"`
	JSON    BatchMoveFileEntryJSON `json:"json"`
	NewNeid float64                `json:"newNeid"`
	NewPath string                 `json:"newPath"`
}

type BatchMoveFileEntryJSON struct {
	Result string  `json:"result"`
	Neid   float64 `json:"neid"`
	Path   string  `json:"path"`
	Nsid   int64   `json:"nsid"`
}
