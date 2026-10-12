package handler

import (
	"log"
	"net/http"
	"time"
)

type statusRecorder struct {
	http.ResponseWriter
	status int
	bytes  int
}

func (r *statusRecorder) WriteHeader(status int) {
	if r.status != 0 {
		return
	}

	r.status = status
	r.ResponseWriter.WriteHeader(status)
}

func (r *statusRecorder) Write(data []byte) (int, error) {
	if r.status == 0 {
		r.WriteHeader(http.StatusOK)
	}

	n, err := r.ResponseWriter.Write(data)
	r.bytes += n

	return n, err
}

func RequestLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(
		w http.ResponseWriter,
		r *http.Request,
	) {
		start := time.Now()

		recorder := &statusRecorder{
			ResponseWriter: w,
		}

		next.ServeHTTP(recorder, r)

		status := recorder.status
		if status == 0 {
			status = http.StatusOK
		}

		duration := time.Since(start)

		log.Printf(
			"method=%s path=%s status=%d bytes=%d duration=%s",
			r.Method,
			r.URL.Path,
			status,
			recorder.bytes,
			duration,
		)
	})
}

func MaxBodySize(maxBytes int64) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(
			w http.ResponseWriter,
			r *http.Request,
		) {
			r.Body = http.MaxBytesReader(
				w,
				r.Body,
				maxBytes,
			)

			next.ServeHTTP(w, r)
		})
	}
}