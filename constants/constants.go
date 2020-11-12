package constants

const (
	// CacheDir Directory for package cache
	CacheDir = "~/.rpmac/cache"
	// CachePerm is the permissions for cache objects
	CachePerm = 0644
	// RepoDir Directory for repo files
	RepoDir = "~/.rpmac/repos.d"
	// PackageCacheDir to cache rpm files
	PackageCacheDir = "~/.rpmac/packagecache"
	// Lockfile for locking during operations
	Lockfile = "~/.rpmac/.LOCK"
)

// Installed Opts
const (
	InstalledTrue  = iota
	InstalledFalse = iota
	InstalledAny   = iota
)
