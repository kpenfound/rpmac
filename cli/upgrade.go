package cli

import (
	"fmt"
	"os"
	"strings"

	"github.com/jedib0t/go-pretty/table"
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
		fmt.Printf("Upgrading %s...\n\n", args[0])

		t := table.NewWriter()
		t.SetOutputMirror(os.Stdout)
		style := t.Style()
		style.Options = table.OptionsNoBordersAndSeparators
		t.SetStyle(*style)

		t.AppendHeader(table.Row{"", "Package", "Version", "Repository", "Size"})
		t.AppendRows([]table.Row{
			{
				"uninstall",
				rpmInstalled.Package.Name,
				rpmInstalled.Package.Version.Version,
				rpmInstalled.Repository.Name,
				rpmInstalled.Package.Size.Installed,
			},
			{
				"install",
				rpmAny.Package.Name,
				rpmAny.Package.Version.Version,
				rpmAny.Repository.Name,
				rpmAny.Package.Size.Installed,
			},
		})
		t.Render()
		fmt.Printf("\n")

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
	} else {
		fmt.Printf("Package %s is already at latest version.\n", args[0])
	}

	return 0
}
