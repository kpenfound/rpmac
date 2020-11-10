package util

import (
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// Download util function to download a file from a given URL to a destination directory
func Download(URL string, destination string, perm os.FileMode) (string, error) {
	URLparts := strings.Split(URL, "/")
	filename := URLparts[len(URLparts)-1]
	resp, err := http.Get(URL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	file := filepath.Join(destination, filename)
	err = ioutil.WriteFile(file, body, perm)
	return file, err
}
