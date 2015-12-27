package main

import (
	"io/ioutil"
	"path"

	"gopkg.in/yaml.v2"
)

func (app *_appContext) loadStates() {
	stateFile := path.Join(app.vimDir, "states.yml")
	data, err := ioutil.ReadFile(stateFile)
	if err != nil {
		return
	}

	if yaml.Unmarshal(data, &app.states) != nil {
		return
	}
}

func (app *_appContext) saveStates() {
	stateFile := path.Join(app.vimDir, "states.yml")
	data, err := yaml.Marshal(&app.states)
	if err != nil {
		return
	}
	ioutil.WriteFile(stateFile, data, 0644)
}

func (app *_appContext) getStringState(key string) string {
	if s, ok := app.states[key].(string); ok {
		return s
	}
	return ""
}

func (app *_appContext) getIntState(key string) int {
	if s, ok := app.states[key].(int); ok {
		return s
	}
	return 0
}

func (app *_appContext) getBoolState(key string) bool {
	if s, ok := app.states[key].(bool); ok {
		return s
	}
	return false
}

func (app *_appContext) setState(key string, value interface{}) {
	app.states[key] = value
}
