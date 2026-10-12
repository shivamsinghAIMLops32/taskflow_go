package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"taskflow/internal/repository"
	"taskflow/internal/service"
)

func newTestHandler() *TaskHandler {
	store := repository.NewMemory()
	taskService := service.NewTaskService(store)

	return NewTaskHandler(taskService)
}

func TestCreateTask(t *testing.T) {
	handler := newTestHandler()

	body := []byte(`{"title":"Learn Go testing"}`)

	req := httptest.NewRequest(
		http.MethodPost,
		"/tasks",
		bytes.NewReader(body),
	)

	rec := httptest.NewRecorder()

	handler.CreateTask(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf(
			"expected status %d, got %d; body: %s",
			http.StatusCreated,
			rec.Code,
			rec.Body.String(),
		)
	}

	if contentType := rec.Header().Get("Content-Type"); contentType != "application/json" {
		t.Fatalf("expected application/json, got %q", contentType)
	}

	var response struct {
		ID    int64  `json:"id"`
		Title string `json:"title"`
	}

	if err := json.Unmarshal(rec.Body.Bytes(), &response); err != nil {
		t.Fatalf("decode response: %v", err)
	}

	if response.ID != 1 {
		t.Errorf("expected ID 1, got %d", response.ID)
	}

	if response.Title != "Learn Go testing" {
		t.Errorf("unexpected title: %q", response.Title)
	}
}

func TestCreateTaskRejectsEmptyTitle(t *testing.T) {
	handler := newTestHandler()

	req := httptest.NewRequest(
		http.MethodPost,
		"/tasks",
		bytes.NewBufferString(`{"title":"   "}`),
	)

	rec := httptest.NewRecorder()

	handler.CreateTask(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusBadRequest,
			rec.Code,
		)
	}

	var response ErrorResponse

	if err := json.Unmarshal(rec.Body.Bytes(), &response); err != nil {
		t.Fatalf("decode error response: %v", err)
	}

	if response.Error.Code != "VALIDATION_ERROR" {
		t.Errorf("unexpected error code: %q", response.Error.Code)
	}
}

func TestGetMissingTask(t *testing.T) {
	handler := newTestHandler()

	req := httptest.NewRequest(
		http.MethodGet,
		"/tasks/999",
		nil,
	)

	// PathValue is normally populated by ServeMux.
	// Because this test calls the handler directly, set it manually.
	req.SetPathValue("id", "999")

	rec := httptest.NewRecorder()

	handler.GetTask(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusNotFound,
			rec.Code,
		)
	}
}