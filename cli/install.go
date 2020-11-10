package cli

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
	// r, err := repository.InitRepositories()
	// if err != nil {
	// 	fmt.Printf("Error initializing repositories: %s", err)
	// 	return 1
	// }

	return 0
}
