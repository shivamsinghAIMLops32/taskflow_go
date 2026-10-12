package main

import (
	"log"
	"net/http"
	"time"

	"taskflow/internal/handler"
	"taskflow/internal/repository"
	"taskflow/internal/service"
)

func main() {
	store := repository.NewMemory()
	taskService := service.NewTaskService(store)
	taskHandler := handler.NewTaskHandler(taskService)

	mux := http.NewServeMux()

	mux.HandleFunc("POST /tasks", taskHandler.CreateTask)
	mux.HandleFunc("GET /tasks", taskHandler.GetTasks)
	mux.HandleFunc("GET /tasks/{id}", taskHandler.GetTask)
	mux.HandleFunc("DELETE /tasks/{id}", taskHandler.DeleteTask)

	server := &http.Server{
		Addr:              ":8080",
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}

	log.Println("TaskFlow API listening on http://localhost:8080")

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
}
