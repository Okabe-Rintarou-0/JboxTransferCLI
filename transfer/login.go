package transfer

import (
	"fmt"
	"jtrans/login"

	"github.com/spf13/cobra"
)

var (
	loginCmd    *cobra.Command
	loginMethod login.Method
	useQRCode   bool
)

func init() {
	loginCmd = &cobra.Command{
		Use:   "login",
		Short: "登录",
		Long:  "登录",
		Run: func(cmd *cobra.Command, args []string) {
			if useQRCode {
				loginMethod = login.Qrcode
			}
			if err := login.Login(loginMethod); err != nil {
				fmt.Println(err.Error())
			}
		},
	}
	loginCmd.Flags().BoolVarP(&useQRCode, "qrcode", "", false, "使用二维码登录")
	rootCmd.AddCommand(loginCmd)
}
