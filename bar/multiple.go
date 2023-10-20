package bar

import (
	"jtrans/db/models"
	"sync"
	"time"

	"github.com/fatih/color"
	"github.com/vbauerster/mpb/v8"
	"github.com/vbauerster/mpb/v8/decor"
)

type barEntry struct {
	Bar         *mpb.Bar
	DBModel     *models.FileSyncTask
	Downloaded  int64
	Uploaded    int64
	Downloading bool
}

type MultipleBarManager struct {
	p          *mpb.Progress
	barMap     map[string]*barEntry
	numWorkers int
	rwLock     sync.RWMutex
}

func (m *MultipleBarManager) Wait() {
	m.p.Wait()
}

func (m *MultipleBarManager) Error(dbModel *models.FileSyncTask) {
	m.rwLock.Lock()
	entry := m.barMap[dbModel.FilePath]
	entry.Bar.Abort(false)
	delete(m.barMap, dbModel.FilePath)
	m.rwLock.Unlock()
}

func (m *MultipleBarManager) Finish(dbModel *models.FileSyncTask) {
	m.rwLock.Lock()
	if dbModel.IsDir() {
		m.barMap[dbModel.FilePath].Bar.SetCurrent(1)
	}
	delete(m.barMap, dbModel.FilePath)
	m.rwLock.Unlock()
}

func (m *MultipleBarManager) PrepareDownloadChunk(curChunkSize, chunkNo, chunkCount int64, dbModel *models.FileSyncTask) {
	m.rwLock.Lock()
	m.barMap[dbModel.FilePath].Downloading = true
	m.rwLock.Unlock()
}

func (m *MultipleBarManager) PrepareUploadChunk(curChunkSize, chunkNo, chunkCount int64, dbModel *models.FileSyncTask) {
	m.rwLock.Lock()
	m.barMap[dbModel.FilePath].Downloading = false
	m.rwLock.Unlock()
}

func (m *MultipleBarManager) Set64(path string, val int64, total int64) {
	m.rwLock.RLock()
	entry := m.barMap[path]
	m.rwLock.RUnlock()
	progress := (entry.Downloaded + entry.Uploaded + val) / 2
	if val == total {
		if entry.Downloading {
			entry.Downloaded += total
		} else {
			entry.Uploaded += total
		}
	}
	entry.Bar.SetCurrent(progress)
}

func (m *MultipleBarManager) Prepare(dbModel *models.FileSyncTask) {
	red, green, yellow := color.New(color.FgRed), color.New(color.FgGreen), color.New(color.FgYellow)
	m.rwLock.Lock()
	if dbModel.IsFile() {
		completed := dbModel.GetCompletedSize()
		bar := &barEntry{
			Bar: m.p.AddBar(dbModel.Size,
				mpb.PrependDecorators(
					decor.Meta(
						decor.Name(dbModel.FileName, decor.WC{W: len(dbModel.FileName) + 1, C: decor.DidentRight}),
						toMetaFunc(green),
					),
					decor.Meta(
						decor.OnAbort(
							decor.Name("", decor.WC{W: 0}),
							"同步失败！",
						),
						toMetaFunc(yellow),
					),
					decor.Counters(decor.SizeB1024(0), "% .2f / % .2f"),
				),
				mpb.AppendDecorators(
					decor.Percentage(),
				),
			),
			DBModel:    dbModel,
			Downloaded: completed,
			Uploaded:   completed,
		}
		bar.Bar.SetCurrent(completed)
		m.barMap[dbModel.FilePath] = bar
	} else {
		m.barMap[dbModel.FilePath] = &barEntry{
			Bar: m.p.AddBar(1,
				mpb.PrependDecorators(
					decor.Meta(decor.Name(dbModel.FilePath, decor.WC{W: len(dbModel.FilePath) + 1, C: decor.DidentRight}), toMetaFunc(red)),
					decor.CountersNoUnit("%d / %d", decor.WCSyncWidth),
				),
			),
			DBModel: dbModel,
		}
	}
	m.rwLock.Unlock()
}

func toMetaFunc(c *color.Color) func(string) string {
	return func(s string) string {
		return c.Sprint(s)
	}
}

func NewMultipleBarManager(numWorkers int) IManager {
	return &MultipleBarManager{
		p:          mpb.New(mpb.WithRefreshRate(180 * time.Millisecond)),
		barMap:     map[string]*barEntry{},
		numWorkers: numWorkers,
	}
}
