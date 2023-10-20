package worker

import (
	"fmt"
	"jtrans/bar"
	"jtrans/db"
	dbmodels "jtrans/db/models"
	"jtrans/jbox"
	jmodels "jtrans/jbox/models"
	"jtrans/tbox"
	"strings"

	"github.com/spf13/cast"
)

type DirectorySyncWorker struct {
	retryTimes int
	path       string
	total      int64
	succ       int64
	page       int
	progress   int64
	jboxDir    *jmodels.DirectoryInfo
	tcli       tbox.IClient
	jcli       jbox.IClient
	state      SyncTaskState
	bar        bar.IManager
	dbModel    *dbmodels.FileSyncTask
}

func newDirectorySyncWorker(jcli jbox.IClient, tcli tbox.IClient, path string, bar bar.IManager) *DirectorySyncWorker {
	return &DirectorySyncWorker{
		retryTimes: 3,
		path:       path,
		total:      0,
		succ:       0,
		page:       0,
		progress:   0,
		state:      SyncWait,
		jcli:       jcli,
		tcli:       tcli,
		bar:        bar,
	}
}

func NewDirectorySyncWorkerFromDBModel(jcli jbox.IClient, tcli tbox.IClient, model *dbmodels.FileSyncTask, bar bar.IManager) *DirectorySyncWorker {
	w := newDirectorySyncWorker(jcli, tcli, model.FilePath, bar)
	w.dbModel = model
	if len(model.RemainParts) > 0 {
		w.page = cast.ToInt(model.RemainParts)
		w.succ = cast.ToInt64(w.page * 50)
	}
	return w
}

func (w *DirectorySyncWorker) GetName() string {
	parts := strings.Split(w.path, "/")
	name := parts[len(parts)-1]
	if name == "" {
		name = "根目录"
	}
	return name
}

func (w *DirectorySyncWorker) GetPath() string {
	return w.path
}

func (w *DirectorySyncWorker) GetParentPath() string {
	parts := strings.Split(w.path, "/")
	parent := strings.Join(parts[:len(parts)-1], "/")
	return parent
}

func (w *DirectorySyncWorker) Start() error {
	if err := w.internalStart(); err != nil {
		w.bar.Error(w.dbModel)
		return err
	}
	return nil
}

func (w *DirectorySyncWorker) handleError() {
	w.state = SyncError
	w.dbModel.State = dbmodels.Error
	_ = db.Update(w.dbModel)
}

func (w *DirectorySyncWorker) handleComplete() {
	w.state = SyncComplete
	w.dbModel.State = dbmodels.Done
	_ = db.Update(w.dbModel)
}

func (w *DirectorySyncWorker) internalStart() error {
	w.state = SyncRunning
	w.dbModel.State = dbmodels.Busy
	w.bar.Prepare(w.dbModel)
	res, err := w.tcli.CreateDirectory(w.path)
	if err != nil {
		if res == nil {
			w.handleError()
			return err
		}
		if res.Code != "SameNameDirectoryOrFileExists" && res.Code != "RootDirectoryNotAllowed" {
			w.handleError()
			return fmt.Errorf("创建目录失败：%s", res.Message)
		}
	}

	for {
		t := w.retryTimes
		for t > 0 {
			t -= 1
			w.jboxDir, err = w.jcli.GetDirectoryInfo(w.path, w.page)
			if err != nil {
				err = fmt.Errorf("获取目录信息失败：%s", err.Error())
				continue
			}
			w.total = w.jboxDir.ContentSize
			if len(w.jboxDir.Content) == 0 {
				break
			}

			order := db.GetMinOrder() - 1
			tx := db.Begin()
			for _, file := range w.jboxDir.Content {
				tx = tx.Create(dbmodels.FromJBoxFileInfo(&file.FileInfo, order))
			}
			if w.page != cast.ToInt(w.dbModel.RemainParts) {
				w.dbModel.RemainParts = cast.ToString(w.page)
				tx.Save(w.dbModel)
			}
			tx.Commit()

			w.succ += cast.ToInt64(len(w.jboxDir.Content))
			w.page += 1
			break
		}

		if t <= 0 {
			w.handleError()
			return err
		}
		if w.jboxDir == nil || len(w.jboxDir.Content) == 0 {
			break
		}
	}

	w.handleComplete()
	w.bar.Finish(w.dbModel)
	return nil
}
