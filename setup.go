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

func setupVimPlugins(c *cli.Context) error {
	vimDir := c.GlobalString("vimdir")
	vimrcPath := c.GlobalString("vimrc")
	bundleDir := path.Join(vimDir, "bundle")
	autoloadDir := path.Join(vimDir, "autoload")
	configDir := path.Join(vimDir, "configs")

	os.MkdirAll(bundleDir, 0755)
	os.MkdirAll(autoloadDir, 0755)
	os.MkdirAll(configDir, 0755)

	vimrcBuf := bytes.NewBuffer([]byte{})
	oldVimrcBuf := bytes.NewBuffer([]byte{})
	generatedVimrc := true

	if dry.FileExists(vimrcPath) {
		oldVimrc, err := os.Open(vimrcPath)
		if err != nil {
			return err
		}
		defer oldVimrc.Close()
		oldVimrcBuf.Reset()
		scanner := bufio.NewScanner(oldVimrc)
		generated := false
		for scanner.Scan() {
			l := scanner.Text()
			if _PATHOGEN_C_PATTERN.MatchString(l) {
				// comment the pathongen config line(s)
				oldVimrcBuf.WriteString("\" ")
			}
			if strings.HasPrefix(l, "\" THIS FILE IS GENERATED BY ") {
				generated = true
			}
			oldVimrcBuf.WriteString(l)
			oldVimrcBuf.WriteString("\n")
		}
		generatedVimrc = generated
	}

	tpl, _ := template.New("vimrc").Parse(_VIMRC_TEMPLATE)
	tpl.Execute(vimrcBuf, struct {
		CMDNAME, CONFIGDIR, BUNDLEDIR, AUTOLOADDIR, VIMRCFILE string
	}{
		CMDNAME:     path.Base(os.Args[0]),
		CONFIGDIR:   configDir,
		BUNDLEDIR:   bundleDir,
		AUTOLOADDIR: autoloadDir,
		VIMRCFILE:   vimrcPath,
	})

	pathogenVim := path.Join(autoloadDir, "pathogen.vim")
	if !dry.FileExists(pathogenVim) {
		if err := installPathogen(pathogenVim); err != nil {
			return err
		}
	}

	vimrcBuf.WriteString(_PATHOGEN_CONFIG)

	if !generatedVimrc {
		// save the user defined old vimrc into config-dir
		if oldVimrcBuf.Len() > 0 {
			oldVimrcFile := path.Join(configDir, "_old_config.vimrc")
			saveConfig(oldVimrcFile, oldVimrcBuf, true, true)
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
				saveConfig(path.Join(configDir, fn), asset.bytes, false, false)
			}
		}
	} else {
		// save common and prebuilt-included vim configs
		for _confPath, _func := range _bindata {
			confPath := path.Join(configDir, path.Base(_confPath))
			if asset, err := _func(); err != nil {
				continue
			} else {
				saveConfig(confPath, asset.bytes, false, false)
			}
		}
	}

	return installPluginsByConfigs(vimrcBuf, vimrcPath, configDir, bundleDir)
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

func installPathogen(installPath string) error {
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
var _INSTALL_PLUGIN_PATTERN = regexp.MustCompile("\\s*\"\\s+@require(\\-plugin|)\\s*:\\s*(.*)")

func _writeVimSource(vimrcBuf *bytes.Buffer, configfile string) {
	if u, err := user.Current(); err == nil {
		if strings.HasPrefix(configfile, u.HomeDir) {
			configfile = "~/" + strings.TrimLeft(configfile, u.HomeDir)
		}
	}
	sourcefrom := "so " + configfile + "\n"
	vimrcBuf.WriteString(sourcefrom)
}

func installPluginsByConfigs(vimrcBuf *bytes.Buffer, vimrcPath, configDir, bundleDir string) error {
	fl, err := dry.ListDirFiles(configDir)
	if err != nil {
		return err
	}

	commonRc := path.Join(configDir, "common.vimrc")
	if dry.FileExists(commonRc) {
		_writeVimSource(vimrcBuf, commonRc)
	} else {
		oldVimrc := path.Join(configDir, "_old_config.vimrc")
		if dry.FileExists(oldVimrc) {
			_writeVimSource(vimrcBuf, oldVimrc)
		}
	}

	for _, f := range fl {
		if strings.HasPrefix(f, "_") {
			continue
		}

		configfile := path.Join(configDir, f)
		err := installPuginByConfig(vimrcBuf, bundleDir, configfile)
		if err != nil {
			continue
		}

		_writeVimSource(vimrcBuf, configfile)
	}

	vimrcBuf.WriteString("\n")

	if saveConfig(vimrcPath, vimrcBuf, true, false) {
		return nil
	} else {
		return errors.New("fails to update .vimrc")
	}
}

func installPuginByConfig(vimrcBuf *bytes.Buffer, bundleDir, configFilepath string) error {
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
		fmt.Println("install plugin: ", plugin)
		if err := installPlugin(bundleDir, plugin); err != nil {
			return err
		}
	}

	if installScript.Len() > 0 {
		tmpDir, err := ioutil.TempDir("", "")
		if err != nil {
			return err
		}
		tmpfile, err := ioutil.TempFile(tmpDir, ".vim-setup-script-")
		if err != nil {
			return err
		}
		defer tmpfile.Close()
		if _, err := io.WriteString(tmpfile, installScript.String()); err != nil {
			return err
		}
		scriptFilepath := path.Join(tmpDir, tmpfile.Name())

		cmd := exec.Command("/bin/sh", scriptFilepath)
		cmd.Env = append(os.Environ(),
			"HOST_OS="+runtime.GOOS,
			"HOST_ARCH="+runtime.GOARCH,
			"VIMDIR="+path.Dir(bundleDir),
			"VIMBUNDLEDIR="+bundleDir,
		)
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			return err
		}
	}

	return nil
}
