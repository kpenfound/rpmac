package cli

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/kpenfound/rpmac/constants"
	"github.com/kpenfound/rpmac/util"
)

// AddRepoCommand install
type AddRepoCommand struct{}

// Help text
func (i *AddRepoCommand) Help() string {
	return "rpmac add-repo {aboslute path to repo file}"
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
	filename := filepath.Base(args[0])
	repodir := util.ReplaceHome(constants.RepoDir)
	destfile := filepath.Join(repodir, filename)

	_ = os.Mkdir(repodir, 0755) // Make sure the repo dir exists

	err := util.Copy(args[0], destfile, 0644)
	if err != nil {
		fmt.Printf("Error adding repository: %s\n", err)
		return 1
	}
	return 0
}
