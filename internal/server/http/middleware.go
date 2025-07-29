package internalhttp

import (
	"context"
	"encoding/json"
	"net"
	"net/http"
	"time"

	"github.com/tolikproh/banners-rotation/internal/logger"
	"github.com/tolikproh/banners-rotation/pkg/httperr"
)

type ResponseWriter struct {
	http.ResponseWriter
	StatusCode  int
	BytesLength int
	log         *logger.Logger
}

func newResponseWriter(w http.ResponseWriter, log *logger.Logger) *ResponseWriter {
	return &ResponseWriter{w, http.StatusOK, 0, log}
}

func (w *ResponseWriter) WriteHeader(statusCode int) {
	w.StatusCode = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

func (w *ResponseWriter) Write(data []byte) (int, error) {
	n, err := w.ResponseWriter.Write(data)
	if err != nil {
		w.log.Error("write response error", "error", err)
		return 0, err
	}
	w.BytesLength += n
	return n, nil
}

func (s *Server) loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()

		rw := newResponseWriter(w, s.log)
		next.ServeHTTP(rw, r)

		addr, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			addr = "unknown"
		}
		s.log.Info("http request",
			"address", addr,
			"start time", startTime.UTC(),
			"method", r.Method,
			"path", r.URL.Path,
			"proto", r.Proto,
			"status code", rw.StatusCode,
			"latency [ms]", time.Since(startTime).Microseconds(),
			"length", rw.BytesLength,
			"user agent", r.UserAgent())
	})
}

type HandlerFunc func(ctx context.Context, r *http.Request) (interface{}, error)

func (s *Server) serveHandler(h HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		response := new(http.Response)

		ctx, cansel := context.WithTimeout(context.Background(), 1*time.Minute)
		defer cansel()

		data, err := h(ctx, r)
		if err != nil {
			code := httperr.HTTPStatus(err)

			response.StatusCode = code
			r.Response = response

			errorHTTP := new(errorHTTP)
			errorHTTP.Method = r.Method
			errorHTTP.Path = r.RequestURI
			errorHTTP.Error = err.Error()

			s.log.Error("http error request",
				"method", errorHTTP.Method,
				"path", errorHTTP.Path,
				"status code", code,
				"error", errorHTTP.Error,
				"user agent", r.UserAgent())

			b, _ := json.Marshal(errorHTTP)
			w.WriteHeader(code)
			w.Write(b)

			return
		}

		response.StatusCode = http.StatusOK
		r.Response = response
		b, _ := json.Marshal(data)
		w.WriteHeader(http.StatusOK)
		w.Write(b)
	}
}
