package transfer

import (
	"fmt"
	"jtrans/login"
	"os"

	"github.com/spf13/cobra"
)

var tboxMvCmd *cobra.Command

func init() {
	tboxMvCmd = &cobra.Command{
		Use:   "mv",
		Short: "移动 tbox 文件",
		Long:  "移动 tbox 文件",
		Args:  cobra.MinimumNArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			_, tcli, err := login.GetClient()
			if err != nil {
				fmt.Println(NotLoginHint)
				os.Exit(1)
			}
			fromPaths := args[:len(args)-1]
			for i, fromPath := range fromPaths {
				fromPaths[i] = formatPath(fromPath)
			}
			toPath := formatPath(args[len(args)-1])

			result, err := tcli.BatchMove(fromPaths, toPath)
			if err != nil {
				fmt.Println(err.Error())
				os.Exit(1)
			}

			for _, result := range result.Result {
				if result.Status != 200 {
					fmt.Printf("移动文件\"%s\"时出错：服务器响应%d\n", result.From, result.Status)
				}
			}
		},
	}
}
