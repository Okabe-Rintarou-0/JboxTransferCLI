package transfer

import (
	"jtrans/db"
	"os"

	"github.com/spf13/cobra"
)

var clearCmd *cobra.Command

func init() {
	clearCmd = &cobra.Command{
		Use:   "clear",
		Short: "清除同步状态和任务",
		Long:  "清除同步状态和任务",
		Run: func(cmd *cobra.Command, args []string) {
			_ = os.Remove(db.FilePath)
		},
	}
}
