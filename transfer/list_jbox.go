package transfer

import (
	"fmt"
	"jtrans/jbox"
	"jtrans/jbox/models"
	"jtrans/login"
	"jtrans/utils"
	"os"
	"sort"

	"github.com/spf13/cobra"
)

var listJboxCmd *cobra.Command

func listJbox(jcli jbox.IClient, targetDir string) error {
	var (
		i     int
		file  *models.FileInfo
		dir   *models.DirectoryInfo
		err   error
		dirs  []*models.DirectoryInfo
		files []*models.FileInfo
	)
	targetDir = formatPath(targetDir)
	page := 0
	for {
		dir, err = jcli.GetDirectoryInfo(targetDir, page)
		if err != nil {
			return err
		}
		if len(dir.Content) == 0 {
			break
		}
		curJDirs, curJFiles := listJboxDir(dir)
		for _, dir = range curJDirs {
			dirs = append(dirs, dir)
		}
		for _, file = range curJFiles {
			files = append(files, file)
		}
		page += 1
	}

	sort.SliceStable(dirs, func(i, j int) bool {
		return dirs[i].Path < dirs[j].Path
	})
	sort.SliceStable(files, func(i, j int) bool {
		return files[i].Path < files[j].Path
	})

	for i, dir = range dirs {
		if i > 0 {
			fmt.Printf(" ")
		}
		fmt.Printf(utils.GetFileName(dir.Path))
	}

	for _, file = range files {
		if i > 0 {
			fmt.Printf(" ")
		}
		fmt.Printf(utils.GetFileName(file.Path))
		i += 1
	}

	fmt.Println()
	return nil
}

func init() {
	listJboxCmd = &cobra.Command{
		Use:   "ls",
		Short: "列出 jbox 指定目录下的文件和目录",
		Long:  "列出 jbox 指定目录下的文件和目录",
		Args:  cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			jcli, _, err := login.GetClient()
			if err != nil {
				fmt.Println(NotLoginHint)
				os.Exit(1)
			}
			dir := "/"
			if len(args) > 0 {
				dir = args[0]
			}
			if err = listJbox(jcli, dir); err != nil {
				fmt.Printf("出现了预期之外的错误：%s\n", err.Error())
				os.Exit(1)
			}
		},
	}
}
