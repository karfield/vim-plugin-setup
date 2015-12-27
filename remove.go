package main

import "github.com/codegangsta/cli"

var removeCommand = cli.Command{
	Name:    "remove",
	Usage:   "remove vim plugin(s)",
	Aliases: []string{"rm"},
	Action: func(c *cli.Context) {
	},
}
