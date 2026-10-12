package repository

import (
	"context"
	"errors"
	"sync"
	"time"

	"taskflow/internal/task"
)

var ErrTaskNotFound = errors.New("task not found")

type Memory struct {
	mu     sync.RWMutex
	tasks  map[int64]task.Task
	nextID int64
}

func NewMemory() *Memory {
	return &Memory{
		tasks:  make(map[int64]task.Task),
		nextID: 1,
	}
}

func (m *Memory) Create(
	ctx context.Context,
	title string,
) (task.Task, error) {
	if err := ctx.Err(); err != nil {
		return task.Task{}, err
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	// Check again after acquiring the lock.
	if err := ctx.Err(); err != nil {
		return task.Task{}, err
	}

	newTask := task.Task{
		ID:        m.nextID,
		Title:     title,
		Completed: false,
		CreatedAt: time.Now().UTC(),
	}

	m.tasks[newTask.ID] = newTask
	m.nextID++

	return newTask, nil
}

func (m *Memory) GetByID(
	ctx context.Context,
	id int64,
) (task.Task, error) {
	if err := ctx.Err(); err != nil {
		return task.Task{}, err
	}

	m.mu.RLock()
	defer m.mu.RUnlock()

	if err := ctx.Err(); err != nil {
		return task.Task{}, err
	}

	foundTask, exists := m.tasks[id]
	if !exists {
		return task.Task{}, ErrTaskNotFound
	}

	return foundTask, nil
}

func (m *Memory) GetAll(
	ctx context.Context,
) ([]task.Task, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	m.mu.RLock()
	defer m.mu.RUnlock()

	if err := ctx.Err(); err != nil {
		return nil, err
	}

	tasks := make([]task.Task, 0, len(m.tasks))

	for _, t := range m.tasks {
		tasks = append(tasks, t)
	}

	return tasks, nil
}

func (m *Memory) Delete(
	ctx context.Context,
	id int64,
) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	if err := ctx.Err(); err != nil {
		return err
	}

	if _, exists := m.tasks[id]; !exists {
		return ErrTaskNotFound
	}

	delete(m.tasks, id)

	return nil
}

var _ TaskRepository = (*Memory)(nil)