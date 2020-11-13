package constants

import "errors"

var (
	// ErrorPackageNotFound Package not found
	ErrorPackageNotFound = errors.New("Package not found")
	// ErrorLockfileExists for when the lockfile is already there
	ErrorLockfileExists = errors.New("Cannot obtain lock on rpmac")
	// ErrorPackageInstalled for when a version of the package is already installed
	ErrorPackageInstalled = errors.New("A version of this package is already installed")
	// ErrorPackageNotInstalled for when the package is supposed to be installed but isnt
	ErrorPackageNotInstalled = errors.New("Package is not installed")
	// ErrorDependencyNotFound for when you cant find all the package dependencies in repos
	ErrorDependencyNotFound = errors.New("Cannot resolve all dependencies")
)
