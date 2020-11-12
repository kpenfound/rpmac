package cli

import (
	"fmt"
	"os"

	"github.com/jedib0t/go-pretty/table"
	"github.com/kpenfound/rpmac/repository"
)

// SearchCommand search
type SearchCommand struct{}

// Help text
func (i *SearchCommand) Help() string {
	return "rpmac search {package}"
}

// Name text
func (i *SearchCommand) Name() string {
	return "search"
}

// Synopsis text
func (i *SearchCommand) Synopsis() string {
	return "Search for a specified package in repositories"
}

// Run operation
func (i *SearchCommand) Run(args []string) int {
	r, err := repository.InitRepositories()
	if err != nil {
		fmt.Printf("Error initializing repositories: %s\n", err)
		return 1
	}

	fmt.Printf("Searching for %s...\n\n", args[0])
	qo := repository.MakeQueryOptions(args[0])
	rpm, err := r.Query(qo)
	if err != nil {
		fmt.Printf("Error querying for package: %s\n", err)
		return 1
	}

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	style := t.Style()
	style.Options = table.OptionsNoBordersAndSeparators
	t.SetStyle(*style)

	t.AppendHeader(table.Row{"", "Package", "Version", "Repository", "Size", "Installed"})
	t.AppendRow(table.Row{
		"",
		rpm.Package.Name,
		rpm.Package.Version.Version,
		rpm.Repository.Name,
		rpm.Package.Size.Installed,
		rpm.Package.Installed})
	t.Render()
	fmt.Printf("\n")

	return 0
}
