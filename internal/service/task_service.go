package service

import (
	"errors"
	"strings"

	"taskflow/internal/repository"
	"taskflow/internal/task"
)

var ErrEmptyTitle = errors.New("task title cannot be empty")

type TaskService struct {
	repo repository.TaskRepository
}

func NewTaskService(repo repository.TaskRepository) *TaskService {
	return &TaskService{
		repo: repo,
	}
}

func (s *TaskService) CreateTask(title string) (task.Task, error) {
	title = strings.TrimSpace(title)

	if title == "" {
		return task.Task{}, ErrEmptyTitle
	}

	createdTask := s.repo.Create(title)

	return createdTask, nil
}

func (s *TaskService) GetTask(id int64) (task.Task, error) {
	if id <= 0 {
		return task.Task{}, errors.New("task ID must be positive")
	}

	return s.repo.GetByID(id)
}

func (s *TaskService) GetAllTasks() []task.Task {
	return s.repo.GetAll()
}

func (s *TaskService) DeleteTask(id int64) error {
	if id <= 0 {
		return errors.New("task ID must be positive")
	}

	return s.repo.Delete(id)
}