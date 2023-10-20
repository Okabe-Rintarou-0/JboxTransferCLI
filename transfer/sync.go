package transfer

import (
	"fmt"
	"jtrans/bar"
	"jtrans/db"
	"jtrans/db/models"
	dbmodels "jtrans/db/models"
	"jtrans/jbox"
	"jtrans/login"
	"jtrans/tbox"
	"jtrans/utils"
	"jtrans/worker"
	"os"
	"sync"

	ignore "github.com/sabhiram/go-gitignore"
	"github.com/spf13/cobra"
)

var (
	syncCmd        *cobra.Command
	syncDirPath    string
	syncFilePath   string
	syncIgnorePath string

	syncAll         bool
	syncRecursively bool
	useBfs          bool
	useDfs          bool

	numWorkers int
	pbar       bar.IManager

	syncErrors []error
	errMtx     = sync.Mutex{}

	workerChan chan worker.IWorker
	wg         sync.WaitGroup
)

func printErrors() {
	fmt.Println("同步过程中发生的错误：")
	for _, err := range syncErrors {
		fmt.Println(err.Error())
	}
}

func getBarManager() bar.IManager {
	if numWorkers > 1 {
		return bar.NewMultipleBarManager(numWorkers)
	}
	return bar.NewSingleBarManager()
}

func syncDirectory(jcli jbox.IClient, tcli tbox.IClient, dirPath string, ignores *ignore.GitIgnore) error {
	dirPath = formatPath(dirPath)
	order := db.GetMinOrder() - 1
	if ignores != nil && ignores.MatchesPath(dirPath) {
		// fmt.Printf("根据\"%s\"配置，目录\"%s\"将被忽略！\n", syncIgnorePath, syncDirPath)
		return nil
	}
	model := db.FindByPath(dirPath)

	if model == nil {
		fileInfo, err := jcli.GetFileInfo(dirPath)
		if err != nil {
			return fmt.Errorf("获取目录信息失败：%s", err.Error())
		}
		model = dbmodels.FromJBoxFileInfo(fileInfo, order)

		err = db.Insert(model)
		if err != nil {
			return fmt.Errorf("插入数据库时出错：%s", err.Error())
		}
	}

	if model.State == dbmodels.Done {
		return fmt.Errorf("指定的目录\"%s\"已同步", dirPath)
	}

	syncWorker := worker.NewDirectorySyncWorkerFromDBModel(jcli, tcli, model, pbar)
	return syncWorker.Start()
}

func syncFile(jcli jbox.IClient, tcli tbox.IClient, filePath string, ignores *ignore.GitIgnore) error {
	filePath = formatPath(filePath)
	order := db.GetMinOrder() - 1
	if ignores != nil && ignores.MatchesPath(filePath) {
		// return fmt.Errorf("根据\"%s\"配置，文件\"%s\"将被忽略！", syncIgnorePath, filePath)
		return nil
	}
	model := db.FindByPath(filePath)
	if model == nil {
		fileInfo, err := jcli.GetFileInfo(filePath)
		if err != nil {
			return fmt.Errorf("获取文件信息失败：%s", err.Error())
		}
		model = dbmodels.FromJBoxFileInfo(fileInfo, order)
	}
	if model.State == dbmodels.Done {
		return fmt.Errorf("指定的文件\"%s\"已同步", filePath)
	}
	syncWorker := worker.NewFileSyncWorkerFromDBModel(jcli, tcli, model, pbar)
	if numWorkers > 1 {
		model.State = models.Busy
		if err := db.Update(model); err != nil {
			return err
		}
		workerChan <- syncWorker
		return nil
	}
	return syncWorker.Start()
}

func syncDirectoryInnerDfs(jcli jbox.IClient, tcli tbox.IClient, dir string, ignores *ignore.GitIgnore) error {
	order := db.GetMinOrder()
	err := syncDirectory(jcli, tcli, dir, ignores)
	if err != nil {
		return err
	}

	for {
		tasks := db.FindIdleTasksWithSmallerOrder(order)
		var filtered []*dbmodels.FileSyncTask
		for _, task := range tasks {
			if ignores != nil && ignores.MatchesPath(task.FilePath) {
				// fmt.Printf("根据\"%s\"配置，文件（目录）\"%s\"将被忽略！\n", syncIgnorePath, task.FilePath)
			} else {
				filtered = append(filtered, task)
			}
		}

		if len(filtered) == 0 {
			break
		}

		for _, task := range filtered {
			if task.IsFile() {
				err = syncFile(jcli, tcli, task.FilePath, ignores)
			} else {
				err = syncDirectoryInnerDfs(jcli, tcli, task.FilePath, ignores)
			}
			if err != nil {
				syncErrors = append(syncErrors, err)
			}
		}
	}
	return nil
}

