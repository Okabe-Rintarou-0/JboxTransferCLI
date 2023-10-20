package transfer

import (
	"fmt"
	"jtrans/db"
	"jtrans/login"
	"os"

	"github.com/spf13/cobra"
)

var (
	syncRestartCmd *cobra.Command
	restartAll     bool
)

func init() {
	syncRestartCmd = &cobra.Command{
		Use:   "restart",
		Short: "重启失败的任务",
		Long:  "重启失败的任务",
		Run: func(cmd *cobra.Command, args []string) {
			if !login.CheckLogin() {
				fmt.Println(NotLoginHint)
				os.Exit(1)
			}
			paths := args
			if restartAll {
				db.RestartAllFailedTasks()
			} else if len(paths) == 0 {
				fmt.Println("请正确指定需要重启的任务！")
				os.Exit(1)
			} else {
				db.RestartFailedTasks(paths)
			}
			fmt.Println("重启成功！")
		},
	}
	syncRestartCmd.Flags().BoolVarP(&restartAll, "all", "A", false, "重启全部")
}
