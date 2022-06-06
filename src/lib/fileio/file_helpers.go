package fileio

import (
	"os"
)

// FileExists will simply tell you if the file can be found at the given path
func FileExists(fpath string) bool {
	if _, err := os.Stat(fpath); err != nil {
		return false
	}

	return true
}
