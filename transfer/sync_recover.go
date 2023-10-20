package transfer

import (
	"fmt"
	"jtrans/db"
	"jtrans/login"
	"os"

	"github.com/spf13/cobra"
)

var (
	syncRecoverCmd *cobra.Command
	recoverAll     bool
)

func init() {
	syncRecoverCmd = &cobra.Command{
		Use:   "recover",
		Short: "恢复失败的任务",
		Long:  "恢复失败的任务",
		Run: func(cmd *cobra.Command, args []string) {
			if !login.CheckLogin() {
				fmt.Println(NotLoginHint)
				os.Exit(1)
			}
			paths := args

			if recoverAll {
				db.RecoverAllFailedTasks()
			} else if len(paths) == 0 {
				fmt.Println("请正确指定需要恢复的任务！")
				os.Exit(1)
			} else {
				db.RecoverFailedTasks(paths)
			}
			fmt.Println("恢复成功！")
		},
	}
	syncRecoverCmd.Flags().BoolVarP(&recoverAll, "all", "A", false, "恢复全部")
}
