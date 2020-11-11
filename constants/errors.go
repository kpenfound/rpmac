package constants

import "errors"

var (
	// ErrorPackageNotFound Package not found
	ErrorPackageNotFound = errors.New("Package not found")
	// ErrorLockfileExists for when the lockfile is already there
	ErrorLockfileExists = errors.New("Cannot obtain lock on rpmac")
)
