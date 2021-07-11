package metrics

import (
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	HttpRequestHistogram = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: "watchdog",
		Subsystem: "http",
		Name:      "request_duration_seconds",
		Help:      "The latency of the HTTP requests.",
	}, []string{"handler", "method"})
)

// Midleware metrics middleware
type Midleware struct {
}

func (middleware *Midleware) Wrap(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		defer func() {
			duration := time.Since(start).Seconds()
			HttpRequestHistogram.WithLabelValues(r.URL.Path, r.Method).Observe(duration)
		}()
		next.ServeHTTP(w, r)
	})
}
