//go:build windows

package fs

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsHidden(t *testing.T) {
	fmt.Println(" TestIsHidden")

	flag, err := IsHidden("C:/")
	assert.True(t, err == nil && !flag)

	flag, err = IsHidden("C:/D")
	assert.True(t, err == nil && !flag)

	flag, err = IsHidden("D:/$RECYCLE.BIN")
	assert.True(t, err == nil && flag)
}