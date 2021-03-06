package main

import (
	"bytes"
	"errors"
	"os"
	"os/user"
	"path"
	"strings"

	"github.com/codegangsta/cli"
	"github.com/fatih/color"
	"github.com/ungerik/go-dry"
)

type _appContext struct {
	cmdName        string
	vimDir         string
	vimrcPath      string
	bundleDir      string
	configDir      string
	autoloadDir    string
	tmpDir         string
	vimrcBuf       *bytes.Buffer
	oldVimrcBuf    *bytes.Buffer
	generatedVimrc bool
	verboseFlag    bool
	enableDebug    bool
	forceUpdate    bool
	states         map[string]interface{}
}

var _app *_appContext

func main() {
	_user, err := user.Current()
	if err != nil {
		os.Exit(-1)
	}

	app := cli.NewApp()
	app.Name = path.Base(os.Args[0])
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
		cli.BoolFlag{
			Name:  "force,f",
			Usage: "force to update",
		},
		cli.BoolFlag{
			Name:  "verbose",
			Usage: "print more console infomation",
		},
		cli.BoolFlag{
			Name:  "debug",
			Usage: "show debug information",
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

var _PREREQUISITES = []string{"bash", "git", "vim", "wget", "cmake"}

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

	_app = new(_appContext)
	_app.states = make(map[string]interface{})
	_app.verboseFlag = !c.GlobalBool("verbose")
	_app.enableDebug = c.GlobalBool("debug")
	_app.forceUpdate = c.GlobalBool("force")

	return _app.setupVimPlugins(c)
}
