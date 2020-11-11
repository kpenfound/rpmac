package cli

import (
	"fmt"
	"os"

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

	packages := os.Args[2:]
	fmt.Printf("Searching %s\n", packages[0])
	rpm, err := r.Query(packages[0])
	if err != nil {
		fmt.Printf("Error querying for package: %s\n", err)
		return 1
	}
	fmt.Printf("Found package %s/%s\n", rpm.Repository.ID, rpm.Package.Name)
	return 0
}
