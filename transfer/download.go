package transfer

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"jtrans/constants"
	"jtrans/encrypt"
	"jtrans/jbox"
	"jtrans/login"
	"jtrans/utils"
	"os"

	"github.com/schollz/progressbar/v3"
	"github.com/spf13/cast"
	"github.com/spf13/cobra"
)

var (
	downloadCmd          *cobra.Command
	shouldOverwrite      bool
	downloadAsByteStream bool
)

func refreshBar(bar *progressbar.ProgressBar, max int64, desc string) {
	bar.Reset()
	bar.ChangeMax64(max)
	bar.Describe(desc)
}

func downloadOneAsByteStream(jcli jbox.IClient, path string) error {
	path = formatPath(path)
	fileInfo, err := jcli.GetFileInfo(path)
	if err != nil {
		return err
	}

	if fileInfo.IsDir {
		return fmt.Errorf("无法下载目录！")
	}

	size := fileInfo.Bytes
	expectedHash := fileInfo.Hash
	chunkCount := utils.ComputeChunkCount(size, constants.ChunkSize)
	var (
		chunkNo    int64 = 1
		chunkData  []byte
		sha256     string
		sha256List string
	)
	for ; chunkNo <= chunkCount; chunkNo++ {
		curChunkSize := utils.ComputeCurrentChunkSize(chunkNo, chunkCount, size)
		chunkData, err = jcli.GetChunk(path, chunkNo, curChunkSize, nil)
		if err != nil {
			return err
		}

		_, err = io.CopyN(os.Stdout, bytes.NewReader(chunkData), curChunkSize)
		if err != nil {
			return err
		}

		sha256 = encrypt.SHA256Hash(chunkData)
		if len(sha256List) > 0 {
			sha256List += ","
		}
		sha256List += sha256
	}

	actualHash := encrypt.MD5Hash([]byte(sha256List))
	if actualHash != expectedHash {
		return fmt.Errorf("校验和错误！")
	}
	return nil
}

func downloadOne(jcli jbox.IClient, bar *progressbar.ProgressBar, path string) error {
	path = formatPath(path)
	fileInfo, err := jcli.GetFileInfo(path)
	if err != nil {
		return err
	}

	if fileInfo.IsDir {
		return fmt.Errorf("无法下载目录！")
	}

	outPath := utils.GetFileName(path)
	var file *os.File
	if _, err = os.Stat(outPath); err == nil && !shouldOverwrite {
		fmt.Printf("文件\"%s\"已存在，跳过下载\n", outPath)
		return nil
	}

	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return err
	}

	file, err = os.OpenFile(outPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	defer file.Close()

	size := fileInfo.Bytes
	expectedHash := fileInfo.Hash
	chunkCount := utils.ComputeChunkCount(size, constants.ChunkSize)
	var (
		chunkNo    int64 = 1
		chunkData  []byte
		sha256     string
		sha256List string
	)
	for ; chunkNo <= chunkCount; chunkNo++ {
		curChunkSize := utils.ComputeCurrentChunkSize(chunkNo, chunkCount, size)
		refreshBar(bar, curChunkSize, fmt.Sprintf("[cyan][%d/%d][reset] [red]%s[reset] 下载分块...", chunkNo, chunkCount, outPath))
		chunkData, err = jcli.GetChunk(path, chunkNo, curChunkSize, func(downloaded int64, total int64) {
			_ = bar.Set64(downloaded)
		})
		if err != nil {
			return err
		}

		_, err = io.CopyN(file, bytes.NewReader(chunkData), curChunkSize)
		if err != nil {
			return err
		}

		sha256 = encrypt.SHA256Hash(chunkData)
		if len(sha256List) > 0 {
			sha256List += ","
		}
		sha256List += sha256
	}

	fmt.Printf("\n下载完毕！正在检查校验和...")
	actualHash := encrypt.MD5Hash([]byte(sha256List))
	if actualHash != expectedHash {
		return fmt.Errorf("校验和错误！")
	}

	fmt.Println("成功！")
	return nil
}

func init() {
	bar := progressbar.NewOptions(
		cast.ToInt(constants.ChunkSize),
		progressbar.OptionEnableColorCodes(true),
		progressbar.OptionShowBytes(true),
		progressbar.OptionSetWidth(15),
	)

	downloadCmd = &cobra.Command{
		Use:   "download",
		Short: "下载指定文件",
		Long:  "下载指定文件",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			jcli, _, err := login.GetClient()
			if err != nil {
				fmt.Println(NotLoginHint)
				os.Exit(1)
			}

			// only support download one file
			if downloadAsByteStream {
				file := formatPath(args[0])
				if err = downloadOneAsByteStream(jcli, file); err != nil {
					fmt.Printf("\n下载文件\"%s\"时出错：%s\n", file, err.Error())
					os.Exit(1)
				}
				return
			}

			for _, file := range args {
				file = formatPath(file)
				if err = downloadOne(jcli, bar, file); err != nil {
					fmt.Printf("\n下载文件\"%s\"时出错：%s\n", file, err.Error())
				}
			}
		},
	}
	downloadCmd.Flags().BoolVarP(&downloadAsByteStream, "bytes", "", false, "是否以字节流的形式输出")
	downloadCmd.Flags().BoolVarP(&shouldOverwrite, "overwrite", "", false, "是否覆盖同名文件")
}
