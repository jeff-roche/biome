package fileio

import (
	"os"
	"os/user"
)

// GetHomeDir will get the path to the current users home directory
func GetHomeDir() (dir string, err error) {
	var cu *user.User

	if cu, err = user.Current(); err != nil {
		dir = ""
	} else {
		dir = cu.HomeDir
	}

	return
}

// GetCD will get the path to the current directory
// this is just a wrapper around os.Getwd()
func GetCD() (string, error) {
	return os.Getwd()
}
