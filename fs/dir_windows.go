//go:build windows

package fs

import (
	"syscall"
)

func IsHidden(path string) (bool, error) {
	if path[1:] == ":/" {
		return false, nil
	}
    pointer, err := syscall.UTF16PtrFromString(path)
    if err != nil {
        return false, err
    }
    attributes, err := syscall.GetFileAttributes(pointer)
    if err != nil {
        return false, err
    }
    return attributes&syscall.FILE_ATTRIBUTE_HIDDEN != 0, nil
}