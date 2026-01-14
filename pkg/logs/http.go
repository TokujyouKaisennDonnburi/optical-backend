package logs

import (
	"bufio"
	"fmt"
	"net"
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

func UnwrapWriter(writer any) (http.ResponseWriter, bool) {
	statusWriter, ok := writer.(statusResponseWriter)
	if !ok {
		return nil, false
	}
	return statusWriter.ResponseWriter, true
}

func (rw *statusResponseWriter) Flush() {
	if f, ok := rw.ResponseWriter.(http.Flusher); ok {
		f.Flush()
	}
}

func (rw *statusResponseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	if h, ok := rw.ResponseWriter.(http.Hijacker); ok {
		return h.Hijack()
	}
	return nil, nil, http.ErrNotSupported
}

func (rw *statusResponseWriter) Push(target string, opts *http.PushOptions) error {
	if p, ok := rw.ResponseWriter.(http.Pusher); ok {
		return p.Push(target, opts)
	}
	return http.ErrNotSupported
}

func HttpLogger(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		sw := &statusResponseWriter{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}
		now := time.Now().UTC()
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
