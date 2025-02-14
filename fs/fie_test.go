package fs

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestReadAllText(t *testing.T) {
	fmt.Println(" TestReadAllText")
	dir := createTempDir(t)
	createFile(t, dir + "/a.txt", 8, 'a')
	content, err := ReadAllText(dir + "/a.txt")
	if err != nil {
		t.Error(err)
	}
	assert.True(t, content == "aaaaaaaa")
	os.RemoveAll(dir)	
}

func TestFileCopy(t *testing.T) {
	fmt.Println(" TestFileCopy")
	dir := createTempDir(t)
	createFile(t, dir + "/a.txt", 8, 'a')
	CopyFile(dir + "/a.txt", dir + "/b.txt")
	content, err := ReadAllText(dir + "/b.txt")
	if err != nil {
		t.Error(err)
	}
	assert.True(t, content == "aaaaaaaa")
	os.RemoveAll(dir)	
}

func TestFileCopyStat(t *testing.T) {
	fmt.Println(" TestFileCopyStat")
	dir := createTempDir(t)
	createFile(t, dir + "/a.txt", 8, 'a')
	time.Sleep(time.Second)
	CopyFile(dir + "/a.txt", dir + "/b.txt")

	a_stat, err := os.Stat(dir + "/a.txt")
	if err != nil {
		t.Error(err)
	}
	b_stat, err := os.Stat(dir + "/b.txt")
	if err != nil {
		t.Error(err)
	}
	assert.True(t, a_stat.ModTime().UnixMilli() < b_stat.ModTime().UnixMilli())

	CopyFileTime(dir + "/a.txt", dir + "/b.txt")
	a_stat, err = os.Stat(dir + "/a.txt")
	if err != nil {
		t.Error(err)
	}
	b_stat, err = os.Stat(dir + "/b.txt")
	if err != nil {
		t.Error(err)
	}
	assert.True(t, a_stat.ModTime() == b_stat.ModTime())
	os.RemoveAll(dir)	
}
