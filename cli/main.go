package cli

import (
	"log"
	"os"

	"github.com/kpenfound/rpmac/util"
	"github.com/mitchellh/cli"
)

// Main cli entrypoint
func Main(args []string) int {
	c := cli.NewCLI("rpmac", "0.0.1")
	c.Args = os.Args[1:]

	c.Commands = map[string]cli.CommandFactory{
		"add-repo": func() (cli.Command, error) {
			return &AddRepoCommand{}, nil
		},
		"install": func() (cli.Command, error) {
			return &InstallCommand{}, nil
		},
		"uninstall": func() (cli.Command, error) {
			return &UninstallCommand{}, nil
		},
		"search": func() (cli.Command, error) {
			return &SearchCommand{}, nil
		},
		"upgrade": func() (cli.Command, error) {
			return &UpgradeCommand{}, nil
		},
	}

	err := util.Lock()
	if err != nil {
		log.Println(err)
		return 1
	}

	exitStatus, err := c.Run()
	if err != nil {
		log.Println(err)
	}

	err = util.Unlock()
	if err != nil {
		log.Println(err)
	}

	return exitStatus
}
