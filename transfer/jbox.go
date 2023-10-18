package transfer

import (
	"github.com/spf13/cobra"
)

var jboxCmd *cobra.Command

func init() {
	jboxCmd = &cobra.Command{
		Use:   "jbox",
		Short: "jbox 相关指令",
		Long:  "jbox 相关指令",
		Run: func(cmd *cobra.Command, args []string) {

		},
	}
}
