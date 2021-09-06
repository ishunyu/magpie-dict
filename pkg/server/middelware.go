package main

import (
	"net/http"
	"time"

	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

type RequestStats map[string]interface{}
type HandlerFunc func(http.ResponseWriter, *http.Request, *RequestStats)

func RequestLogHandler(f HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		request_id := uuid.New().String()
		stats := &RequestStats{
			"request_id": request_id,
		}

		timeStart := time.Now()
		f(w, req, stats)
		duration := time.Since(timeStart)

		fields := log.Fields{
			"client_ip": req.RemoteAddr,
			"endpoint":  req.URL.Path,
			"duration":  duration.Milliseconds(),
		}

		var msg interface{} = nil
		for k, v := range *stats {
			if k == "msg" {
				msg = v
				continue
			}
			fields[k] = v
		}

		if msg == nil {
			RequestLogger().WithFields(fields).Info("")
		} else {
			RequestLogger().WithFields(fields).Info(msg)
		}
	}
}

func (stats *RequestStats) RequestID() string {
	return (*stats)["request_id"].(string)
}

func (stats *RequestStats) Add(key string, value interface{}) {
	(*stats)[key] = value
}

func (stats *RequestStats) Message(message string) {
	stats.Add("msg", message)
}
