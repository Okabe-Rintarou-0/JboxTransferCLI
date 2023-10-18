package transfer

import (
	"fmt"
	"jtrans/login"
	"jtrans/tbox/models"
	"os"

	"github.com/spf13/cobra"
)

var mkdirCmd *cobra.Command

func init() {
	mkdirCmd = &cobra.Command{
		Use:   "mkdir",
		Short: "在 tbox 中创建目录",
		Long:  "在 tbox 中创建目录",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			_, tcli, err := login.GetClient()
			if err != nil {
				fmt.Println(NotLoginHint)
				os.Exit(1)
			}
			dir := formatPath(args[0])
			var errMsg *models.ErrorMessage
			if errMsg, err = tcli.CreateDirectory(dir); err == nil {
				switch errMsg.Code {
				case "SameNameDirectoryOrFileExists":
					fmt.Printf("指定的目录\"%s\"已存在\n", dir)
					break
				case "RootDirectoryNotAllowed":
					fmt.Println("不允许创建根目录！")
					break
				default:
					fmt.Printf("创建目录\"%s\"成功！\n", dir)
				}
			} else {
				fmt.Printf("创建目录时出错：%s", err.Error())
			}
		},
	}
}
