package transfer

import (
	"fmt"
	"jtrans/jbox"
	jmodels "jtrans/jbox/models"
	"jtrans/login"
	"jtrans/tbox"
	tmodels "jtrans/tbox/models"
	"os"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/spf13/cobra"
)

var diffCmd *cobra.Command

func init() {
	diffCmd = &cobra.Command{
		Use:   "diff",
		Short: "比较新老云盘",
		Long:  "比较新老云盘",
		Run: func(cmd *cobra.Command, args []string) {
			jcli, tcli, err := login.GetClient()
			if err != nil {
				fmt.Println(NotLoginHint)
				os.Exit(1)
			}
			if result, err := computeDiff(jcli, tcli); err == nil {
				fmt.Println("可以直接同步传输的目录或文件：")
				transferTbl := canTransferTbl()
				for _, entry := range result.CanTransfer {
					file := entry.GetFile()
					fileType := "文件"
					if file.IsDir {
						fileType = "目录"
					}
					transferTbl.AppendRow(table.Row{file.Path, fileType, file.Size, file.Hash})
				}
				fmt.Println(transferTbl.Render())
				fmt.Println("冲突的目录或文件：")
				cxTbl := conflictTbl()
				for _, entry := range result.Conflicts {
					file := entry.GetFile()
					fileType := "文件"
					if file.IsDir {
						fileType = "目录"
					}
					cxTbl.AppendRow(table.Row{file.Path, fileType, entry.Reason})
				}
				fmt.Println(cxTbl.Render())
			} else {
				fmt.Printf("出错：%s\n", err.Error())
				os.Exit(1)
			}
		},
	}
}

const (
	TBoxExistsSameNameDirectoryConflict = "新版云盘存在相同路径的目录"
	TBoxExistsSameNameFileConflict      = "新版云盘存在相同路径的文件"
)

type fileEntry struct {
	File *jmodels.FileInfo
	Dir  *jmodels.DirectoryInfo
}

func (e fileEntry) GetFile() *jmodels.FileInfo {
	if e.File != nil {
		return e.File
	} else if e.Dir != nil {
		return &e.Dir.FileInfo
	}
	return nil
}

type conflictEntry struct {
	fileEntry
	Reason string
}

type diffResult struct {
	CanTransfer []fileEntry
	Conflicts   []conflictEntry
}

func displayJboxDir(dir *jmodels.DirectoryInfo) {
	fmt.Printf("当前目录：%s\n", dir.Path)
	for _, file := range dir.Content {
		if file.IsDir {
			fmt.Printf("目录底下有目录：%s hash:%s\n", file.Path, file.Hash)
		} else {
			fmt.Printf("目录底下有文件：%s hash:%s\n", file.Path, file.Hash)
		}
	}
}

func listJboxDir(dir *jmodels.DirectoryInfo) (dirs map[string]*jmodels.DirectoryInfo, files map[string]*jmodels.FileInfo) {
	dirs = make(map[string]*jmodels.DirectoryInfo)
	files = make(map[string]*jmodels.FileInfo)
	for _, file := range dir.Content {
		if file.IsDir {
			dirs[file.Path] = file
		} else {
			files[file.Path] = &file.FileInfo
		}
	}
	return
}

func displayTboxDir(dir *tmodels.DirectoryInfo) {
	fmt.Printf("当前目录：%s\n", dir.FullPath())
	for _, file := range dir.Contents {
		if file.Type != "file" {
			fmt.Printf("目录底下有目录：%s hash:%s\n", file.FullPath(), file.Crc64)
		} else {
			fmt.Printf("目录底下有文件：%s hash:%s\n", file.FullPath(), file.Crc64)
		}
	}
}

func listTboxDir(dir *tmodels.DirectoryInfo) (dirs map[string]*tmodels.FileInfo, files map[string]*tmodels.FileInfo) {
	dirs = make(map[string]*tmodels.FileInfo)
	files = make(map[string]*tmodels.FileInfo)
	for _, file := range dir.Contents {
		if file.IsDir() {
			dirs[file.FullPath()] = file
		} else {
			files[file.FullPath()] = file
		}
	}
	return
}

