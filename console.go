package main

import (
	"fmt"
	"os"
	"path"
	"reflect"
	"runtime"
	"strconv"
	"strings"

	"github.com/fatih/color"
)

func prehandleArgs(a ...interface{}) (string, []interface{}, bool) {
	if len(a) == 0 {
		return "", nil, false
	}
	format := ""
	nargs := 0
	args := a
	if reflect.TypeOf(a[0]).Kind() == reflect.String {
		format = a[0].(string)
		if len(a) == 1 {
			return format, nil, true
		} else {
			nargs = strings.Count(format, "%")
		}
		args = a[1:]
	}

	for len(args) > nargs {
		format += " %+v"
		nargs++
	}
	return format, args, true

}

func (app *_appContext) println(a ...interface{}) {
	if app.verboseFlag {
		return
	}
	format, args, _ := prehandleArgs(a...)
	fmt.Fprintf(os.Stdout, format+"\n", args...)
}

func (app *_appContext) printf(a ...interface{}) {
	if app.verboseFlag {
		return
	}
	if format, args, ok := prehandleArgs(a...); ok {
		fmt.Fprintf(os.Stdout, format, args...)
	}
}

func codePosition() string {
	_, file, line, ok := runtime.Caller(2)
	if ok {
		p := strings.Index(file, "fulcrumos.com")
		if p >= 0 {
			file = file[p:]
		} else {
			file = path.Base(file)
		}
		return file + "@" + strconv.Itoa(line)
	} else {
		return ""
	}
}

var _ERROR_MSG = color.New(color.FgHiWhite, color.BgRed, color.Bold).SprintFunc()("ERROR:")
var _error = color.New(color.FgRed).SprintFunc()

func (app *_appContext) err(a ...interface{}) {
	if app.verboseFlag {
		return
	}
	if format, args, ok := prehandleArgs(a...); ok {
		str := fmt.Sprintf(" "+format+" ", args...)
		fmt.Fprintf(os.Stderr, _ERROR_MSG+" "+_error(str)+"\n")
	}
}

var _WARNING_MSG = color.New(color.FgBlack, color.BgYellow).SprintFunc()("WARNING:")
var _warning = color.New(color.FgHiYellow, color.Bold, color.Underline).SprintFunc()

func (app *_appContext) warn(a ...interface{}) {
	if app.verboseFlag {
		return
	}
	if format, args, ok := prehandleArgs(a...); ok {
		str := fmt.Sprintf(format, args...)
		fmt.Fprintf(os.Stdout, _WARNING_MSG+" "+_warning(str)+"\n")
	}
}

func (app *_appContext) warning(a ...interface{}) {
	app.warn(a...)
}

var _SUCCESS_MSG = color.New(color.FgWhite, color.BgGreen).SprintFunc()("SUCCESS:")
var _success = color.New(color.FgHiGreen, color.Bold).SprintFunc()

func (app *_appContext) success(a ...interface{}) {
	if app.verboseFlag {
		return
	}
	if format, args, ok := prehandleArgs(a...); ok {
		str := fmt.Sprintf(format, args...)
		fmt.Fprintf(os.Stdout, _SUCCESS_MSG+" "+_success(str)+"\n")
	}
}

var _FATAL_MSG = color.New(color.FgHiYellow, color.BgRed, color.Bold).SprintFunc()("FATAL:")
var _fatal = color.New(color.FgHiRed, color.Bold).SprintFunc()

func (app *_appContext) fatal(a ...interface{}) {
	if format, args, ok := prehandleArgs(a...); ok {
		str := fmt.Sprintf(format, args...)
		fmt.Fprintf(os.Stderr, _FATAL_MSG+" "+_fatal(str)+"\n")
		os.Exit(-1)
	}
}

var _DEBUG_MSG = color.New(color.FgHiCyan, color.Bold).SprintFunc()("DBG")
var _debug = color.New(color.FgCyan).SprintFunc()

func (app *_appContext) debug(a ...interface{}) {
	if app.verboseFlag || !app.enableDebug {
		return
	}
	if format, args, ok := prehandleArgs(a...); ok {
		str := fmt.Sprintf(format, args...)
		fmt.Fprintf(os.Stdout, _DEBUG_MSG+"("+codePosition()+"): "+_debug(str)+"\n")
	}
}

var _INFO_MSG = color.New(color.FgGreen, color.Bold).SprintFunc()("INFO:")
var _info = color.New(color.FgGreen).SprintFunc()

func (app *_appContext) info(a ...interface{}) {
	if app.verboseFlag {
		return
	}
	if format, args, ok := prehandleArgs(a...); ok {
		str := fmt.Sprintf(format, args...)
		fmt.Fprintf(os.Stdout, _INFO_MSG+" "+_info(str)+"\n")
	}
}
