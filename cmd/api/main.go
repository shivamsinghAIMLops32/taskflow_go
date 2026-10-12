package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
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

	app := handler.RequestLogger(
	handler.MaxBodySize(1 << 20)(mux),
)

server := &http.Server{
	Addr:              ":8080",
	Handler:           app,
	ReadHeaderTimeout: 5 * time.Second,
}

	// Listen for Ctrl+C or a termination signal.
	signalCtx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()

	serverErr := make(chan error, 1)

	go func() {
		log.Println("TaskFlow API listening on http://localhost:8080")
		serverErr <- server.ListenAndServe()
	}()

	select {
	case <-signalCtx.Done():
		log.Println("Shutdown signal received")

	case err := <-serverErr:
		if !errors.Is(err, http.ErrServerClosed) {
			log.Printf("HTTP server failed: %v", err)
		}
		return
	}

	shutdownCtx, cancel := context.WithTimeout(
		context.Background(),
		10*time.Second,
	)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Printf("Graceful shutdown failed: %v", err)

		if closeErr := server.Close(); closeErr != nil {
			log.Printf("Force close failed: %v", closeErr)
		}
	}

	log.Println("TaskFlow API stopped")
}