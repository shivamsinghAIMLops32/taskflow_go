package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"taskflow/internal/repository"
	"taskflow/internal/service"
)

type TaskHandler struct {
	service *service.TaskService
}

func NewTaskHandler(service *service.TaskService) *TaskHandler {
	return &TaskHandler{
		service: service,
	}
}

func (h *TaskHandler) CreateTask(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Title string `json:"title"`
	}

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&input); err != nil {
		writeError(
			w,
			http.StatusBadRequest,
			"INVALID_JSON",
			"Request body must contain valid task JSON",
		)
		return
	}

	// Reject a second JSON value or other trailing data.
	if err := decoder.Decode(&struct{}{}); !errors.Is(err, io.EOF) {
		writeError(
			w,
			http.StatusBadRequest,
			"INVALID_JSON",
			"Request body must contain exactly one JSON object",
		)
		return
	}

createdTask, err := h.service.CreateTask(
	r.Context(),
	input.Title,
)
	if err != nil {
		if errors.Is(err, service.ErrEmptyTitle) {
			writeError(
				w,
				http.StatusBadRequest,
				"VALIDATION_ERROR",
				err.Error(),
			)
			return
		}

		writeError(
			w,
			http.StatusInternalServerError,
			"INTERNAL_ERROR",
			"An unexpected error occurred",
		)
		return
	}

	writeJSON(w, http.StatusCreated, createdTask)
}

func (h *TaskHandler) GetTasks(w http.ResponseWriter, r *http.Request) {
tasks, err := h.service.GetAllTasks(r.Context())
if err != nil {
	writeError(
		w,
		http.StatusInternalServerError,
		"INTERNAL_ERROR",
		"An unexpected error occurred",
	)
	return
}

writeJSON(w, http.StatusOK, tasks)
}

func (h *TaskHandler) GetTask(w http.ResponseWriter, r *http.Request) {
	id, err := parseTaskID(r)
	if err != nil {
		writeError(
			w,
			http.StatusBadRequest,
			"INVALID_TASK_ID",
			"Task ID must be a positive integer",
		)
		return
	}

	foundTask, err := h.service.GetTask(
	r.Context(),
	id,
)
	if err != nil {
		if errors.Is(err, repository.ErrTaskNotFound) {
			writeError(
				w,
				http.StatusNotFound,
				"TASK_NOT_FOUND",
				"Task does not exist",
			)
			return
		}

		writeError(
			w,
			http.StatusInternalServerError,
			"INTERNAL_ERROR",
			"An unexpected error occurred",
		)
		return
	}

	writeJSON(w, http.StatusOK, foundTask)
}

func (h *TaskHandler) DeleteTask(w http.ResponseWriter, r *http.Request) {
	id, err := parseTaskID(r)
	if err != nil {
		writeError(
			w,
			http.StatusBadRequest,
			"INVALID_TASK_ID",
			"Task ID must be a positive integer",
		)
		return
	}

	if err := h.service.DeleteTask(r.Context(),id); err != nil {
		if errors.Is(err, repository.ErrTaskNotFound) {
			writeError(
				w,
				http.StatusNotFound,
				"TASK_NOT_FOUND",
				"Task does not exist",
			)
			return
		}

		writeError(
			w,
			http.StatusInternalServerError,
			"INTERNAL_ERROR",
			"An unexpected error occurred",
		)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func parseTaskID(r *http.Request) (int64, error) {
	rawID := r.PathValue("id")

	id, err := strconv.ParseInt(rawID, 10, 64)
	if err != nil || id <= 0 {
		return 0, fmt.Errorf("invalid task ID")
	}

	return id, nil
}