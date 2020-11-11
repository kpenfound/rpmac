package util

import (
	"io/ioutil"
	"os"
)

// Copy a src file to a dst file
func Copy(src string, dst string, perm os.FileMode) error {
	b, err := ioutil.ReadFile(src)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(dst, b, perm)
}
