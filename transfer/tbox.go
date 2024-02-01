package transfer

import (
	"github.com/spf13/cobra"
)

var tboxCmd *cobra.Command

func init() {
	tboxCmd = &cobra.Command{
		Use:   "tbox",
		Short: "tbox 相关指令",
		Long:  "tbox 相关指令",
		Args:  cobra.NoArgs,
		Run:   func(cmd *cobra.Command, args []string) {},
	}
}
