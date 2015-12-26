package main

import (
"errors"
"strings"
"os"
"os/user"
"path"
  "github.com/scalingo/codegangsta-cli"
"github.com/ungerik/go-dry"
)

func main() {
_user, err:=user.Current()
if err != nil {
os.Exit(-1)
}

	app:=cli.NewApp()
app.Name = "vimplugin"
app.Version = "1.0.0"
app.Usage = "simple util to help to install/manage vim plugins"
app.Flags = []cli.Flag{
cli.StringFlag{
Name: "bundle-dir,d",
Usage: "set bundle directory",
Value: path.Join(_user.HomeDir, ".vim/bundle"),
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
for _, preq := range _PREREQUISITES {
exists := false
	paths:= strings.Split(os.Getenv("PATH"), ":")
for _, p := range paths {
	if dry.FileExists(path.Join(p, preq)) {
exists = true
}

}
if !exists {
return errors.New("missing '" + preq+"'")
}
}
return nil
}

