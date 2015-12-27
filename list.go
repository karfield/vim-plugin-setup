package main

import (
	"fmt"
	"path"

	"github.com/codegangsta/cli"
	"github.com/ungerik/go-dry"
)

var listCommand = cli.Command{
	Name:    "list",
	Usage:   "list installed vim plugins",
	Aliases: []string{"ls"},
	Action: func(c *cli.Context) {
		vimdir := c.GlobalString("vimdir")
		fmt.Println("List plugins:")
		fl, err := dry.ListDirDirectories(path.Join(vimdir, "bundle"))
		if err != nil {
			fmt.Printf("cannot access to '%s' (error: %s)\n", vimdir, err)
			return
		}
		for _, plugin := range fl {
			fmt.Println(" ", plugin)
		}
	},
}
