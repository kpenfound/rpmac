package cli

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
	return 0
}
