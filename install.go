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
	bundleDir := path.Join(c.GlobalString("vimdir"), "bundle")
	for _, plugin := range c.Args() {
		installPlugin(bundleDir, plugin)
	}
}

func installPlugin(bundleDir, url string) error {
	pluginName := path.Base(url)
	installDir := path.Join(bundleDir, pluginName)
	if strings.HasPrefix(url, "github.com/") {
		url = "https://" + url
	}

	var cmd *exec.Cmd
	if dry.FileExists() {
		cmd = exec.Command("git", "pull")
	} else {
		cmd = exec.Command("git", "clone", url, installDir)
	}
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		// cannot access to the git
		return err
	}
	return nil
}
