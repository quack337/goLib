package fs

import (
	"errors"
	"os"
	"sort"
)

type FileInfo struct {
	Name string
	ModTime int64
	Size int64 
}

func GetEntries(dirPath string) (files []FileInfo, dirs []string, err error) {
    entries, err := os.ReadDir(dirPath)
    if err != nil { return files, dirs, err } 
    for _, entry := range entries {
		if entry.IsDir() {
			dirs = append(dirs, entry.Name())
		} else {
			var fileInfo, err = entry.Info()
			if err != nil { return files, dirs, err } 
			files = append(files, FileInfo{
				entry.Name(),
				fileInfo.ModTime().UnixMilli(),
				fileInfo.Size()})
		}
	}
	return files, dirs, nil
}

func CompareByNameSizeTime(a *FileInfo, b *FileInfo) bool {
	return a.Name == b.Name && a.ModTime == b.ModTime && a.Size == b.Size
}

type FileInfos []FileInfo

func (fi FileInfos) Len() int { return len(fi) }
func (fi FileInfos) Less(i, j int) bool { return fi[i].Name < fi[j].Name }
func (fi FileInfos) Swap(i, j int) { fi[i], fi[j] = fi[j], fi[i] }

func SortFileInfos(fileInfos []FileInfo) {
	sort.Sort(FileInfos(fileInfos))
}

func BinarySearchFileInfoByName(fileInfos []FileInfo, name string) int {
	var left = 0
	var right = len(fileInfos) - 1
	for left <= right {
		var middle = (left + right) / 2
		if fileInfos[middle].Name < name {
			left = middle + 1
		} else if fileInfos[middle].Name > name {
			right = middle - 1
		} else {
			return middle		
		}
	}
	return -1
}

func DirExists(dirPath string) (bool, error) {
	var info, err = os.Stat(dirPath)
	if os.IsNotExist(err)  {
		return false, nil
	}
	if info.IsDir() {
		return true, nil
	} else {
		return false, errors.New(dirPath + " is file")
	}
}