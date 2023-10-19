package transfer

import (
	"fmt"
	"jtrans/db"
	"jtrans/login"
	"jtrans/utils"
	"os"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/spf13/cobra"
)

var (
	syncTasksCmd *cobra.Command
	orderType    = "asc"
	showFinished = false
	maxRows      = -1
)

func syncTaskQueueTbl() table.Writer {
	t := table.NewWriter()
	t.SetAutoIndex(true)
	t.Style().Format.Header = text.FormatTitle
	t.AppendHeader(table.Row{"路径", "类型", "大小", "同步状态"})
	return t
}

func printTaskQueue() {
	tasks := db.GetAllOrderedByOrder(orderType == "desc")
	tbl := syncTaskQueueTbl()
	/// 0:Idle
	/// 1:Busy
	/// 2:Error
	/// 3:Done
	/// 4:Cancel
	stateMap := map[int]string{
		0: "等待执行",
		1: "忙碌",
		2: "同步出错",
		3: "同步完成",
		4: "取消",
	}
	rows := 0
	for _, task := range tasks {
		if task.State == 3 && !showFinished {
			continue
		}
		fileType := "文件"
		if task.Type == 1 {
			fileType = "目录"
		}
		state := stateMap[task.State]
		tbl.AppendRow(table.Row{task.FilePath, fileType, utils.FormatBytes(float64(task.Size)), state})

		rows += 1
		if rows == maxRows {
			break
		}
	}

	fmt.Println("目前同步任务：")
	fmt.Println(tbl.Render())
}

func init() {
	syncTasksCmd = &cobra.Command{
		Use:   "tasks",
		Short: "查看同步队列文件",
		Long:  "查看同步队列文件",
		Run: func(cmd *cobra.Command, args []string) {
			if !login.CheckLogin() {
				fmt.Println(NotLoginHint)
				os.Exit(1)
			}

			printTaskQueue()
		},
	}
	syncTasksCmd.Flags().StringVarP(&orderType, "order", "", "asc", "排列顺序，asc 或 desc，默认 asc")
	syncTasksCmd.Flags().IntVarP(&maxRows, "max", "", -1, "打印的最大行数")
	syncTasksCmd.Flags().BoolVarP(&showFinished, "show-finished", "", false, "是否打印已完成的任务")
}
