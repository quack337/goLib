package fs

import (
	"io"
	"os"
)

func CopyFile(src, dst string) (int64, error) {
	source, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return 0, err
	}
	defer destination.Close()
	nBytes, err := io.Copy(destination, source)
	return nBytes, err
}

func CopyFileTime(src, dst string) error {
	stat, err := os.Stat(src)
	if err != nil {
		return err
	}
	err = os.Chtimes(dst, stat.ModTime(), stat.ModTime())
	return err
}

func ReadAllText(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(data), nil
}