package log

import (
	"bytes"
	"fmt"
	"path/filepath"
	"runtime"

	"github.com/fatih/color"
)

type Level struct {
	priority  int
	name      string
	colorFunc func(...interface{}) string
}

var (
	gitVersion   string
	buildVersion string

	minPriority int

	level = struct {
		fatal  *Level
		error  *Level
		warn   *Level
		notice *Level
		info   *Level
		debug  *Level
	}{
		fatal:  &Level{5, "Fatal", color.New(color.BgRed).SprintFunc()},
		error:  &Level{4, "Error", color.New(color.FgHiRed).SprintFunc()},
		warn:   &Level{3, "Warn", color.New(color.FgHiYellow).SprintFunc()},
		notice: &Level{2, "Notice", color.New(color.FgHiGreen).SprintFunc()},
		info:   &Level{1, "Info", color.New(color.FgHiBlue).SprintFunc()},
		debug:  &Level{0, "Debug", color.New(color.FgHiBlack).SprintFunc()}}

	goTagColorFunc    = color.New(color.FgWhite).SprintFunc()
	goIdColorFunc     = color.New(color.FgWhite).SprintFunc()
	funcNameColorFunc = color.New(color.FgCyan).SprintFunc()
	fileColorFunc     = color.New(color.FgHiMagenta).SprintFunc()
	lineColorFunc     = color.New(color.FgYellow).SprintFunc()
)

func SetMinLevel(minLevel string) {
	switch minLevel {
	case level.fatal.name:
		minPriority = level.fatal.priority
		return
	case level.error.name:
		minPriority = level.error.priority
		return
	case level.warn.name:
		minPriority = level.warn.priority
		return
	case level.notice.name:
		minPriority = level.notice.priority
		return
	case level.info.name:
		minPriority = level.info.priority
		return
	case level.debug.name:
		minPriority = level.debug.priority
		return
	}

	minPriority = level.debug.priority
}

func Fatal(v ...interface{}) {
	printf(level.fatal, "", v...)
}

func Fatalf(format string, v ...interface{}) {
	printf(level.fatal, format, v...)
}

func Error(v ...interface{}) {
	printf(level.error, "", v...)
}

func Errorf(format string, v ...interface{}) {
	printf(level.error, format, v...)
}

func Warn(v ...interface{}) {
	printf(level.warn, "", v...)
}

func Warnf(format string, v ...interface{}) {
	printf(level.warn, format, v...)
}

func Notice(v ...interface{}) {
	printf(level.notice, "", v...)
}

func Noticef(format string, v ...interface{}) {
	printf(level.notice, format, v...)
}

func Info(v ...interface{}) {
	printf(level.info, "", v...)
}

func Infof(format string, v ...interface{}) {
	printf(level.info, format, v...)
}

func Debug(v ...interface{}) {
	printf(level.debug, "", v...)
}

func Debugf(format string, v ...interface{}) {
	printf(level.debug, format, v...)
}

func printf(level *Level, format string, v ...interface{}) {
	if level.priority < minPriority {
		return
	}

	buffer := make([]byte, 64)
	buffer = buffer[:runtime.Stack(buffer, false)]
	bufList := bytes.Fields(buffer)
	goId := "-1"
	if len(bufList) >= 2 {
		goId = fmt.Sprintf("%s", string(bufList[1]))
	}

	funcName := "???"
	pc, file, line, ok := runtime.Caller(2)
	if !ok {
		file = "???"
		line = -1
	} else {
		funcName = runtime.FuncForPC(pc).Name()
		file = filepath.Base(file)
	}

	msg := ""
	if len(format) != 0 {
		msg = fmt.Sprintf(format, v...)
	} else {
		if len(v) != 0 {
			msg = fmt.Sprint(v...)
		}
	}
	if len(msg) != 0 {
		msg = " " + msg
	}

	fmt.Printf(
		"%s %s%s:%s %s:%s:%s {git:%s, build:%s}\n",
		level.colorFunc("["+level.name+"]"),
		goTagColorFunc("Go"),
		goIdColorFunc(goId),
		msg,
		funcNameColorFunc(funcName+"()"),
		fileColorFunc(file),
		lineColorFunc(line),
		gitVersion,
		buildVersion)
}
