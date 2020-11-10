package cli

import (
	"fmt"

	"github.com/kpenfound/rpmac/repository"
)

// InstallCommand install
type InstallCommand struct{}

// Help text
func (i *InstallCommand) Help() string {
	return "rpmac install {package}"
}

// Name text
func (i *InstallCommand) Name() string {
	return "install"
}

// Synopsis text
func (i *InstallCommand) Synopsis() string {
	return "Install a specified package"
}

// Run operation
func (i *InstallCommand) Run(args []string) int {
	r, err := repository.InitRepositories()
	if err != nil {
		fmt.Printf("Error initializing repositories: %s\n", err)
		return 1
	}

	rpm, err := r.Query("rpmac-test")
	if err != nil {
		fmt.Printf("Error querying for package: %s\n", err)
		return 1
	}
	fmt.Printf("Found package '%s' in repository '%s'\n", rpm.Package.Name, rpm.Repository.Name)

	err = rpm.Package.Install()
	if err != nil {
		fmt.Printf("Error installing package: %s\n", err)
		return 1
	}

	return 0
}
