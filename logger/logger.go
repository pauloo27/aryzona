package logger

import (
	"log"
	"os"
	"runtime/debug"
	"time"

	"github.com/Pauloo27/aryzona/utils"
)

var (
	stdout = log.New(os.Stdout, "", 0)
	stderr = log.New(os.Stderr, "", 0)
)

// levels
var (
	DEBUG   = Level{"DEBUG", ColorBlue, false, false}
	SUCCESS = Level{"SUCCESS", ColorGreen, false, false}
	INFO    = Level{"INFO", ColorWhite, false, false}
	WARN    = Level{"WARN", ColorYellow, false, false}
	ERROR   = Level{"ERROR", ColorRed, true, false}
	FATAL   = Level{"FATAL", ColorRed, true, true}
)

func Log(level Level, msg string) {
	now := time.Now().Format("15:04:05")
	formattedMsg := utils.Fmt(
		"%s%s[%s @ %s]%s %s", level.Color, ColorBold, level.Name, now, ColorReset, msg,
	)
	if level.Error {
		stderr.Println(formattedMsg)
		debug.PrintStack()
	} else {
		stdout.Println(formattedMsg)
	}
	if level.Fatal {
		os.Exit(-1)
	}
}

func Logf(level Level, format string, v ...interface{}) {
	Log(level, utils.Fmt(format, v...))
}

func Debug(msg string) {
	Log(DEBUG, msg)
}

func Success(msg string) {
	Log(SUCCESS, msg)
}

func Info(msg string) {
	Log(INFO, msg)
}

func Warn(msg string) {
	Log(WARN, msg)
}

func Error(msg string) {
	Log(ERROR, msg)
}

func Fatal(msg string) {
	Log(FATAL, msg)
}

func Debugf(format string, v ...interface{}) {
	Logf(DEBUG, format, v...)
}

func Successf(format string, v ...interface{}) {
	Logf(SUCCESS, format, v...)
}

func Infof(format string, v ...interface{}) {
	Logf(INFO, format, v...)
}

func Warnf(format string, v ...interface{}) {
	Logf(WARN, format, v...)
}

func Errorf(format string, v ...interface{}) {
	Logf(ERROR, format, v...)
}

func Fatalf(format string, v ...interface{}) {
	Logf(FATAL, format, v...)
}

func HandleFatal(err error, msg string) {
	if err == nil {
		return
	}
	Fatalf("%s: %v", msg, err.Error())
}
