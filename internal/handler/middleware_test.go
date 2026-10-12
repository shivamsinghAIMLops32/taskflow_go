package handler

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestMaxBodySize(t *testing.T) {
	taskHandler := newTestHandler()

	mux := http.NewServeMux()
	mux.HandleFunc("POST /tasks", taskHandler.CreateTask)

	app := RequestLogger(
		MaxBodySize(16)(mux),
	)

	body := `{"title":"This request is too large"}`

	req := httptest.NewRequest(
		http.MethodPost,
		"/tasks",
		strings.NewReader(body),
	)

	rec := httptest.NewRecorder()

	app.ServeHTTP(rec, req)

	if rec.Code != http.StatusRequestEntityTooLarge {
		t.Fatalf(
			"expected status %d, got %d; body: %s",
			http.StatusRequestEntityTooLarge,
			rec.Code,
			rec.Body.String(),
		)
	}
}