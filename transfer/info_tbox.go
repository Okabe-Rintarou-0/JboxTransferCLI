package transfer

import (
	"fmt"
	"jtrans/login"
	"jtrans/tbox"
	"jtrans/tbox/models"
	"jtrans/utils"
	"os"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/spf13/cast"
	"github.com/spf13/cobra"
)

var infoTboxCmd *cobra.Command

func tboxFileInfoTbl() table.Writer {
	t := table.NewWriter()
	t.SetAutoIndex(true)
	t.Style().Format.Header = text.FormatTitle
	t.AppendHeader(table.Row{"路径", "类型", "大小", "校验和(crc64)"})
	return t
}

func infoTbox(tcli tbox.IClient, t table.Writer, filePath string) error {
	var (
		file *models.FileInfo
		err  error
	)
	file, err = tcli.GetFileInfo(filePath)
	if err != nil {
		return err
	}
	tp := "文件"
	if file.IsDir() {
		tp = "目录"
	}
	t.AppendRow(table.Row{file.FullPath(), tp, utils.FormatBytes(cast.ToFloat64(file.Size)), file.Crc64})
	return nil
}

func init() {
	infoTboxCmd = &cobra.Command{
		Use:   "info",
		Short: "查看 tbox 文件或目录信息",
		Long:  "查看 tbox 文件或目录信息",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			_, tcli, err := login.GetClient()
			if err != nil {
				fmt.Println(NotLoginHint)
				os.Exit(1)
			}
			tbl := tboxFileInfoTbl()
			for _, path := range args {
				if err = infoTbox(tcli, tbl, formatPath(path)); err != nil {
					fmt.Printf("在读取文件（目录）\"%s\"信息时出错：%s，跳过该文件（目录）\n", path, err.Error())
				}
			}
			fmt.Println(tbl.Render())
		},
	}
}
