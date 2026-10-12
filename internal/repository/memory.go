package repository

import (
	"sync"
	"taskflow/internal/task"
	"time"
	"errors"
)

var ErrTaskNotFound = errors.New("task not found")

type Memory struct {
	mu sync.RWMutex
	tasks map[int64]task.Task
	nextID int64
}


func NewMemory() *Memory{
 return &Memory{
	  tasks: make(map[int64]task.Task),
	  nextID: 1,
 }
}

func (m *Memory) Create(title string) task.Task{
 m.mu.Lock()
 defer m.mu.Unlock()
 	newTask := task.Task{
		ID:        m.nextID,
		Title:     title,
		Completed: false,
		CreatedAt: time.Now().UTC(),
	}
	m.tasks[newTask.ID] = newTask
	m.nextID++

	return newTask
}

func (m *Memory) GetByID(id int64) (task.Task, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	foundTask, exists := m.tasks[id]
	if !exists {
		return task.Task{}, ErrTaskNotFound
	}

	return foundTask, nil
}
func (m *Memory) GetAll() []task.Task {
	m.mu.RLock()
	defer m.mu.RUnlock()

	tasks := make([]task.Task, 0, len(m.tasks))

	for _, t := range m.tasks {
		tasks = append(tasks, t)
	}

	return tasks
}

func (m *Memory) Delete(id int64) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.tasks[id]; !exists {
		return ErrTaskNotFound
	}

	delete(m.tasks, id)

	return nil
}