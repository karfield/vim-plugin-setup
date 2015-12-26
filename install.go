package main


import (
		"github.com/scalingo/codegangsta-cli"
		"path"
		"os"
		"io"
"fmt"
"net/http"
"bufio"
		"github.com/fatih/color"
		"github.com/ungerik/go-dry"
       )

var _PATHOGEN_VIM_URL = "https://raw.githubusercontent.com/tpope/vim-pathogen/master/autoload/pathogen.vim"

var installCommand = cli.Command{
Name: "install",
	      Usage: "install vim plugin(s)",
	      Aliases: []string{"i"},
	      Action: installPlugin,
}

func installPlugin(c *cli.Context) {
	for _, plugin:=range c.Args() {
bundleDir := c.GlobalString("bundle-dir")
		   if !dry.FileIsDir(bundleDir) {
			   if !prepareBundleDir(bundleDir) {
				   color.Red("Unable to access bundle dir (%s).", bundleDir)
					   os.Exit(-1)
			   }
		   }

	   if tryInstallByGit(bundleDir, plugin) {
		   color.Green("plugin \"%s\" has been installed.", plugin)
			   continue
	   }

	   if url, isGit, ok:=searchPlugin(plugin);!ok {
		   color.Yellow("Cannot find the plugin: %s", plugin)
			   continue
	   } else {
result:=false
	       if isGit {

		       result= tryInstallByGit(bundleDir, url) 
	       } else {
		       result= downloadAndInstallPluginTarball(bundleDir, url) 
	       }

       if result {
	       color.Green("plugin \"%s\" has been installed.", plugin)
       } else {
	       color.Red("Failed to install plugin: \"%s\"", plugin)
       }
	   }

	}
}

func prepareBundleDir(bundleDir string) bool {
	if os.MkdirAll(bundleDir, 0755) != nil {
		return false
	}
	vimDir:= path.Dir(bundleDir)
	autoloadDir := path.Join(vimDir, "autoload")
	if os.MkdirAll(autoloadDir, 0755) != nil {
		return false
	}
	configDir := path.Join(vimDir, "configs")
	if os.MkdirAll(configDir, 0755) != nil {
		return false
	}

	if !dry.FileExists(path.Join(autoloadDir, "pathogen.vim")) {
		if resp, err := http.Get(_PATHOGEN_VIM_URL); err != nil {
			return false
		} else {
			defer resp.Body.Close()
			pathogen, err := os.Create(path.Join(autoloadDir, "pathogen.vim"))
			if err != nil {
				return false
			}
			defer pathogen.Close()
			if _, err:=io.Copy(pathogen, resp.Body); err != nil {
			return false
			}
		}
	}

	return generateVimrc(vimDir)
}

var _VIMRC_HEADER = `
execute pathogen#infect()
syntax on
filetype plugin indent on

`

func generateVimrc(vimDir string) bool {
	vimrcFile, err := os.Create(path.Join(vimDir, "vimrc.inc"))
	if err != nil {
		return false
	}
	defer vimrcFile.Close()
	writer := bufio.NewWriter(vimrcFile)
	
	fmt.Println(_VIMRC_HEADER)

	files, err := dry.ListDirFiles(path.Join(vimDir, "configs"))
	if err != nil {
	return false
	}
	for _, f := range files {
		fmt.Printf("source %s/configs/%s\n", vimDir, f)
	}

	return writer.Flush() == nil
}

func tryInstallByGit(bundleDir, plugin string) bool {
return false
}

func searchPlugin(plugin string) (string,bool, bool){
return "",true, false
}

func downloadAndInstallPluginTarball(bundleDir, plugin string) bool {
return false
}
