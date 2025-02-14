package fs

import (
	"fmt"
	"os"
	"time"

	"testing"

	"github.com/stretchr/testify/assert"
)

func createFile(t *testing.T, path string, size int, data byte) {
	f, err := os.Create(path)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	a := make([]byte, size)
	for i := 0; i < size; i++ {
		a[i] = data
	}
	_, err = f.Write(a)
	if err != nil {
		t.Fatal(err)
	}
}

func createTempDir(t *testing.T) string {
	dir, err := os.MkdirTemp("", "sampledir")
	if err != nil {
		t.Fatal(err)
	}
	return dir
}

func createDir(t *testing.T, path string) {
	err := os.Mkdir(path, os.ModePerm)
	if err != nil {
		t.Fatal(err)
	}
}

func TestGetEntries(t *testing.T) {
	fmt.Println(" TestGetEntires")

	time0 := time.Now().UnixMilli()
	dir := createTempDir(t)
	defer os.RemoveAll(dir)
	createFile(t, dir + "/a.txt", 8, 'a')
	createFile(t, dir + "/b.txt", 6, 'b')
	createDir(t, dir + "/sub")
	createFile(t, dir + "/sub/c.txt", 10, 'c')
	createFile(t, dir + "/sub/d.txt", 12, 'd')
	time1 := time.Now().UnixMilli()

	files, dirs, err := GetEntries(dir)
	if err != nil {
		t.Fatal(err)
	}
	assert.True(t, len(files) == 2, "len(files) == ", len(files))
	assert.True(t, files[0].Name == "a.txt", "files[0].Name == ", files[0].Name)
	assert.True(t, files[1].Name == "b.txt", "files[1].Name == ", files[1].Name)
	assert.True(t, files[0].Size == 8, "files[0].Size == ", files[0].Size)
	assert.True(t, files[1].Size == 6, "files[1].Size == ", files[1].Size)
	assert.True(t, files[0].ModTime >= time0 && files[0].ModTime <= time1, "files[0].ModTime == ", files[0].ModTime)
	assert.True(t, files[1].ModTime >= time0 && files[1].ModTime <= time1, "files[1].ModTime == ", files[1].ModTime)

	assert.True(t, len(dirs) == 1, "len(dirs) == ", len(dirs))
	assert.True(t, dirs[0] == "sub", "dirs[0] == ", dirs[0])

	files, dirs, err = GetEntries(dir + "/sub")
	if err != nil {
		t.Fatal(err)
	}
	assert.True(t, len(files) == 2, "len(files) == ", len(files))
	assert.True(t, files[0].Name == "c.txt", "files[0].Name == ", files[0].Name)
	assert.True(t, files[1].Name == "d.txt", "files[1].Name == ", files[1].Name)
	assert.True(t, files[0].Size == 10, "files[0].Size == ", files[0].Size)
	assert.True(t, files[1].Size == 12, "files[1].Size == ", files[1].Size)
	assert.True(t, files[0].ModTime >= time0 && files[0].ModTime <= time1, "files[0].ModTime == ", files[0].ModTime)
	assert.True(t, files[1].ModTime >= time0 && files[1].ModTime <= time1, "files[1].ModTime == ", files[1].ModTime)

	assert.True(t, len(dirs) == 0, "len(dirs) == ", len(dirs))
}

func TestCompareByNameSizeTime(t *testing.T) {
	fmt.Println(" TestCompareByNameSizeTime")
	dir := createTempDir(t)
	defer os.RemoveAll(dir)

	createFile(t, dir + "/d.txt", 8, 'a')
	createFile(t, dir + "/b.txt", 6, 'b')
	var files, _, err = GetEntries(dir)
	if err != nil {
		t.Fatal(err)
	}
	assert.True(t, len(files) == 2)
	assert.True(t, CompareByNameSizeTime(&files[0], &files[0]))
	assert.True(t, CompareByNameSizeTime(&files[1], &files[1]))
	assert.False(t, CompareByNameSizeTime(&files[0], &files[1]))
	assert.False(t, CompareByNameSizeTime(&files[1], &files[0]))
	var file0 = files[0]
	var file1 = files[1]
	assert.True(t, CompareByNameSizeTime(&files[0], &file0))
	assert.True(t, CompareByNameSizeTime(&files[1], &file1))
	assert.False(t, CompareByNameSizeTime(&files[0], &file1))
	assert.False(t, CompareByNameSizeTime(&files[1], &file0))
}

