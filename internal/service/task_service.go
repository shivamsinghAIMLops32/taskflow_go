package service

import (
	"context"
	"errors"
	"strings"

	"taskflow/internal/repository"
	"taskflow/internal/task"
)

var ErrEmptyTitle = errors.New("task title cannot be empty")
var ErrInvalidTaskID = errors.New("task ID must be positive")

type TaskService struct {
	repo repository.TaskRepository
}

func NewTaskService(repo repository.TaskRepository) *TaskService {
	return &TaskService{
		repo: repo,
	}
}

func (s *TaskService) CreateTask(
	ctx context.Context,
	title string,
) (task.Task, error) {
	if err := ctx.Err(); err != nil {
		return task.Task{}, err
	}

	title = strings.TrimSpace(title)

	if title == "" {
		return task.Task{}, ErrEmptyTitle
	}

	return s.repo.Create(ctx, title)
}

func (s *TaskService) GetTask(
	ctx context.Context,
	id int64,
) (task.Task, error) {
	if id <= 0 {
		return task.Task{}, ErrInvalidTaskID
	}

	return s.repo.GetByID(ctx, id)
}

func (s *TaskService) GetAllTasks(
	ctx context.Context,
) ([]task.Task, error) {
	return s.repo.GetAll(ctx)
}

func (s *TaskService) DeleteTask(
	ctx context.Context,
	id int64,
) error {
	if id <= 0 {
		return ErrInvalidTaskID
	}

	return s.repo.Delete(ctx, id)
}