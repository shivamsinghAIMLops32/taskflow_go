package main

import (
	"fmt"
	"log"

	"taskflow/internal/repository"
	"taskflow/internal/service"
)

func main() {
	store := repository.NewMemory()

	taskService := service.NewTaskService(store)

	firstTask, err := taskService.CreateTask("Learn Go interfaces")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Created: %+v\n", firstTask)

	secondTask, err := taskService.CreateTask("   ")
	if err != nil {
		fmt.Printf("Expected validation error: %v\n", err)
	}

	_ = secondTask
}