package util

import (
	"fmt"

	"github.com/mitchellh/go-homedir"
)

// ReplaceHome is a wrapper for homedir.Expand
func ReplaceHome(filepath string) string {
	expanded, err := homedir.Expand(filepath)
	if err != nil {
		fmt.Printf(err.Error())
		return filepath
	}
	return expanded
}
