package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"os/user"
	"path"
	"regexp"
	"runtime"
	"strings"
	"text/template"

	"github.com/codegangsta/cli"
	"github.com/ungerik/go-dry"
)

const _VIMRC_TEMPLATE = `
"
" THIS FILE IS GENERATED BY '{{.CMDNAME}}'
" DO NOT MIDIFY VIM CONFIGS HERE!
"
" If you want to change configs, Change directory to {{.CONFIGDIR}}
" Then you could modify/add/remove isolated vim configs
"
" {{.CMDNAME}} is a tool to help you manage vim plugins
" You can install plugin like this:
"	{{.CMDNAME}} install <plugin-name> [<other-plugin-list>]
"
" The actual plugin manage is using 'pathogen', it will autoload the other
" plugins which stored in {{.BUNDLEDIR}}
"
" After install one plugin, you could add some vim config in {{.CONFIGDIR}}
"
" For more information about {{.CMDNAME}}, use '-h' to print more help information
"
`

const _PATHOGEN_CONFIG = `
" Config for github.com/tpope/vim-pathogen
" 
execute pathogen#infect()
syntax on
filetype plugin indent on

`

const _PATHOGEN_VIM_URL = "https://raw.githubusercontent.com/tpope/vim-pathogen/master/autoload/pathogen.vim"

var _PATHOGEN_C_PATTERN = regexp.MustCompile("^\\s*exec(?:ute|)\\s+pathogen#.*")

func (app *_appContext) setupVimPlugins(c *cli.Context) error {
	app.vimDir = c.GlobalString("vimdir")
	app.vimrcPath = c.GlobalString("vimrc")
	app.bundleDir = path.Join(app.vimDir, "bundle")
	app.autoloadDir = path.Join(app.vimDir, "autoload")
	app.configDir = path.Join(app.vimDir, "configs")
	app.tmpDir = path.Join(app.vimDir, "tmp")
	app.cmdName = path.Base(os.Args[0])

	os.MkdirAll(app.bundleDir, 0755)
	os.MkdirAll(app.autoloadDir, 0755)
	os.MkdirAll(app.configDir, 0755)
	os.RemoveAll(app.tmpDir)
	os.MkdirAll(app.tmpDir, 0755)
	defer os.RemoveAll(app.tmpDir)

	app.vimrcBuf = bytes.NewBuffer([]byte{})
	app.oldVimrcBuf = bytes.NewBuffer([]byte{})
	app.generatedVimrc = true

	if dry.FileExists(app.vimrcPath) {
		oldVimrc, err := os.Open(app.vimrcPath)
		if err != nil {
			return err
		}
		defer oldVimrc.Close()
		app.oldVimrcBuf.Reset()
		scanner := bufio.NewScanner(oldVimrc)
		generated := false
		for scanner.Scan() {
			l := scanner.Text()
			if _PATHOGEN_C_PATTERN.MatchString(l) {
				// comment the pathongen config line(s)
				app.oldVimrcBuf.WriteString("\" ")
			}
			if strings.HasPrefix(l, "\" THIS FILE IS GENERATED BY ") {
				generated = true
			}
			app.oldVimrcBuf.WriteString(l)
			app.oldVimrcBuf.WriteString("\n")
		}
		app.generatedVimrc = generated
	}

	tpl, _ := template.New("vimrc").Parse(_VIMRC_TEMPLATE)
	tpl.Execute(app.vimrcBuf, struct {
		CMDNAME, CONFIGDIR, BUNDLEDIR, AUTOLOADDIR, VIMRCFILE string
	}{
		CMDNAME:     app.cmdName,
		CONFIGDIR:   app.configDir,
		BUNDLEDIR:   app.bundleDir,
		AUTOLOADDIR: app.autoloadDir,
		VIMRCFILE:   app.vimrcPath,
	})

	pathogenVim := path.Join(app.autoloadDir, "pathogen.vim")
	if !dry.FileExists(pathogenVim) {
		if err := app.installPathogen(pathogenVim); err != nil {
			return err
		}
	}

	app.vimrcBuf.WriteString(_PATHOGEN_CONFIG)

	if !app.generatedVimrc {
		// save the user defined old vimrc into config-dir
		if app.oldVimrcBuf.Len() > 0 {
			oldVimrcFile := path.Join(app.configDir, "_old_config.vimrc")
			saveConfig(oldVimrcFile, app.oldVimrcBuf, true, true)
		}
		// save prebuilt-included vim configs except common.vimrc
		for _confPath, _func := range _bindata {
			fn := path.Base(_confPath)
			if fn == "common.vimrc" {
				continue
			}
			if asset, err := _func(); err != nil {
				continue
			} else {
				saveConfig(path.Join(app.configDir, fn), asset.bytes, false, false)
			}
		}
	} else {
		// save common and prebuilt-included vim configs
		for _confPath, _func := range _bindata {
			confPath := path.Join(app.configDir, path.Base(_confPath))
			if asset, err := _func(); err != nil {
				continue
			} else {
				saveConfig(confPath, asset.bytes, false, false)
			}
		}
	}

	return app.installPluginsByConfigs()
}

