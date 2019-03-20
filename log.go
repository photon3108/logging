package logging

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"sync"

	"github.com/fatih/color"
)

var (
	defaultLogger struct {
		value Logger
		once  sync.Once
	}

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
	goIDColorFunc     = color.New(color.FgWhite).SprintFunc()
	funcNameColorFunc = color.New(color.FgCyan).SprintFunc()
	fileColorFunc     = color.New(color.FgHiMagenta).SprintFunc()
	lineColorFunc     = color.New(color.FgYellow).SprintFunc()
)

// DefaultLogger initializes defaultLogger.value and returns it.
func DefaultLogger() Logger {
	defaultLogger.once.Do(func() {
		var err error
		defaultLogger.value, err = NewLogger()
		if err != nil {
			panic(err)
		}
	})
	return defaultLogger.value
}

// SetMinLevel sets minPriority.
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

func printf(output io.Writer, depth int, level *Level, msg string) {
	if level.priority < minPriority {
		return
	}

	buffer := make([]byte, 64)
	buffer = buffer[:runtime.Stack(buffer, false)]
	bufList := bytes.Fields(buffer)
	goID := "-1"
	if len(bufList) >= 2 {
		goID = fmt.Sprintf("%s", string(bufList[1]))
	}

	funcName := "???"
	pc, file, line, ok := runtime.Caller(depth)
	if !ok {
		file = "???"
		line = -1
	} else {
		funcName = runtime.FuncForPC(pc).Name()
		file = filepath.Base(file)
	}

	if len(msg) != 0 {
		msg = " " + msg
	}

	fmt.Fprintf(
		output,
		"%s %s%s:%s %s:%s:%s {git:%s, build:%s}\n",
		level.colorFunc("["+level.name+"]"),
		goTagColorFunc("Go"),
		goIDColorFunc(goID),
		msg,
		funcNameColorFunc(funcName+"()"),
		fileColorFunc(file),
		lineColorFunc(line),
		gitVersion,
		buildVersion)
}

// Logger is an interface.
type Logger interface {
	Fatal(v ...interface{})
	Fatalf(format string, v ...interface{})
	Error(v ...interface{})
	Errorf(format string, v ...interface{})
	Warn(v ...interface{})
	Warnf(format string, v ...interface{})
	Notice(v ...interface{})
	Noticef(format string, v ...interface{})
	Info(v ...interface{})
	Infof(format string, v ...interface{})
	Debug(v ...interface{})
	Debugf(format string, v ...interface{})
}

// LoggerWithDepth is a logger you can specify the depth of callers.
type LoggerWithDepth interface {
	Fatald(depth int, v ...interface{})
	Fataldf(depth int, format string, v ...interface{})
	Errord(depth int, v ...interface{})
	Errordf(depth int, format string, v ...interface{})
	Warnd(depth int, v ...interface{})
	Warndf(depth int, format string, v ...interface{})
	Noticed(depth int, v ...interface{})
	Noticedf(depth int, format string, v ...interface{})
	Infod(depth int, v ...interface{})
	Infodf(depth int, format string, v ...interface{})
	Debugd(depth int, v ...interface{})
	Debugdf(depth int, format string, v ...interface{})
}

// Level represents the level of log.
type Level struct {
	priority  int
	name      string
	colorFunc func(...interface{}) string
}

// NewLogger creates a new logger.
func NewLogger() (Logger, error) {
	logger := new(concreteLogger)
	logger.output = os.Stdout
	return logger, nil
}

type concreteLogger struct {
	output io.Writer
}

func (logger *concreteLogger) SetOutput(output io.Writer) {
	logger.output = output
}

func (logger *concreteLogger) Fatal(v ...interface{}) {
	logger.Fatald(3, v...)
}

func (logger *concreteLogger) Fatalf(format string, v ...interface{}) {
	logger.Fataldf(3, format, v...)
}

func (logger *concreteLogger) Error(v ...interface{}) {
	logger.Errord(3, v...)
}

func (logger *concreteLogger) Errorf(format string, v ...interface{}) {
	logger.Errordf(3, format, v...)
}

func (logger *concreteLogger) Warn(v ...interface{}) {
	logger.Warnd(3, v...)
}

func (logger *concreteLogger) Warnf(format string, v ...interface{}) {
	logger.Warndf(3, format, v...)
}

func (logger *concreteLogger) Notice(v ...interface{}) {
	logger.Noticed(3, v...)
}

func (logger *concreteLogger) Noticef(format string, v ...interface{}) {
	logger.Noticedf(3, format, v...)
}

func (logger *concreteLogger) Info(v ...interface{}) {
	logger.Infod(3, v...)
}

func (logger *concreteLogger) Infof(format string, v ...interface{}) {
	logger.Infodf(3, format, v...)
}

func (logger *concreteLogger) Fatald(depth int, v ...interface{}) {
	printf(logger.output, depth, level.fatal, fmt.Sprint(v...))
}

func (logger *concreteLogger) Fataldf(
	depth int, format string, v ...interface{}) {
	printf(logger.output, depth, level.fatal, fmt.Sprintf(format, v...))
}

func (logger *concreteLogger) Errord(depth int, v ...interface{}) {
	printf(logger.output, depth, level.error, fmt.Sprint(v...))
}

func (logger *concreteLogger) Errordf(
	depth int, format string, v ...interface{}) {
	printf(logger.output, depth, level.error, fmt.Sprintf(format, v...))
}

func (logger *concreteLogger) Warnd(depth int, v ...interface{}) {
	printf(logger.output, depth, level.warn, fmt.Sprint(v...))
}

func (logger *concreteLogger) Warndf(
	depth int, format string, v ...interface{}) {
	printf(logger.output, depth, level.warn, fmt.Sprintf(format, v...))
}

func (logger *concreteLogger) Noticed(depth int, v ...interface{}) {
	printf(logger.output, depth, level.notice, fmt.Sprint(v...))
}

func (logger *concreteLogger) Noticedf(
	depth int, format string, v ...interface{}) {
	printf(logger.output, depth, level.notice, fmt.Sprintf(format, v...))
}

func (logger *concreteLogger) Infod(depth int, v ...interface{}) {
	printf(logger.output, depth, level.info, fmt.Sprint(v...))
}

func (logger *concreteLogger) Infodf(
	depth int, format string, v ...interface{}) {
	printf(logger.output, depth, level.info, fmt.Sprintf(format, v...))
}

func (logger *concreteLogger) Debugd(depth int, v ...interface{}) {
	printf(logger.output, depth, level.debug, fmt.Sprint(v...))
}

func (logger *concreteLogger) Debugdf(
	depth int, format string, v ...interface{}) {
	printf(logger.output, depth, level.debug, fmt.Sprintf(format, v...))
}
