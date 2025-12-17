package logs

import (
	"fmt"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

type statusResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (s *statusResponseWriter) WriteHeader(code int) {
	s.statusCode = code
	s.ResponseWriter.WriteHeader(code)
}

func (s *statusResponseWriter) Write(b []byte) (int, error) {
	return s.ResponseWriter.Write(b)
}

func HttpLogger(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		sw := &statusResponseWriter{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}
		now := time.Now()
		defer func() {
			log(now, sw.statusCode, r)
		}()

		next.ServeHTTP(sw, r)
	}
	return http.HandlerFunc(fn)
}

func log(startAt time.Time, status int, r *http.Request) {
	entry := logrus.WithFields(logrus.Fields{
		"status":  status,
		"path":    r.URL.Path,
		"method":  r.Method,
		"latency": fmt.Sprintf("%.4f", time.Since(startAt).Seconds()),
	})
	switch {
	case status < 400:
		entry.Info()
	case status < 500:
		entry.Warning()
	default:
		entry.Error()
	}
}
