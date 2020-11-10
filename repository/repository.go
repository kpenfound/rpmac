package repository

import (
	"bytes"
	"compress/gzip"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/kpenfound/rpmac/constants"
	"github.com/kpenfound/rpmac/rpm"
	"github.com/kpenfound/rpmac/util"
)

// Repository package repository struct
type Repository struct {
	Name       string
	BaseURL    string
	Enabled    bool
	Gpgcheck   bool
	Revision   int
	CacheFiles []string
	Packages   []rpm.RPM
}

// Sync a repository metadata with local cache
func (r *Repository) Sync() error {
	cachePath := filepath.Join(constants.CacheDir, r.Name)

	// Read repomd.xml
	repomdURL := fmt.Sprintf("%s/repodata/repomd.xml", r.BaseURL)
	resp, err := http.Get(repomdURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	rmd := RepoMd{}

	err = xml.Unmarshal(body, &rmd)
	if err != nil {
		return err
	}

	// Check update TODO:Compare against cached revision
	if rmd.Revision > 0 {
		err = r.ClearCache()
		if err != nil {
			return err
		}

		var cacheFiles []string

		for _, item := range rmd.Items {
			itemURL := fmt.Sprintf("%s/%s", r.BaseURL, item.Location.Href)
			err = util.Download(itemURL, cachePath, constants.CachePerm)
			if err != nil {
				return err
			}

			downloadedFile := filepath.Join(cachePath, item.Location.Href)
			cacheFiles = append(cacheFiles, downloadedFile)
		}

		r.CacheFiles = cacheFiles
	}

	return nil
}

// Query for a package by name in local cache
func (r *Repository) Query(name string) (rpm.RPM, error) {
	p := rpm.RPM{}
	return p, nil
}

// ClearCache clears the repo cache
func (r *Repository) ClearCache() error {
	cachePath := filepath.Join(constants.CacheDir, r.Name)
	err := os.RemoveAll(cachePath)
	if err != nil {
		return err
	}

	err = os.Mkdir(cachePath, constants.CachePerm)
	return err
}

// LoadCache loads packages from cache files
func (r *Repository) LoadCache() error {
	var p []rpm.RPM

	for _, f := range r.CacheFiles {
		if strings.HasSuffix(f, "-primary.xml.gz") { // Only read primary for now
			gzdat, err := ioutil.ReadFile(f)
			if err != nil {
				return err
			}

			reader := bytes.NewReader(gzdat)
			dat, err := gzip.NewReader(reader)
			if err != nil {
				return err
			}

			dats, err := ioutil.ReadAll(dat)
			if err != nil {
				return err
			}

			mf := MetadataFile{}
			err = xml.Unmarshal(dats, &mf)

			for _, rpm := range mf.PackageList {
				p = append(p, rpm)
			}
		}
	}

	r.Packages = p
	return nil
}

// ************************
// repomd.xml structs
// ************************

// RepoMdItemHref struct for the repomd.xml item hrefs
type RepoMdItemHref struct {
	Href string `xml:"href,attr"`
}

// RepoMdItem struct for the repomd.xml data items
type RepoMdItem struct {
	Type         string         `xml:"type,attr"`
	Checksum     string         `xml:"checksum"`
	OpenChecksum string         `xml:"open-checksum"`
	Location     RepoMdItemHref `xml:"location"`
	Timestamp    int            `xml:"timestamp"`
	Size         int            `xml:"size"`
	OpenSize     int            `xml:"open-size"`
}

// RepoMd struct for the repomd.xml
type RepoMd struct {
	Revision int          `xml:"revision"`
	Items    []RepoMdItem `xml:"Group>data"`
}

// MetadataFile struct for the -primary.xml.gz
type MetadataFile struct {
	Packages    string    `xml:"packages,attr"`
	PackageList []rpm.RPM `xml:"Group>package"`
}
