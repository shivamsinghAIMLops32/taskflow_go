package main

import (
	"fmt"
	"taskflow/internal/task"
	"time"
)

func main() {
	newTask := &task.Task{
		ID:        1,
		Title:     "Learn go backend development",
		Completed: false,
		CreatedAt: time.Now(),
	}

	fmt.Printf("%+v\n", *newTask)
}
