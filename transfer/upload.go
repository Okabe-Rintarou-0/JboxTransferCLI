package transfer

import (
	"fmt"
	"io"
	"jtrans/constants"
	"jtrans/login"
	"jtrans/tbox"
	"jtrans/tbox/models"
	"jtrans/utils"
	"jtrans/worker"
	"os"

	"github.com/schollz/progressbar/v3"
	"github.com/spf13/cast"
	"github.com/spf13/cobra"
)

var (
	uploadCmd            *cobra.Command
	uploadFrom           string
	uploadTo             string
	uploadFromByteStream bool
)

func uploadOne(tcli tbox.IClient) error {
	var (
		err     error
		file    *os.File
		confirm *models.ConfirmChunkUploadResult
		data    []byte
		chunk   int64
		stat    os.FileInfo
	)

	bar := progressbar.NewOptions(cast.ToInt(constants.ChunkSize),
		progressbar.OptionEnableColorCodes(true),
		progressbar.OptionShowBytes(true),
		progressbar.OptionSetWidth(15))

	if !uploadFromByteStream {
		file, err = os.Open(uploadFrom)
		if err != nil {
			return fmt.Errorf("指定的文件\"%s\"不存在！", uploadFrom)
		}

		stat, err = file.Stat()
		if err != nil {
			return fmt.Errorf("读取指定文件\"%s\"的信息时出错：%s", uploadFrom, err.Error())
		}

		if stat.IsDir() {
			return fmt.Errorf("\"%s\"必须是一个文件！", uploadFrom)
		}
	} else {
		file = os.Stdin
	}

	data, err = io.ReadAll(file)
	if err != nil {
		return fmt.Errorf("读取指定文件\"%s\"时出错：%s", uploadFrom, err.Error())
	}

	size := int64(len(data))
	chunkCount := utils.ComputeChunkCount(size, constants.ChunkSize)

	w := worker.NewTBoxUploadWorker(tcli, uploadTo, size, func(uploaded int64, total int64) {
		_ = bar.Set64(uploaded)
	})
	if err = w.PrepareForUpload(); err != nil {
		return fmt.Errorf("准备上传时出现错误:%s", err.Error())
	}

	for chunk = 1; chunk <= chunkCount; chunk++ {
		if err = w.EnsureNoExpire(chunk); err != nil {
			return fmt.Errorf("上传时出现错误:%s", err.Error())
		}
		chunkSize := utils.ComputeCurrentChunkSize(chunk, chunkCount, size)
		start := (chunk - 1) * constants.ChunkSize
		toUpload := data[start : start+chunkSize]

		bar.Reset()
		bar.Describe(fmt.Sprintf("[cyan][%d/%d][reset] [red]%s[reset] 上传分块...", chunk, chunkCount, uploadTo))
		bar.ChangeMax64(chunkSize)

		if err = w.Upload(toUpload, chunk); err != nil {
			return fmt.Errorf("上传时出现错误:%s", err.Error())
		}
	}

	confirm, err = w.Confirm()
	if err != nil {
		return fmt.Errorf("上传时出现错误:%s", err.Error())
	}
	fmt.Printf("\n上传成功！文件校验值：%s\n", confirm.Crc64)
	return nil
}

func init() {
	uploadCmd = &cobra.Command{
		Use:   "upload",
		Short: "上传指定文件",
		Long:  "上传指定文件",
		Args:  cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			specifiedFrom := len(uploadFrom) > 0
			specifiedTo := len(uploadTo) > 0
			if len(args) > 0 {
				if !specifiedFrom {
					uploadFrom = args[0]
				}
				if !specifiedTo {
					uploadTo = args[0]
				}
			}

			if uploadFromByteStream && len(uploadTo) == 0 {
				fmt.Println("使用 stdin 字节流上传必须指定保存路径！")
				os.Exit(1)
			}
			uploadTo = formatPath(uploadTo)

			_, tcli, err := login.GetClient()
			if err != nil {
				fmt.Println(NotLoginHint)
				os.Exit(1)
			}

			if err = uploadOne(tcli); err != nil {
				fmt.Println(err.Error())
				os.Exit(1)
			}
		},
	}
	uploadCmd.Flags().BoolVarP(&uploadFromByteStream, "bytes", "", false, "是否从 stdin 字节流上传文件")
	uploadCmd.Flags().StringVarP(&uploadFrom, "from", "f", "", "指定要上传的文件")
	uploadCmd.Flags().StringVarP(&uploadTo, "to", "t", "", "指定上传路径")
}
