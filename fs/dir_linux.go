//go:build linux

package fs

func IsHidden(path string) (bool, error) {
	return false, nil
}