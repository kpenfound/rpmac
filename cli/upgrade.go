package cli

import (
	"fmt"
	"strings"

	"github.com/kpenfound/rpmac/constants"
	"github.com/kpenfound/rpmac/repository"
)

// UpgradeCommand upgrade
type UpgradeCommand struct{}

// Help text
func (i *UpgradeCommand) Help() string {
	return "rpmac upgrade [package]"
}

// Name text
func (i *UpgradeCommand) Name() string {
	return "upgrade"
}

// Synopsis text
func (i *UpgradeCommand) Synopsis() string {
	return "Upgrade a specified package or all packages"
}

// Run operation
func (i *UpgradeCommand) Run(args []string) int {
	r, err := repository.InitRepositories()
	if err != nil {
		fmt.Printf("Error initializing repositories: %s\n", err)
		return 1
	}

	qoInstalled := repository.MakeQueryOptions(args[0])
	qoInstalled.Installed = constants.InstalledTrue
	rpmInstalled, err := r.Query(qoInstalled)
	if err != nil {
		fmt.Printf("Error querying for package: %s\n", err)
		return 1
	}

	qoAny := repository.MakeQueryOptions(args[0])
	qoAny.Installed = constants.InstalledAny
	rpmAny, err := r.Query(qoAny)
	if err != nil {
		fmt.Printf("Error querying for package: %s\n", err)
		return 1
	}

	v1 := rpmInstalled.Package.Version.Version
	v2 := rpmAny.Package.Version.Version
	if strings.Compare(v1, v2) != 0 && repository.CompatibleVersion(v1, v2, false) {
		fmt.Printf("Upgrading %s from %s to %s\n", args[0], v1, v2)

		// Uninstall Old Version
		err = rpmInstalled.Package.Uninstall()
		if err != nil {
			fmt.Printf("Error uninstalling package: %s\n", err)
			return 1
		}
		err = rpmInstalled.Repository.Save()
		if err != nil {
			fmt.Printf("Error saving package state: %s\n", err)
			return 1
		}

		// Install New Version
		err = rpmAny.Package.Install(rpmAny.Repository.BaseURL)
		if err != nil {
			fmt.Printf("Error installing package: %s\n", err)
			return 1
		}
		err = rpmInstalled.Repository.Save()
		if err != nil {
			fmt.Printf("Error saving package state: %s\n", err)
			return 1
		}

		fmt.Println("Completed upgrade successfully.")
	}

	return 0
}
