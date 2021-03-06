package repository

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/kpenfound/rpmac/constants"
	"github.com/kpenfound/rpmac/rpm"
	"github.com/kpenfound/rpmac/util"
)

// Repository package repository struct
type Repository struct {
	ID         string
	Name       string
	BaseURL    string
	Enabled    bool
	Gpgcheck   bool
	Revision   int
	CacheFiles []string
	Packages   []*rpm.RPM
}

// ReadRepoFile returns a repository slice for a given repo file
func ReadRepoFile(repofile string) ([]Repository, error) {
	r := []Repository{} // A repo file can contain multiple repos

	file, err := os.Open(repofile)
	if err != nil {
		return r, err
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	var line string
	var repo Repository
	repoCounter := 0
	for {
		line, err = reader.ReadString('\n')
		if err != nil && err != io.EOF {
			break
		}
		line = strings.Trim(line, "\n ")

		// Process the line here.
		if len(line) > 1 {
			if strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]") {
				if repoCounter > 0 { // More than one repo in this file, make a new repo
					r = append(r, repo)
				}
				repo = Repository{}
				repoCounter++
				id := strings.Trim(line, "[]")
				repo.ID = id
				repo.CacheFiles = []string{}
				repo.Packages = []*rpm.RPM{}
				repo.Revision = 0
			} else if strings.Contains(line, "=") {
				lineParts := strings.Split(line, "=")
				switch lineParts[0] {
				case "name":
					repo.Name = lineParts[1]
				case "baseurl":
					repo.BaseURL = lineParts[1]
				case "enabled":
					repo.Enabled = lineParts[1] == "1"
				case "gpgcheck":
					repo.Gpgcheck = lineParts[1] == "1"
				}
			}
		}

		if err != nil {
			break
		}
	}
	if err != io.EOF {
		return r, err
	}
	r = append(r, repo)

	return r, nil
}

// Sync a repository metadata with local cache
func (r *Repository) Sync() error {
	cacheDir := util.ReplaceHome(constants.CacheDir)
	cachePath := filepath.Join(cacheDir, r.ID)

	err := r.Load()
	if err != nil {
		return err
	}

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

	// Check update
	if rmd.Revision > r.Revision {
		err = r.ClearCache()
		if err != nil {
			return err
		}

		var cacheFiles []string

		for _, item := range rmd.Items {
			itemURL := fmt.Sprintf("%s/%s", r.BaseURL, item.Location.Href)
			fname, err := util.Download(itemURL, cachePath, constants.CachePerm)
			if err != nil {
				return err
			}

			cacheFiles = append(cacheFiles, fname)
		}

		r.CacheFiles = cacheFiles
		r.Revision = rmd.Revision

		err = r.LoadCache()
		if err != nil {
			return err
		}

		return r.Save()
	}

	return nil
}

// CompatibleVersion checks version1 to be compatible with version2
func CompatibleVersion(version1 string, version2 string, fuzzy bool) bool {
	v1 := strings.Split(version1, ".")
	v2 := strings.Split(version2, ".")

	for i, v1part := range v1 {
		// 1.0.0 is compatible with 1 and 1.0 if fuzzy==true
		if (len(v2) - 1) < i {
			return fuzzy
		}
		// Convert version parts to int for comparison
		v1partInt, err := strconv.ParseInt(v1part, 10, 32)
		if err != nil {
			return false
		}
		v2partInt, err := strconv.ParseInt(v2[i], 10, 32)
		if err != nil {
			return false //TODO: Handle beta and rc
		}

		// Check if v1 > v2
		if v1partInt > v2partInt {
			return false
		}
	}
	return true
}

// Query for a package by name in local cache
func (r *Repository) Query(opts QueryOptions) (*rpm.RPM, error) {
	packageIndex := -1
	var e error

	for i, rpm := range r.Packages {
		if rpm.Name == opts.Name {
			// If we dont want an installed version and we find one, error
			if opts.Installed == constants.InstalledFalse && rpm.Installed {
				return rpm, constants.ErrorPackageInstalled
			}

			// Check version against query version
			if opts.FuzzyVersion != "" && !CompatibleVersion(rpm.Version.Version, opts.FuzzyVersion, true) {
				continue
			}

			// If we want an installed version and package is not installed, keep looking
			if opts.Installed == constants.InstalledTrue && !rpm.Installed {
				continue
			}

			// Return latest compatible version of package
			if packageIndex > -1 && CompatibleVersion(r.Packages[i].Version.Version, rpm.Version.Version, false) {
				packageIndex = i
			}

			if packageIndex == -1 {
				packageIndex = i
			}
		}
	}

	// https://media.giphy.com/media/IHOOMIiw5v9VS/giphy.gif
	if packageIndex == -1 {
		e = constants.ErrorPackageNotFound
		packageIndex = 0
	}
	return r.Packages[packageIndex], e
}

// ClearCache clears the repo cache
func (r *Repository) ClearCache() error {
	cacheDir := util.ReplaceHome(constants.CacheDir)
	cachePath := filepath.Join(cacheDir, r.ID)
	err := os.RemoveAll(cachePath)
	if err != nil {
		return err
	}

	err = os.Mkdir(cachePath, 0755)
	return err
}

// LoadCache loads packages from cache files
func (r *Repository) LoadCache() error {
	var p []*rpm.RPM

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

			p = mf.PackageList
		}
	}

	r.Packages = p
	return nil
}

// Save repo struct as xml to cache dir
func (r *Repository) Save() error {
	cacheDir := util.ReplaceHome(constants.CacheDir)
	saveFile := filepath.Join(cacheDir, r.ID, "repo.cached.xml")

	file, err := xml.Marshal(r)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(saveFile, file, 0644)
}

// Load repo struct from cache dir
func (r *Repository) Load() error {
	cacheDir := util.ReplaceHome(constants.CacheDir)
	saveFile := filepath.Join(cacheDir, r.ID, "repo.cached.xml")

	// Check if the savefile exists. return if it doesn't
	if _, err := os.Stat(saveFile); os.IsNotExist(err) {
		return nil
	}

	file, err := ioutil.ReadFile(saveFile)
	if err != nil {
		return err
	}

	loaded := Repository{}

	err = xml.Unmarshal(file, &loaded)
	if err != nil {
		return err
	}

	// Only overwrite these things. Other attributes should be refreshed from the .repo file
	r.Revision = loaded.Revision
	r.CacheFiles = loaded.CacheFiles
	r.Packages = loaded.Packages

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
	Items    []RepoMdItem `xml:"data"`
}

// MetadataFile struct for the -primary.xml.gz
type MetadataFile struct {
	Packages    string     `xml:"packages,attr"`
	PackageList []*rpm.RPM `xml:"package"`
}
