package util

import (
	"os"

	"github.com/kpenfound/rpmac/constants"
)

// Lock creates the lockfile. Errors if it already exists
func Lock() error {
	lf := lockfile()

	if _, err := os.Stat(lf); os.IsNotExist(err) {
		_, err = os.Create(lf)
		return err
	}

	return constants.ErrorLockfileExists
}

// Unlock removes the lockfile
func Unlock() error {
	lf := lockfile()

	return os.Remove(lf)
}

func lockfile() string {
	return ReplaceHome(constants.Lockfile)
}
