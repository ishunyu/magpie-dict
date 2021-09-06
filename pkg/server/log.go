package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	log "github.com/sirupsen/logrus"
)

const (
	timeFormat string = "2006-01-02 15:04:05.999 -0700"
)

var (
	requestLog *log.Logger
)

type Formatter struct {
}

func RequestLogger() *log.Logger {
	return requestLog
}

func setupLogger() {
	log.SetFormatter(new(Formatter))
	log.SetReportCaller(true)
}

func setupRequestLogger(config *Config) {
	if requestLog != nil {
		return
	}

	requestLog = log.New()

	path := filepath.Join(config.TempPath, "requests.log")
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Errorf("Problem trying to open %s for logging.", path)
	}
	requestLog.SetOutput(f)

	requestLog.SetFormatter(&log.JSONFormatter{
		TimestampFormat: timeFormat,
	})
}

func (f *Formatter) Format(entry *log.Entry) ([]byte, error) {
	time := entry.Time.Format(timeFormat)
	level := strings.ToUpper(entry.Level.String())[0]

	callerFile := ""
	callerLine := 0

	caller := entry.Caller
	if entry.Logger.ReportCaller && caller != nil {
		if caller.File != "" {
			callerFile = filepath.Base(caller.File)
		}
		callerLine = caller.Line
	}

	msg := fmt.Sprintf("%v [%c] %s:%d %s\n", time, level, callerFile, callerLine, entry.Message)

	return []byte(msg), nil
}
