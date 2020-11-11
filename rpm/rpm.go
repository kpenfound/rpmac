package rpm

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/cavaliercoder/go-rpm"
	"github.com/kpenfound/rpmac/constants"
	"github.com/kpenfound/rpmac/util"
	"github.com/sassoftware/go-rpmutils"
)

// Version Version information
type Version struct {
	Epoch   int    `xml:"epoch,attr"`
	Version string `xml:"ver,attr"`
	Rel     int    `xml:"rel,attr"`
}

// Checksum checksum information
type Checksum struct {
	Type     string `xml:"type,attr"`
	Pkgid    string `xml:"pkgid,attr"`
	Checksum string
}

// Size size information
type Size struct {
	Package   int `xml:"package,attr"`
	Installed int `xml:"installed,attr"`
	Archive   int `xml:"archive,attr"`
}

// Location href
type Location struct {
	Href string `xml:"href,attr"`
}

// Provides files provided by the RPM
type Provides struct {
	Name    string `xml:"name,attr"`
	Flags   string `xml:"flags,attr"`
	Epoch   int    `xml:"epoch,attr"`
	Version string `xml:"ver,attr"`
	Rel     int    `xml:"rel,attr"`
}

// Format metadata information
type Format struct {
	License   string   `xml:"rpm:license"`
	Vendor    string   `xml:"rpm:vendor"`
	Group     string   `xml:"rpm:group"`
	Buildhost string   `xml:"rpm:buildhost"`
	SourceRPM string   `xml:"rpm:sourcerpm"`
	Provides  Provides `xml:"rpm:provides"`
	File      string   `xml:"file"`
}

// RPM package type
type RPM struct {
	Name        string   `xml:"name"`
	Type        string   `xml:"type,attr"`
	Arch        string   `xml:"arch"`
	Version     Version  `xml:"version"`
	Checksum    Checksum `xml:"checksum"`
	Summary     string   `xml:"summary"`
	Description string   `xml:"description"`
	Packager    string   `xml:"packager"`
	URL         string   `xml:"url"`
	Size        Size     `xml:"size"`
	Location    Location `xml:"location"`
	Format      Format   `xml:"format"`
	Installed   bool
}

// Uninstall uninstalls a package from the system
func (r *RPM) Uninstall() error {
	if r.Installed {
		// Find our package
		cacheDir := util.ReplaceHome(constants.PackageCacheDir)
		packageFile := filepath.Join(cacheDir, r.Location.Href)
		pf, err := rpm.OpenPackageFile(packageFile)
		if err != nil {
			return err
		}
		// Remove files
		for _, fi := range pf.Files() {
			if _, err := os.Stat(fi.Name()); err == nil {
				err = os.Remove(fi.Name())
				if err != nil {
					return err
				}
			}
		}
		// Finally, remove package from cache
		err = os.Remove(packageFile)
		if err != nil {
			return err
		}

		r.Installed = false
		return nil
	}

	return errors.New("Package is not installed")
}

// Install installs the RPM to the system
func (r *RPM) Install(baseurl string) error {
	if r.Installed {
		return errors.New("Package is already installed")
	}

	// Download the package
	packageURL := fmt.Sprintf("%s/%s", baseurl, r.Location.Href)
	cacheDir := util.ReplaceHome(constants.PackageCacheDir)
	_ = os.Mkdir(cacheDir, 0755) // Make sure the cache dir exists

	fname, err := util.Download(packageURL, cacheDir, constants.CachePerm)
	if err != nil {
		return err
	}

	// Install the package
	pf, err := rpm.OpenPackageFile(fname)
	if err != nil {
		return err
	}

	// Check provided files for existing files
	files := pf.Files()
	packagedFiles := []string{}
	for _, fi := range files {
		if _, err := os.Stat(fi.Name()); err == nil {
			errstr := fmt.Sprintf("file %s already exists", fi.Name())
			return errors.New(errstr)
		}
		packagedFiles = append(packagedFiles, fi.Name())
	}

	err = extract(pf.Name(), fname, files)
	if err != nil {
		return err
	}

	r.Installed = true
	return nil
}

func extract(packagename string, filename string, files []rpm.FileInfo) error {

	// Create temp dir for extraction
	tmp, err := ioutil.TempDir("", packagename)
	if err != nil {
		return err
	}
	defer os.RemoveAll(tmp)

	f, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	r := bytes.NewReader(f)

	// Load and extract rpm to temp directory
	rrpm, err := rpmutils.ReadRpm(r)
	if err != nil {
		return err
	}
	err = rrpm.ExpandPayload(tmp)
	if err != nil {
		return err
	}

	// Move extracted files to filesystem
	for _, fi := range files {
		tmpFile := filepath.Join(tmp, fi.Name())
		err = os.Link(tmpFile, fi.Name())
		if err != nil {
			return err
		}
		err = os.Chmod(fi.Name(), fi.Mode())
		if err != nil {
			return err
		}
		// TODO: set user, group, and stuff
	}

	return nil
}
