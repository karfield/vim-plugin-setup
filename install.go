package main

import (
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/codegangsta/cli"
	"github.com/fatih/color"
	"github.com/ungerik/go-dry"
)

var installCommand = cli.Command{
	Name:    "install",
	Usage:   "install vim plugin(s)",
	Aliases: []string{"i"},
	Action:  doInstallPlugin,
}

func doInstallPlugin(c *cli.Context) {
	if len(c.Args()) == 0 {
		color.Yellow("Missing vim plugin")
		return
	}
	for _, plugin := range c.Args() {
		_app.installPlugin(plugin)
	}
}

func (app *_appContext) installPlugin(url string) error {
	pluginName := path.Base(url)
	installDir := path.Join(app.bundleDir, pluginName)
	if strings.HasPrefix(url, "github.com/") {
		url = "https://" + url
	}

	var cmd *exec.Cmd
	if dry.FileExists(installDir) {
		cmd = exec.Command("git", "pull")
	} else {
		cmd = exec.Command("git", "clone", url, installDir)
	}
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	if app.verboseFlag {
		cmd.Stdout = os.Stdout
	}

	submoduleFile := path.Join(installDir, ".gitmodules")
	if dry.FileExists(submoduleFile) {
		if err := cmd.Run(); err != nil {
			// cannot access to the git
			return err
		}
		cmd = exec.Command("git", "submodule", "update", "--init", "--recursive")
		cmd.Stdin = os.Stdin
		cmd.Stderr = os.Stderr
		if app.verboseFlag {
			cmd.Stdout = os.Stdout
		}
		if err := cmd.Run(); err != nil {
			// cannot access to the git
			return err
		}
	}
	return nil
}
