package cli

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
	return 0
}
