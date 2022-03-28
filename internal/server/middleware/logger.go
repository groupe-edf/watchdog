package middleware

import (
	"net/http"
	"runtime/debug"
	"time"

	"github.com/groupe-edf/watchdog/pkg/logging"
)

// Entry log entry
type Entry struct {
	Host               string
	Latency            time.Duration
	Protocole          string
	ReceivedTime       time.Time
	Referer            string
	RemoteIP           string
	ResponseHeaderSize int64
	ResponseBodySize   int64
	RequestBodySize    int64
	RequestHeaderSize  int64
	RequestMethod      string
	RequestURL         string
	ServerIP           string
	Status             int
	UserAgent          string
}

// Logger middleware logs http requests
type Logger struct {
	Logger logging.Interface
}

// Wrap implements Middleware
func (middleware Logger) Wrap(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				middleware.Logger.WithFields(logging.Fields{
					"error": err,
					"trace": string(debug.Stack()),
				}).Error()
				w.Write(debug.Stack())
			}
		}()
		start := time.Now()
		entry := Entry{
			Host:          r.Host,
			Protocole:     r.Proto,
			ReceivedTime:  start,
			Referer:       r.Referer(),
			RemoteIP:      r.RemoteAddr,
			RequestMethod: r.Method,
			RequestURL:    r.URL.String(),
			UserAgent:     r.UserAgent(),
		}
		next.ServeHTTP(w, r)
		entry.Latency = time.Since(start)
		middleware.Logger.WithFields(logging.Fields{
			"latency":       entry.Latency,
			"method":        entry.RequestMethod,
			"protocole":     entry.Protocole,
			"received_time": entry.ReceivedTime,
			"referer":       entry.Referer,
			"remote_ip":     entry.RemoteIP,
			"server_ip":     entry.ServerIP,
			"status":        entry.Status,
			"url":           entry.RequestURL,
			"user_agent":    entry.UserAgent,
		}).Info()
	})
}
