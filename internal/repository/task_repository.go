package repository

import (
	"context"

	"taskflow/internal/task"
)

type TaskRepository interface {
	Create(ctx context.Context, title string) (task.Task, error)
	GetByID(ctx context.Context, id int64) (task.Task, error)
	GetAll(ctx context.Context) ([]task.Task, error)
	Delete(ctx context.Context, id int64) error
}