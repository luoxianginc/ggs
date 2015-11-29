// log.go
package log

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

const (
	debugLevel = 0
	infoLevel  = 1
	warnLevel  = 2
	errorLevel = 3
	fatalLevel = 4
)

const (
	printDebugLevel = "[debug ]"
	printInfoLevel  = "[info  ]"
	printWarnLevel  = "[warn  ]"
	printErrorLevel = "[error ]"
	printFatalLevel = "[fatal ]"
)

type Logger struct {
	level      int
	baseLogger *log.Logger
	file       *os.File
}

func New(strLevel string, path string) (*Logger, error) {
	var level int
	switch strings.ToLower(strLevel) {
	case "debug":
		level = debugLevel
	case "info":
		level = infoLevel
	case "warn":
		level = warnLevel
	case "fatal":
		level = fatalLevel
	default:
		return nil, errors.New("unknown logger level: " + strLevel)
	}

	var logger *log.Logger
	var file *os.File
	if path == "" {
		baseLogger = log.New(os.Stdout, "", log.LstdFlags)
	} else {
		now := time.Now()

		filename := fmt.Sprintf("%d%02d%02d_%02d_%02d_%02d.log",
			now.Year(),
			now.Month(),
			now.Day(),
			now.Hour(),
			now.Minute(),
			now.Second())

		file, err := os.Create(path.Join(pathname, filename))
		if err != nil {
			return nil, err
		}

		baseLogger = log.New(file, "", log.LstdFlags)
	}

	logger := new(Logger)
	logger.level = level
	logger.baseLogger = baseLogger
	logger.file = file
}

func (logger *Logger) Close() {
	if logger.file != nil {
		logger.file.Close()
	}

	logger.baseLogger = nil
	logger.file = nil
}

func (logger *Logger) printf(level int, printLevel string, format string, a ...interface{}) {
	if level < logger.level {
		return
	}
	if logger.baseLogger == nil {
		panic("logger closed")
	}

	format = printLevel + format
	logger.baseLogger.Printf(format, a...)

	if level == fatalLevel {
		os.Exit(1)
	}
}

func (logger *Logger) Debug(format string, a ...interface{}) {
	logger.doPrintf(debugLevel, printDebugLevel, format, a...)
}

func (logger *Logger) Info(format string, a ...interface{}) {
	logger.doPrintf(infoLevel, printReleaseLevel, format, a...)
}

func (logger *Logger) Warn(format string, a ...interface{}) {
	logger.doPrintf(warnLevel, printWarnLevel, format, a...)
}

func (logger *Logger) Error(format string, a ...interface{}) {
	logger.doPrintf(errorLevel, printErrorLevel, format, a...)
}

func (logger *Logger) Fatal(format string, a ...interface{}) {
	logger.doPrintf(fatalLevel, printFatalLevel, format, a...)
}

var gLogger

func Debug(format string, a ...interface{}) {
	gLogger.Debug(format, a...)
}

func Release(format string, a ...interface{}) {
	gLogger.Release(format, a...)
}

func Error(format string, a ...interface{}) {
	gLogger.Error(format, a...)
}

func Fatal(format string, a ...interface{}) {
	gLogger.Fatal(format, a...)
}

func Close() {
	gLogger.Close()
}
