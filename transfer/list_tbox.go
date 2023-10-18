package transfer

import (
	"fmt"
	"jtrans/login"
	"jtrans/tbox"
	"jtrans/tbox/models"
	"os"
	"sort"

	"github.com/spf13/cobra"
)

var listTboxCmd *cobra.Command

func listTbox(tcli tbox.IClient, targetDir string) error {
	var (
		i     int
		err   error
		dir   *models.DirectoryInfo
		file  *models.FileInfo
		dirs  []*models.FileInfo
		files []*models.FileInfo
	)
	targetDir = formatPath(targetDir)

	dir, err = tcli.GetDirectoryInfo(targetDir, nil, nil, "")
	if err != nil {
		return err
	}

	dirsMap, filesMap := listTboxDir(dir)
	for _, file = range dirsMap {
		dirs = append(dirs, file)
	}
	for _, file = range files {
		files = append(files, file)
	}

	sort.SliceStable(dirs, func(i, j int) bool {
		return dirs[i].Name < dirs[j].Name
	})
	sort.SliceStable(files, func(i, j int) bool {
		return files[i].Name < files[j].Name
	})

	i = 0
	for _, file = range dirsMap {
		if i > 0 {
			fmt.Printf(" ")
		}
		fmt.Printf(file.Name)
		i += 1
	}

	for _, file = range filesMap {
		if i > 0 {
			fmt.Printf(" ")
		}
		fmt.Printf(file.Name)
		i += 1
	}

	fmt.Println()
	return nil
}

func init() {
	listTboxCmd = &cobra.Command{
		Use:   "ls",
		Short: "列出 tbox 指定目录下的文件和目录",
		Long:  "列出 tbox 指定目录下的文件和目录",
		Args:  cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			_, tcli, err := login.GetClient()
			if err != nil {
				fmt.Println(NotLoginHint)
				os.Exit(1)
			}
			dir := "/"
			if len(args) > 0 {
				dir = args[0]
			}
			if err = listTbox(tcli, dir); err != nil {
				fmt.Printf("出现了预期之外的错误：%s\n", err.Error())
				os.Exit(1)
			}
		},
	}
}
