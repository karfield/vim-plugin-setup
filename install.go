package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path"
	"regexp"
	"strconv"

	"github.com/codegangsta/cli"
	"github.com/fatih/color"
	"github.com/ungerik/go-dry"
)

var installCommand = cli.Command{
	Name:    "install",
	Usage:   "install vim plugin(s)",
	Aliases: []string{"i"},
	Action: func(c *cli.Context) {
		if len(c.Args()) == 0 {
			color.Yellow("Missing vim plugin")
			return
		}
		for _, plugin := range c.Args() {
			_app.installPlugin(plugin)
		}
	},
}

type vimPluginInfo struct {
	name               string   `json:"name"`
	normalizedName     string   `json:"normalized_name"`
	createDate         uint64   `json:"create_at"`
	author             string   `json:"author"`
	slug               string   `json:"slug"`
	tags               []string `json:"tags"`
	shortDesc          string   `json:"short_desc"`
	category           string   `json:"category"`
	keywords           string   `json:"keywords"`
	pluginManagerUsers int      `json:"plugin_manager_users"`
	updatedDate        uint64   `json:"updated_at"`

	githubRepoName        string `json:"github_repo_name"`
	githubHomepage        string `github_homepage`
	githubReadme          string `json:"github_readme_filename"`
	githubUrl             string `json:"github_url"`
	githubVimStars        int    `json:"github_vim_scripts_stars"`
	githubScriptBundles   int    `json:"github_vim_script_bundles"`
	githubStars           int    `json:"github_stars"`
	githubScriptsRepoName string `json:"github_vim_scripts_repo_name"`
	githubOwner           string `json:"github_owner"`
	githubRepoId          string `json:"github_repo_id"`
	githubShortDesc       string `json:"github_short_desc"`
	githubBundles         int    `json:"github_bundles"`
	githubAuthor          string `json:"github_author"`

	vimorgName      string `json:"vimorg_name"`
	vimorgType      string `json:"vimorg_type"`
	vimorgAuthor    string `json:"vimorg_user"`
	vimorgUrl       string `json:"vimorg_url"`
	vimorgShortDesc string `json:"vimorg_short_desc"`
	vimorgRating    int    `vimorg_rating`
	vimorgNumRaters int    `json:"vimorg_num_raters"`
	vimorgDownloads int    `json:"vimorg_downloads"`
}

type searchResult struct {
	totalResults   int             `json:"total_pages"`
	resultsPerPage int             `json:"results_per_page"`
	totalPages     int             `json:"total_pages"`
	plugins        []vimPluginInfo `json:"plugins"`
}

func (app *_appContext) searchVimPlugin(keyword string, pageIndex int) (*searchResult, bool) {
	queries := make(url.Values)
	queries.Add("q", keyword)
	queries.Add("page", strconv.Itoa(pageIndex))
	resp, err := http.Get("http://vimawesome.com/?" + queries.Encode())
	if err != nil {
		return nil, false
	}
	_data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, false
	}
	searchResult := searchResult{}
	if json.Unmarshal([]byte(_data), &searchResult) != nil {
		return nil, false
	}

	return &searchResult, true
}

var _GIT_HTTP_URL_PATTERN = regexp.MustCompile("(https?\\:\\/\\/|)(?:([^\\/]+\\.(?:com|org|io|net))\\/|)(.*)(\\.git|)")
var _GIT_SSH_PATTERN = regexp.MustCompile("([^@]+)@([^\\:]+)\\:(.*)\\.git")

func getPluginNameFromUrl(url string) (string, string) {
	pluginName := ""
	if ss := _GIT_HTTP_URL_PATTERN.FindStringSubmatch(url); len(ss) > 0 {
		pluginName = path.Base(ss[3])
		if ss[1] != "" {
			url = ss[1]
		} else {
			url = "https://"
		}
		url += ss[2] + "/" + ss[3] + ss[4]
	} else if ss := _GIT_SSH_PATTERN.FindStringSubmatch(url); len(ss) > 0 {
		pluginName = path.Base(ss[3])
	} else {
		return "", ""
	}
	return pluginName, url
}

func (app *_appContext) installPlugin(url string) error {
	gitflag := false

	var pluginName string
	pluginName, url = getPluginNameFromUrl(url)

	app.info("Install plugin:", pluginName)

	if pluginName == "" {
		if result, ok := app.searchVimPlugin(path.Base(url), 1); ok {
			if result.totalResults > 0 {
				app.info("Find %d result:", result.totalResults)
				for i, plugin := range result.plugins {
					app.info("[%d] plugin: %s(%s)", i, plugin.name, plugin.normalizedName)
					app.info("     description: %s", plugin.shortDesc)
					app.info("     author: %s", plugin.author)
					if plugin.githubUrl != "" {
						app.info("   github: %s", plugin.githubUrl)
						if i == 0 {
							pluginName = plugin.normalizedName
							url = plugin.githubUrl
							gitflag = true
						}
					} else {
						app.info("   vim.org: %s", plugin.vimorgUrl)
						if i == 0 {
							pluginName = plugin.normalizedName
							url = plugin.vimorgUrl
						}
					}
				}

				if result.totalPages > 1 {
					app.info("Press [space] to load more plugins")
					// FIXME
				}

				if result.totalResults > 1 {
					app.info("type the index number to decide which plugin to be installed?")
					// FXIME
				}
			} else {
				app.err("No plugin matches for: ", url)
				return errors.New("name error")
			}
		} else {
			app.err("Sorry! Cannot recognize the plugin url/name pattern")
			return errors.New("name error")
		}
	} else {
		gitflag = true
	}

	if app.getBoolState("plugin:" + pluginName) {
		app.info("%s has been installed.", pluginName)
		return nil
	}

	installDir := path.Join(app.bundleDir, pluginName)

	if gitflag {
		var cmd *exec.Cmd
		if dry.FileIsDir(path.Join(installDir, ".git")) {
			app.info("Updating", url)
			cmd = exec.Command("git", "pull")
			cmd.Dir = installDir
		} else {
			app.info("Cloning", url)
			os.RemoveAll(installDir)
			cmd = exec.Command("git", "clone", url, installDir)
		}
		cmd.Stdin = os.Stdin
		if app.enableDebug {
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
		}

		if err := cmd.Run(); err != nil {
			// cannot access to the git
			app.err("Unable to sync:", url)
			return err
		}

		submoduleFile := path.Join(installDir, ".gitmodules")
		if dry.FileExists(submoduleFile) {
			cmd = exec.Command("git", "submodule", "update", "--init", "--recursive")
			cmd.Dir = installDir
			cmd.Stdin = os.Stdin
			if app.enableDebug {
				cmd.Stdout = os.Stdout
				cmd.Stderr = os.Stderr
			}
			if err := cmd.Run(); err != nil {
				// cannot access to the git
				return err
			}
		}
	} else {
	}

	app.setState("plugin:"+pluginName, true)
	return nil
}
