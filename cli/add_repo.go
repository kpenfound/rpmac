package cli

// AddRepoCommand install
type AddRepoCommand struct{}

// Help text
func (i *AddRepoCommand) Help() string {
	return "rpmac add-repo {repo file}"
}

// Name text
func (i *AddRepoCommand) Name() string {
	return "add-repo"
}

// Synopsis text
func (i *AddRepoCommand) Synopsis() string {
	return "Add a repository to rpmac"
}

// Run operation
func (i *AddRepoCommand) Run(args []string) int {
	return 0
}
