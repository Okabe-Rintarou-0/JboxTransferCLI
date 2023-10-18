package transfer

import (
	"fmt"
	"github.com/spf13/cobra"
	"jtrans/login"
	"os"
)

var (
	logoutCmd *cobra.Command
)

func init() {
	logoutCmd = &cobra.Command{
		Use:   "logout",
		Short: "登出",
		Long:  "登出",
		Run: func(cmd *cobra.Command, args []string) {
			_ = os.Remove(login.SessionPath)
			fmt.Println("登出成功！")
		},
	}
	rootCmd.AddCommand(logoutCmd)
}
