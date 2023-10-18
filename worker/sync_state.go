package worker

type SyncTaskState int

const (
	SyncWait SyncTaskState = iota
	SyncRunning
	SyncError
	SyncComplete
	SyncPause
)