func syncDirectoryInnerBfs(jcli jbox.IClient, tcli tbox.IClient, dir string, ignores *ignore.GitIgnore) error {
	err := syncDirectory(jcli, tcli, dir, ignores)
	if err != nil {
		return err
	}
	for {
		tasks := db.FindIdleTasks()
		var filtered []*dbmodels.FileSyncTask
		for _, task := range tasks {
			if ignores != nil && ignores.MatchesPath(task.FilePath) {
				// fmt.Printf("根据\"%s\"配置，文件（目录）\"%s\"将被忽略！\n", syncIgnorePath, task.FilePath)
			} else {
				filtered = append(filtered, task)
			}
		}

		if len(filtered) == 0 {
			break
		}

		for _, task := range filtered {
			if task.IsFile() {
				err = syncFile(jcli, tcli, task.FilePath, ignores)
			} else {
				err = syncDirectory(jcli, tcli, task.FilePath, ignores)
			}
			if err != nil {
				syncErrors = append(syncErrors, err)
			}
		}
	}
	return nil
}

func syncDirectoryRecursively(jcli jbox.IClient, tcli tbox.IClient, dir string, ignores *ignore.GitIgnore) error {
	if useBfs {
		return syncDirectoryInnerBfs(jcli, tcli, dir, ignores)
	}
	return syncDirectoryInnerDfs(jcli, tcli, dir, ignores)
}

func syncAllItems(jcli jbox.IClient, tcli tbox.IClient, ignores *ignore.GitIgnore) error {
	return syncDirectoryRecursively(jcli, tcli, "/", ignores)
}

func syncQueueItems(jcli jbox.IClient, tcli tbox.IClient, ignores *ignore.GitIgnore) error {
	var (
		err      error
		filtered []*dbmodels.FileSyncTask
	)
	tasks := db.FindExecutableTasks()

	for _, task := range tasks {
		if ignores != nil && ignores.MatchesPath(task.FilePath) {
			// fmt.Printf("根据\"%s\"配置，文件（目录）\"%s\"将被忽略！\n", syncIgnorePath, task.FilePath)
		} else {
			filtered = append(filtered, task)
		}
	}

	if len(filtered) == 0 {
		fmt.Println("没有可执行的同步任务！")
		return nil
	}

	for _, task := range filtered {
		if task.IsFile() {
			err = syncFile(jcli, tcli, task.FilePath, ignores)
		} else {
			if syncRecursively {
				err = syncDirectoryRecursively(jcli, tcli, task.FilePath, ignores)
			} else {
				err = syncDirectory(jcli, tcli, task.FilePath, ignores)
			}
		}
		if err != nil {
			syncErrors = append(syncErrors, err)
		}
	}

	return nil
}

func parallelWork() {
	workerChan = make(chan worker.IWorker, numWorkers)
	wg = sync.WaitGroup{}
	wg.Add(numWorkers)
	for i := 0; i < numWorkers; i++ {
		go func() {
			var err error
			for worker := range workerChan {
				if err = worker.Start(); err != nil {
					errMtx.Lock()
					syncErrors = append(syncErrors, err)
					errMtx.Unlock()
				}
			}
			wg.Done()
		}()
	}
}

func init() {
	syncCmd = &cobra.Command{
		Use:   "sync",
		Short: "同步文件",
		Long:  "同步文件",
		Run: func(cmd *cobra.Command, args []string) {
			var ignores *ignore.GitIgnore
			jcli, tcli, err := login.GetClient()
			if err != nil {
				fmt.Println(NotLoginHint)
				os.Exit(1)
			}

			if len(syncIgnorePath) > 0 {
				ignores, err = ignore.CompileIgnoreFile(syncIgnorePath)
				if err != nil {
					fmt.Println("指定的 .ignore 文件无效！")
					os.Exit(1)
				}
			}

			numWorkers = utils.Max(numWorkers, 1)
			if numWorkers > 1 {
				parallelWork()
			}
			pbar = getBarManager()

			if syncAll {
				err = syncAllItems(jcli, tcli, ignores)
			} else if len(syncDirPath) > 0 {
				if syncRecursively {
					err = syncDirectoryRecursively(jcli, tcli, syncDirPath, ignores)
				} else {
					err = syncDirectory(jcli, tcli, syncDirPath, ignores)
				}
			} else if len(syncFilePath) > 0 {
				err = syncFile(jcli, tcli, syncFilePath, ignores)
			} else {
				fmt.Println("开始完成队列中剩余的同步任务...")
				err = syncQueueItems(jcli, tcli, ignores)
			}
			if err != nil {
				syncErrors = append(syncErrors, err)
			}

			if numWorkers > 1 {
				close(workerChan)
				wg.Wait()
			}

			if len(syncErrors) > 0 {
				printErrors()
				os.Exit(1)
			}
		},
	}
	syncCmd.Flags().BoolVarP(&syncRecursively, "recursive", "r", false, "递归同步目录")
	syncCmd.Flags().BoolVarP(&useBfs, "bfs", "", false, "使用 bfs 算法")
	syncCmd.Flags().BoolVarP(&syncAll, "all", "A", false, "同步所有文件和目录")
	syncCmd.Flags().StringVarP(&syncFilePath, "file", "f", "", "指定同步文件")
	syncCmd.Flags().StringVarP(&syncDirPath, "dir", "d", "", "指定同步目录")
	syncCmd.Flags().IntVarP(&numWorkers, "workers", "w", 1, "指定 workers 数量")
	syncCmd.Flags().StringVarP(&syncIgnorePath, "ignore", "", "", "指定 .ignore 文件的路径")
}
