package cli

import (
	"fmt"

	"github.com/kpenfound/rpmac/constants"
	"github.com/kpenfound/rpmac/repository"
)

// UninstallCommand uninstall
type UninstallCommand struct{}

// Help text
func (i *UninstallCommand) Help() string {
	return "rpmac uninstall {package}"
}

// Name text
func (i *UninstallCommand) Name() string {
	return "uninstall"
}

// Synopsis text
func (i *UninstallCommand) Synopsis() string {
	return "Uninstall a specified package"
}

// Run operation
func (i *UninstallCommand) Run(args []string) int {
	r, err := repository.InitRepositories()
	if err != nil {
		fmt.Printf("Error initializing repositories: %s\n", err)
		return 1
	}

	qo := repository.MakeQueryOptions(args[0])
	qo.Installed = constants.InstalledTrue
	rpm, err := r.Query(qo)
	if err != nil {
		fmt.Printf("Error querying for package: %s\n", err)
		return 1
	}
	fmt.Printf("Uninstalling %s/%s\n", rpm.Repository.ID, rpm.Package.Name)

	err = rpm.Package.Uninstall()
	if err != nil {
		fmt.Printf("Error uninstalling package: %s\n", err)
		return 1
	}
	fmt.Printf("Uninstalled %s\n", rpm.Package.Name)
	err = rpm.Repository.Save()
	if err != nil {
		fmt.Printf("Error saving package state: %s\n", err)
		return 1
	}

	return 0
}
