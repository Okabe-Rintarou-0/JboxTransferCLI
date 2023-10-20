package transfer

import (
	"fmt"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "jtrans",
	Short: "Jbox 文件同步命令行工具",
	Long:  "Jbox 文件同步命令行工具",
}

func initCmd() {
	rootCmd.AddCommand(jboxCmd)
	rootCmd.AddCommand(tboxCmd)
	rootCmd.AddCommand(diffCmd)
	rootCmd.AddCommand(syncCmd)
	rootCmd.AddCommand(crc64Cmd)
	rootCmd.AddCommand(md5Cmd)

	syncCmd.AddCommand(clearCmd)
	syncCmd.AddCommand(syncTasksCmd)
	syncCmd.AddCommand(syncRecoverCmd)
	syncCmd.AddCommand(syncRestartCmd)

	jboxCmd.AddCommand(downloadCmd)
	jboxCmd.AddCommand(listJboxCmd)
	jboxCmd.AddCommand(infoJboxCmd)

	tboxCmd.AddCommand(uploadCmd)
	tboxCmd.AddCommand(listTboxCmd)
	tboxCmd.AddCommand(mkdirCmd)
	tboxCmd.AddCommand(infoTboxCmd)
}

func Execute() {
	initCmd()

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err.Error())
	}
}
