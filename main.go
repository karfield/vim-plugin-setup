package main

import (
	"errors"
	"os"
	"os/user"
	"path"
	"strings"

	"github.com/codegangsta/cli"
	"github.com/fatih/color"
	"github.com/ungerik/go-dry"
)

func main() {
	_user, err := user.Current()
	if err != nil {
		os.Exit(-1)
	}

	app := cli.NewApp()
	app.Name = "vimplugin"
	app.Version = "1.0.0"
	app.Usage = "simple util to help to install/manage vim plugins"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "vimdir,d",
			Usage: "change vim directory",
			Value: path.Join(_user.HomeDir, ".vim"),
		},
		cli.StringFlag{
			Name:  "vimrc, rc",
			Usage: "change .vimrc path",
			Value: path.Join(_user.HomeDir, ".vimrc"),
		},
	}
	app.Before = checkBeforeRun
	app.Commands = []cli.Command{
		installCommand,
		listCommand,
		removeCommand,
	}

	app.Run(os.Args)
}

var _PREREQUISITES = []string{"git", "vim", "wget"}

func checkBeforeRun(c *cli.Context) error {
	preqMissing := []string{}
	for _, preq := range _PREREQUISITES {
		exists := false
		paths := strings.Split(os.Getenv("PATH"), ":")
		for _, p := range paths {
			if dry.FileExists(path.Join(p, preq)) {
				exists = true
			}

		}
		if !exists {
			preqMissing = append(preqMissing, preq)
		}
	}
	if len(preqMissing) > 0 {
		color.Red("Missing prequisite(s): %+v", preqMissing)
		return errors.New("missing prequisites")
	}

	return setupVimPlugins(c)
}
