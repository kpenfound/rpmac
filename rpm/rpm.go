package rpm

import (
	"fmt"
	"os"

	"github.com/kpenfound/rpmac/constants"
	"github.com/kpenfound/rpmac/util"
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

// Install installs the RPM to the system
func (r *RPM) Install(baseurl string) error {
	fmt.Printf("Installing %s\n", r.Name)

	// Download the package
	packageURL := fmt.Sprintf("%s/%s", baseurl, r.Location.Href)
	cacheDir := util.ReplaceHome(constants.PackageCacheDir)
	_ = os.Mkdir(cacheDir, 0755) // Make sure the cache dir exists

	fname, err := util.Download(packageURL, cacheDir, constants.CachePerm)
	if err != nil {
		return err
	}

	fmt.Printf("Downloaded to %s\n", fname)

	// Install the package

	return nil
}