func TestSortFileInfos(t *testing.T) {
	fmt.Println(" TestSortFileInfos")
	dir := createTempDir(t)
    defer os.RemoveAll(dir)

	var files, _, err = GetEntries(dir)
	if err != nil {
		t.Fatal(err)
	}
	assert.True(t, len(files) == 0)
	SortFileInfos(files)
	assert.True(t, len(files) == 0)

	createFile(t, dir + "/d.txt", 8, 'a')
	files, _, err = GetEntries(dir)
	if err != nil {
		t.Fatal(err)
	}
	assert.True(t, len(files) == 1)
	SortFileInfos(files)
	assert.True(t, len(files) == 1)
	assert.True(t, files[0].Name == "d.txt")

	createFile(t, dir + "/b.txt", 8, 'a')
	files, _, err = GetEntries(dir)
	if err != nil {
		t.Fatal(err)
	}
	assert.True(t, len(files) == 2)
	SortFileInfos(files)
	assert.True(t, len(files) == 2)
	assert.True(t, files[0].Name == "b.txt")
	assert.True(t, files[1].Name == "d.txt")

	createFile(t, dir + "/a.txt", 8, 'a')
	files, _, err = GetEntries(dir)
	if err != nil {
		t.Fatal(err)
	}
	assert.True(t, len(files) == 3)
	SortFileInfos(files)
	assert.True(t, len(files) == 3)
	assert.True(t, files[0].Name == "a.txt")
	assert.True(t, files[1].Name == "b.txt")
	assert.True(t, files[2].Name == "d.txt")

	createFile(t, dir + "/c.txt", 8, 'a')
	files, _, err = GetEntries(dir)
	if err != nil {
		t.Fatal(err)
	}
	assert.True(t, len(files) == 4)
	SortFileInfos(files)
	assert.True(t, len(files) == 4)
	assert.True(t, files[0].Name == "a.txt")
	assert.True(t, files[1].Name == "b.txt")
	assert.True(t, files[2].Name == "c.txt")
	assert.True(t, files[3].Name == "d.txt")
}

func TestBinarySearchFileInfoByName(t *testing.T) {
	fmt.Println(" TestBinarySearchFileInfoByName")

	dir := createTempDir(t)
    defer os.RemoveAll(dir)

	files, _, err := GetEntries(dir)
	if err != nil {
		t.Fatal(err)
	}
	assert.True(t, BinarySearchFileInfoByName(files, "a.txt") < 0)

	createFile(t, dir + "/b.txt", 8, 'a')
	files, _, err = GetEntries(dir)
	if err != nil {
		t.Fatal(err)
	}
	assert.True(t, BinarySearchFileInfoByName(files, "a.txt") < 0)
	assert.True(t, BinarySearchFileInfoByName(files, "b.txt") == 0)
	assert.True(t, BinarySearchFileInfoByName(files, "c.txt") < 0)

	createFile(t, dir + "/c.txt", 8, 'a')
	files, _, err = GetEntries(dir)
	if err != nil {
		t.Fatal(err)
	}
	assert.True(t, BinarySearchFileInfoByName(files, "a.txt") < 0)
	assert.True(t, BinarySearchFileInfoByName(files, "b.txt") == 0)
	assert.True(t, BinarySearchFileInfoByName(files, "c.txt") == 1)
	assert.True(t, BinarySearchFileInfoByName(files, "d.txt") < 0)

	createFile(t, dir + "/d.txt", 8, 'a')
	files, _, err = GetEntries(dir)
	if err != nil {
		t.Fatal(err)
	}
	assert.True(t, BinarySearchFileInfoByName(files, "a.txt") < 0)
	assert.True(t, BinarySearchFileInfoByName(files, "b.txt") == 0)
	assert.True(t, BinarySearchFileInfoByName(files, "c.txt") == 1)
	assert.True(t, BinarySearchFileInfoByName(files, "d.txt") == 2)
	assert.True(t, BinarySearchFileInfoByName(files, "e.txt") < 0)
}

func TestDirExists(t *testing.T) {
	fmt.Println(" TestDirExists")

	dir := createTempDir(t)
    defer os.RemoveAll(dir)
	
	flag, err := DirExists(dir)
	assert.True(t, flag && err == nil)
	
	flag, err = DirExists(dir + "_a_b_c_d_")
	assert.True(t, !flag && err == nil)

	createFile(t, dir + "/a", 8, 'a')
	flag, err = DirExists(dir + "/a")
	assert.True(t, !flag && err != nil)
}