func computeDiffInner(currentDir string, jcli jbox.IClient, tcli tbox.IClient) (canTransfer []fileEntry, conflicts []conflictEntry, err error) {
	var (
		jdir   *jmodels.DirectoryInfo
		tdir   *tmodels.DirectoryInfo
		jdirs  = make(map[string]*jmodels.DirectoryInfo)
		jfiles = make(map[string]*jmodels.FileInfo)
	)

	page := 0
	currentDir = formatPath(currentDir)
	for {
		jdir, err = jcli.GetDirectoryInfo(currentDir, page)
		if len(jdir.Content) == 0 {
			break
		}
		if err != nil {
			return nil, nil, err
		}
		curJDirs, curJFiles := listJboxDir(jdir)
		for path, dir := range curJDirs {
			jdirs[path] = dir
		}
		for path, file := range curJFiles {
			jfiles[path] = file
		}

		page += 1
	}
	//displayJboxDir(jdir)

	tdir, err = tcli.GetDirectoryInfo(currentDir,
		nil,
		&tmodels.OrderOption{
			By:   "name",
			Type: tmodels.OrderByAsc,
		},
		"")
	if err != nil {
		return nil, nil, err
	}

	//displayTboxDir(tdir)
	tdirs, tfiles := listTboxDir(tdir)

	for jpath, jfile := range jfiles {
		if _, ok := tdirs[jpath]; ok {
			conflicts = append(conflicts, conflictEntry{
				fileEntry: fileEntry{File: jfile},
				Reason:    TBoxExistsSameNameDirectoryConflict,
			})
		} else if _, ok = tfiles[jpath]; ok {
			conflicts = append(conflicts, conflictEntry{
				fileEntry: fileEntry{File: jfile},
				Reason:    TBoxExistsSameNameFileConflict,
			})
		} else {
			canTransfer = append(canTransfer, fileEntry{
				File: jfile,
			})
		}
	}

	for jpath, jdir := range jdirs {
		if _, ok := tdirs[jpath]; ok {
			var (
				subCanTransfer []fileEntry
				subConflict    []conflictEntry
			)
			subCanTransfer, subConflict, err = computeDiffInner(jpath, jcli, tcli)
			if err != nil {
				return nil, nil, err
			}
			canTransfer = append(canTransfer, subCanTransfer...)
			conflicts = append(conflicts, subConflict...)
		} else if _, ok = tfiles[jpath]; ok {
			conflicts = append(conflicts, conflictEntry{
				fileEntry: fileEntry{Dir: jdir},
				Reason:    TBoxExistsSameNameFileConflict,
			})
		} else {
			canTransfer = append(canTransfer, fileEntry{
				Dir: jdir,
			})
		}
	}

	return canTransfer, conflicts, nil
}

func canTransferTbl() table.Writer {
	t := table.NewWriter()
	t.SetAutoIndex(true)
	t.Style().Format.Header = text.FormatTitle
	t.AppendHeader(table.Row{"文件路径", "文件类型", "文件大小", "文件Hash"})
	return t
}

func conflictTbl() table.Writer {
	t := table.NewWriter()
	t.SetAutoIndex(true)
	t.Style().Format.Header = text.FormatTitle
	t.AppendHeader(table.Row{"文件路径", "文件类型", "冲突原因"})
	return t
}

func computeDiff(jcli jbox.IClient, tcli tbox.IClient) (*diffResult, error) {
	var (
		canTransfer []fileEntry
		conflicts   []conflictEntry
		err         error
	)

	canTransfer, conflicts, err = computeDiffInner("/", jcli, tcli)
	if err != nil {
		return nil, err
	}

	return &diffResult{
		CanTransfer: canTransfer,
		Conflicts:   conflicts,
	}, nil
}
