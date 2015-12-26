package main



import (
  "github.com/scalingo/codegangsta-cli"
)

var listCommand = cli.Command{
Name: "list",
Usage: "list installed vim plugins",
Aliases: []string{"ls"},
Action: func(c *cli.Context) {
},
}
