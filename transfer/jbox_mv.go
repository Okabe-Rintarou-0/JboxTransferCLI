package transfer

import (
	"fmt"
	"jtrans/login"
	"os"

	"github.com/spf13/cobra"
)

var jboxMvCmd *cobra.Command

func init() {
	jboxMvCmd = &cobra.Command{
		Use:   "mv",
		Short: "移动 jbox 文件",
		Long:  "移动 jbox 文件",
		Args:  cobra.MinimumNArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			jcli, _, err := login.GetClient()
			if err != nil {
				fmt.Println(NotLoginHint)
				os.Exit(1)
			}
			fromPaths := args[:len(args)-1]
			for i, fromPath := range fromPaths {
				fromPaths[i] = formatPath(fromPath)
			}
			toPath := formatPath(args[len(args)-1])

			result, err := jcli.BatchMove(fromPaths, toPath)
			if err != nil {
				fmt.Println(err.Error())
				os.Exit(1)
			}

			if len(fromPaths) > 1 {
				if result.Result != "success" {
					fmt.Println("服务器返回了错误：", result.Result)
					os.Exit(1)
				}
			} else {
				for _, entry := range result.Failed {
					fmt.Printf("文件\"%s\"移动失败!", entry.Path)
				}
			}
		},
	}
}
