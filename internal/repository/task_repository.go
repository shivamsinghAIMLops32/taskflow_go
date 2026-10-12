package repository

import "taskflow/internal/task"

type TaskRepository interface {
	Create(title string) task.Task
	GetByID(id int64) (task.Task, error)
	GetAll() []task.Task
	Delete(id int64) error
}
