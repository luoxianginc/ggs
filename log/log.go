package log

import (
	"fmt"
	"log"
	"os"
	"path"
	"strings"
	"time"

	"ggs/conf"
)

const (
	debugLevel = 0
	infoLevel  = 1
	warnLevel  = 2
	errorLevel = 3
	fatalLevel = 4
)

const (
	printDebugLevel = "[debug] "
	printInfoLevel  = "[info ] "
	printWarnLevel  = "[warn ] "
	printErrorLevel = "[error] "
	printFatalLevel = "[fatal] "
)

type Logger struct {
	level      int
	baseLogger *log.Logger
	file       *os.File
}

var gLogger *Logger

func init() {
	var level int
	switch strings.ToLower(conf.Env.LogLevel) {
	case "debug":
		level = debugLevel
	case "info":
		level = infoLevel
	case "warn":
		level = warnLevel
	case "error":
		level = errorLevel
	case "fatal":
		level = fatalLevel
	default:
		panic("unknown logger level: " + conf.Env.LogLevel)
	}

	var baseLogger *log.Logger
	var file *os.File
	if conf.Env.LogPath != "" {
		now := time.Now()

		filename := fmt.Sprintf("%d%02d%02d_%02d_%02d_%02d.log",
			now.Year(),
			now.Month(),
			now.Day(),
			now.Hour(),
			now.Minute(),
			now.Second())

		file, err := os.Create(path.Join(conf.Env.LogPath, filename))
		if err != nil {
			panic("cannot create log file")
		}

		baseLogger = log.New(file, "", log.LstdFlags)
	} else {
		baseLogger = log.New(os.Stdout, "", log.LstdFlags)
	}

	gLogger = new(Logger)
	gLogger.level = level
	gLogger.baseLogger = baseLogger
	gLogger.file = file
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
	logger.printf(debugLevel, printDebugLevel, format, a...)
}

func (logger *Logger) Info(format string, a ...interface{}) {
	logger.printf(infoLevel, printInfoLevel, format, a...)
}

func (logger *Logger) Warn(format string, a ...interface{}) {
	logger.printf(warnLevel, printWarnLevel, format, a...)
}

func (logger *Logger) Error(format string, a ...interface{}) {
	logger.printf(errorLevel, printErrorLevel, format, a...)
}

func (logger *Logger) Fatal(format string, a ...interface{}) {
	logger.printf(fatalLevel, printFatalLevel, format, a...)
}

func Debug(format string, a ...interface{}) {
	gLogger.Debug(format, a...)
}

func Info(format string, a ...interface{}) {
	gLogger.Info(format, a...)
}

func Warn(format string, a ...interface{}) {
	gLogger.Warn(format, a...)
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
