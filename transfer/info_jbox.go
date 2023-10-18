package transfer

import (
	"fmt"
	"jtrans/jbox"
	"jtrans/jbox/models"
	"jtrans/login"
	"jtrans/utils"
	"os"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/spf13/cobra"
)

var infoJboxCmd *cobra.Command

func jboxFileInfoTbl() table.Writer {
	t := table.NewWriter()
	t.SetAutoIndex(true)
	t.Style().Format.Header = text.FormatTitle
	t.AppendHeader(table.Row{"路径", "类型", "大小", "校验和(md5)"})
	return t
}

func infoJbox(jcli jbox.IClient, t table.Writer, filePath string) error {
	var (
		file *models.FileInfo
		err  error
	)
	file, err = jcli.GetFileInfo(filePath)
	if err != nil {
		return err
	}
	tp := "文件"
	if file.IsDir {
		tp = "目录"
	}
	t.AppendRow(table.Row{file.Path, tp, utils.FormatBytes(float64(file.Bytes)), file.Hash})
	return nil
}

func init() {
	infoJboxCmd = &cobra.Command{
		Use:   "info",
		Short: "查看 jbox 文件或目录信息",
		Long:  "查看 jbox 文件或目录信息",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			jcli, _, err := login.GetClient()
			if err != nil {
				fmt.Println(NotLoginHint)
				os.Exit(1)
			}
			tbl := jboxFileInfoTbl()
			for _, path := range args {
				if err = infoJbox(jcli, tbl, formatPath(path)); err != nil {
					fmt.Printf("在读取文件（目录）\"%s\"信息时出错：%s，跳过该文件（目录）\n", path, err.Error())
				}
			}
			fmt.Println(tbl.Render())
		},
	}
}
