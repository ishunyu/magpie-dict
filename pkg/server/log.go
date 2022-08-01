package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	log "github.com/sirupsen/logrus"
)

const (
	FORMAT_TIME    string = "2006-01-02 15:04:05.000 -0700"
	KEY_REQUEST_ID string = "request_id"
	LOG_REQUESTS   string = "requests.log"
)

var (
	requestLog *log.Logger
)

type Formatter struct {
}

func RequestLogger() *log.Logger {
	return requestLog
}

func SetupLogger(loggingPath string) error {
	fmt.Println("Configuring logging")
	err := os.MkdirAll(loggingPath, 0700)
	if err != nil {
		return err
	}

	log.SetOutput(os.Stdout)
	log.SetFormatter(new(Formatter))
	log.SetReportCaller(true)

	SetupRequestLogger(loggingPath)
	log.Info("Logging configured")

	return nil
}

func SetupRequestLogger(loggingPath string) {
	if requestLog != nil {
		return
	}

	requestLog = log.New()

	path := filepath.Join(loggingPath, LOG_REQUESTS)
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		fmt.Printf("Problem trying to open %s for logging.\n", path)
		os.Exit(1)
	}
	requestLog.SetOutput(f)

	requestLog.SetFormatter(&log.JSONFormatter{
		TimestampFormat: FORMAT_TIME,
	})
}

func (f *Formatter) Format(entry *log.Entry) ([]byte, error) {
	time := entry.Time.Format(FORMAT_TIME)
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
	callerInfo := fmt.Sprintf("%s:%d", callerFile, callerLine)

	requestID := ""
	if val, ok := entry.Data[KEY_REQUEST_ID]; ok {
		requestID = fmt.Sprintf(" [%v]", val)
	}

	msg := fmt.Sprintf("%v [%c] %-18s %s%s\n", time, level, callerInfo, entry.Message, requestID)

	return []byte(msg), nil
}