func saveConfig(_path string, _data interface{}, force, backup bool) bool {
	if dry.FileExists(_path) {
		if backup {
			os.Rename(_path, path.Join(path.Dir(_path), "."+path.Base(_path)))
		}
		if force {
			os.Remove(_path)
		} else {
			return true
		}
	}
	if data, ok := _data.(*bytes.Buffer); ok {
		var file *os.File
		var err error
		file, err = os.Create(_path)
		if err != nil {
			return false
		}
		defer file.Close()

		w := bufio.NewWriter(file)
		w.Reset(file)
		io.Copy(w, data)
		if err := w.Flush(); err != nil {
			return false
		}
		return true
	} else if data, ok := _data.([]byte); ok {
		return ioutil.WriteFile(_path, data, 0644) == nil
	} else if data, ok := _data.(string); ok {
		return ioutil.WriteFile(_path, []byte(data), 0644) == nil
	}
	return false
}

func (app *_appContext) installPathogen(installPath string) error {
	fmt.Println("Install pathogen")
	if resp, err := http.Get(_PATHOGEN_VIM_URL); err != nil {
		return err
	} else {
		defer resp.Body.Close()
		pathogen, err := os.Create(installPath)
		if err != nil {
			return err
		}
		defer pathogen.Close()
		if _, err := io.Copy(pathogen, resp.Body); err != nil {
			return err
		}
	}
	return nil
}

var _INSTALL_SCRIPT_BEGIN_PATTERN = regexp.MustCompile("\\s*\"\\s+@run\\-script\\s*(?:\\(([^\\)]*)\\)|)")
var _INSTALL_SCRIPT_END_PATTERN = regexp.MustCompile("\\s*\"\\s+@end\\-script")
var _INSTALL_SCRIPT_LINE_PATTERN = regexp.MustCompile("\\s*\"(.*)")
var _INSTALL_PLUGIN_PATTERN = regexp.MustCompile("\\s*\"\\s+@require(?:\\-plugin|)\\s*:\\s*(.*)")

func (app *_appContext) _writeVimSource(configfile string) {
	if u, err := user.Current(); err == nil {
		if strings.HasPrefix(configfile, u.HomeDir) {
			configfile = "~/" + strings.TrimLeft(configfile, u.HomeDir)
		}
	}
	sourcefrom := "so " + configfile + "\n"
	app.vimrcBuf.WriteString(sourcefrom)
}

func (app *_appContext) installPluginsByConfigs() error {
	fl, err := dry.ListDirFiles(app.configDir)
	if err != nil {
		return err
	}

	commonRc := path.Join(app.configDir, "common.vimrc")
	if dry.FileExists(commonRc) {
		app._writeVimSource(commonRc)
	} else {
		oldVimrc := path.Join(app.configDir, "_old_config.vimrc")
		if dry.FileExists(oldVimrc) {
			app._writeVimSource(oldVimrc)
		}
	}

	for _, f := range fl {
		if strings.HasPrefix(f, "_") {
			continue
		}

		configfile := path.Join(app.configDir, f)
		err := app.installPuginByConfig(configfile)
		if err != nil {
			continue
		}

		app._writeVimSource(configfile)
	}

	app.vimrcBuf.WriteString("\n")

	if saveConfig(app.vimrcPath, app.vimrcBuf, true, false) {
		return nil
	} else {
		return errors.New("fails to update .vimrc")
	}
}

func (app *_appContext) installPuginByConfig(configFilepath string) error {
	file, err := os.Open(configFilepath)
	if err != nil {
		return err
	}
	defer file.Close()

	installScript := bytes.NewBufferString("")
	plugins := []string{}

	scanner := bufio.NewScanner(file)

	scriptBegin := false
	scriptEnd := false

	for scanner.Scan() {
		line := scanner.Text()
		if _INSTALL_SCRIPT_BEGIN_PATTERN.MatchString(line) {
			scriptBegin = true
			scriptEnd = false
			continue
		}
		if _INSTALL_SCRIPT_END_PATTERN.MatchString(line) {
			scriptBegin = false
			scriptEnd = true
			continue
		}
		if _INSTALL_PLUGIN_PATTERN.MatchString(line) {
			ss := _INSTALL_PLUGIN_PATTERN.FindStringSubmatch(line)
			plugins = append(plugins, ss[1])
			continue
		}
		if scriptBegin && !scriptEnd {
			ss := _INSTALL_SCRIPT_LINE_PATTERN.FindStringSubmatch(line)
			if len(ss) > 0 {
				installScript.WriteString(ss[1])
				installScript.WriteString("\n")
			}
		}
	}

	for _, plugin := range plugins {
		fmt.Printf("install plugin: [%s]\n", plugin)
		if err := app.installPlugin(plugin); err != nil {
			return err
		}
	}

	if installScript.Len() > 0 {
		tmpfile, err := ioutil.TempFile(app.tmpDir, ".script-")
		if err != nil {
			return err
		}
		defer tmpfile.Close()
		if _, err := io.WriteString(tmpfile, installScript.String()); err != nil {
			return err
		}
		cmd := exec.Command("/bin/bash", tmpfile.Name())
		cmd.Env = append(os.Environ(),
			"HOST_OS="+runtime.GOOS,
			"HOST_ARCH="+runtime.GOARCH,
			"VIMDIR="+path.Dir(app.bundleDir),
			"VIMBUNDLEDIR="+app.bundleDir,
		)
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		println(installScript.String())
		if err := cmd.Run(); err != nil {
			return err
		}
	}

	return nil
}
