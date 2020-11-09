package cli

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
	return 0
}
