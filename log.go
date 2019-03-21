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
	fieldColorFunc    = color.New(color.FgGreen).SprintFunc()
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

func sprint(valueList ...interface{}) string {
	translatedList := make([]interface{}, 0, len(valueList))
	translateField := func(f *Field) {
		translatedList = append(translatedList, f.Sprint()...)
	}
	for _, value := range valueList {
		switch v := value.(type) {
		case *Field:
			translateField(v)
		default:
			translatedList = append(translatedList, value)
		}
	}
	formedList := make([]interface{}, 0, len(translatedList)*2)
	for _, value := range translatedList {
		formedList = append(formedList, value, ", ")
	}
	end := len(formedList)
	if end != 0 {
		// Strip the last comma.
		end--
	}
	return fmt.Sprint(formedList[:end]...)
}

func sprintf(format string, valueList ...interface{}) string {
	if len(valueList) == 0 {
		return fmt.Sprintf(format, valueList...)
	}

	beginField := len(valueList)
Loop:
	for idx, value := range valueList {
		switch value.(type) {
		case *Field:
			beginField = idx
			break Loop
		}
	}
	newList := make([]interface{}, 0, len(valueList)-beginField+1)
	newList = append(newList, fmt.Sprintf(format, valueList[:beginField]...))
	newList = append(newList, valueList[beginField:]...)
	return sprint(newList...)
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

	// go build -tags 'debug'
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
	printf(logger.output, depth, level.fatal, sprint(v...))
}

func (logger *concreteLogger) Fataldf(
	depth int, format string, v ...interface{}) {
	printf(logger.output, depth, level.fatal, sprintf(format, v...))
}

func (logger *concreteLogger) Errord(depth int, v ...interface{}) {
	printf(logger.output, depth, level.error, sprint(v...))
}

func (logger *concreteLogger) Errordf(
	depth int, format string, v ...interface{}) {
	printf(logger.output, depth, level.error, sprintf(format, v...))
}

func (logger *concreteLogger) Warnd(depth int, v ...interface{}) {
	printf(logger.output, depth, level.warn, sprint(v...))
}

func (logger *concreteLogger) Warndf(
	depth int, format string, v ...interface{}) {
	printf(logger.output, depth, level.warn, sprintf(format, v...))
}

func (logger *concreteLogger) Noticed(depth int, v ...interface{}) {
	printf(logger.output, depth, level.notice, sprint(v...))
}

func (logger *concreteLogger) Noticedf(
	depth int, format string, v ...interface{}) {
	printf(logger.output, depth, level.notice, sprintf(format, v...))
}

func (logger *concreteLogger) Infod(depth int, v ...interface{}) {
	printf(logger.output, depth, level.info, sprint(v...))
}

func (logger *concreteLogger) Infodf(
	depth int, format string, v ...interface{}) {
	printf(logger.output, depth, level.info, sprintf(format, v...))
}

func (logger *concreteLogger) Debugd(depth int, v ...interface{}) {
	printf(logger.output, depth, level.debug, sprint(v...))
}

func (logger *concreteLogger) Debugdf(
	depth int, format string, v ...interface{}) {
	printf(logger.output, depth, level.debug, sprintf(format, v...))
}

type Field struct {
	values map[string]interface{}
	keys   []string
}

func NewField(values ...interface{}) *Field {
	f := new(Field)
	f.values = make(map[string]interface{})
	f.keys = make([]string, 0, len(values))
	for idx := 0; idx+1 < len(values); idx += 2 {
		key := values[idx]
		strKey, ok := key.(string)
		if !ok {
			strKey = fmt.Sprint(key)
		}
		f.Add(strKey, values[idx+1])
	}
	return f
}

func NewErrField(value interface{}) *Field {
	return NewField("err", value)
}

func (f *Field) Sprint() []interface{} {
	translatedList := make([]interface{}, 0, len(f.keys))
	for _, key := range f.keys {
		_, exist := f.values[key]
		if !exist {
			translatedList = append(translatedList, key)
			continue
		}
		translatedList = append(
			translatedList,
			fmt.Sprintf("%s(%v)", fieldColorFunc(key), f.values[key]))
	}
	return translatedList
}

func (f *Field) Add(key string, value interface{}) *Field {
	// Disallow overwrite.
	_, exist := f.values[key]
	if exist {
		return f
	}

	f.keys = append(f.keys, key)
	f.values[key] = value
	return f
}
