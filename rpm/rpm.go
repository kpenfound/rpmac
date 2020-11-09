package rpm

type RPMVersion struct {
	Epoch   int    `xml:"epoch,attr"`
	Version string `xml:"ver,attr"`
	Rel     int    `xml:"rel,attr"`
}

type RPMChecksum struct {
	Type     string `xml:"type,attr"`
	Pkgid    string `xml:"pkgid,attr"`
	Checksum string
}

type RPMSize struct {
	Package   int `xml:"package,attr"`
	Installed int `xml:"installed,attr"`
	Archive   int `xml:"archive,attr"`
}

type RPMLocation struct {
	Href string `xml:"href,attr"`
}

type RPMProvides struct {
	Name    string `xml:"name,attr"`
	Flags   string `xml:"flags,attr"`
	Epoch   int    `xml:"epoch,attr"`
	Version string `xml:"ver,attr"`
	Rel     int    `xml:"rel,attr"`
}

type RPMFormat struct {
	License   string      `xml:"rpm:license"`
	Vendor    string      `xml:"rpm:vendor"`
	Group     string      `xml:"rpm:group"`
	Buildhost string      `xml:"rpm:buildhost"`
	SourceRPM string      `xml:"rpm:sourcerpm"`
	Provides  RPMProvides `xml:"rpm:provides"`
	File      string      `xml:"file"`
}

// RPM package type
type RPM struct {
	Name        string      `xml:"name"`
	Type        string      `xml:"type,attr"`
	Arch        string      `xml:"arch"`
	Version     RPMVersion  `xml:"version"`
	Checksum    RPMChecksum `xml:"checksum"`
	Summary     string      `xml:"summary"`
	Description string      `xml:"description"`
	Packager    string      `xml:"packager"`
	URL         string      `xml:"url"`
	Size        RPMSize     `xml:"size"`
	Location    RPMLocation `xml:"location"`
	Format      RPMFormat   `xml:"format"`
	Installed   bool
}
