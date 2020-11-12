package cli

import (
	"fmt"
	"os"

	"github.com/jedib0t/go-pretty/table"
	"github.com/kpenfound/rpmac/constants"
	"github.com/kpenfound/rpmac/repository"
)

// InstallCommand install
type InstallCommand struct{}

// Help text
func (i *InstallCommand) Help() string {
	return "rpmac install {package} {package2=1.0}"
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

	toInstall := []*repository.RepoPackage{}

	for _, a := range args {

		qo := repository.MakeQueryOptions(a)
		qo.Installed = constants.InstalledFalse
		pack, err := r.Query(qo)
		if err != nil {
			fmt.Printf("Error querying for package: %s\n", err)
			return 1
		}
		toInstall = append(toInstall, pack)
	}

	fmt.Printf("Resolved Dependnecies...\n\n")

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	style := t.Style()
	style.Options = table.OptionsNoBordersAndSeparators
	t.SetStyle(*style)

	t.AppendHeader(table.Row{"Package", "Version", "Repository", "Size"})
	for _, i := range toInstall {
		t.AppendRow(table.Row{i.Package.Name, i.Package.Version.Version, i.Repository.Name, i.Package.Size.Installed})
	}
	t.Render()
	fmt.Printf("\n")

	// TODO: confirm install
	for _, rpm := range toInstall {
		err = rpm.Package.Install(rpm.Repository.BaseURL)
		if err != nil {
			fmt.Printf("Error installing package: %s\n", err)
			return 1
		}
		fmt.Printf("Installed %s\n", rpm.Package.Name)
		err = rpm.Repository.Save()
		if err != nil {
			fmt.Printf("Error saving package state: %s\n", err)
			return 1
		}
	}
	fmt.Println("Installation complete!")

	return 0
}
