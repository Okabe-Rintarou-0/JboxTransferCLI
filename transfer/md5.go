package transfer

import (
	"bufio"
	"fmt"
	"io"
	"jtrans/constants"
	"jtrans/encrypt"
	"os"

	"github.com/spf13/cobra"
)

var md5Cmd *cobra.Command

func init() {
	md5Cmd = &cobra.Command{
		Use:   "md5",
		Short: "从 stdin 接收字节流，计算 md5 校验和",
		Long:  "从 stdin 接收字节流，计算 md5 校验和",
		Run: func(cmd *cobra.Command, args []string) {
			reader := bufio.NewReader(os.Stdin)
			md5 := encrypt.NewMD5()
			for {
				buffer := make([]byte, constants.ChunkSize)
				n, err := reader.Read(buffer)
				if err == io.EOF {
					break
				}
				if err != nil {
					fmt.Printf("读取字节流时出错： %s\n", err.Error())
					os.Exit(1)
				}
				encrypt.MD5HashProc(md5, buffer[:n])
			}
			fmt.Println(encrypt.MD5HashProcFinish(md5))
		},
	}
}
