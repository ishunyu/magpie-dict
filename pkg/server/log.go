package main

import (
	"fmt"
	"path/filepath"
	"strings"

	log "github.com/sirupsen/logrus"
)

const (
	timeFormat string = "2006-01-02 15:04:05.999 -0700"
)

type Formatter struct {
}

func setupLogger() {
	log.SetFormatter(new(Formatter))
	log.SetReportCaller(true)
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
