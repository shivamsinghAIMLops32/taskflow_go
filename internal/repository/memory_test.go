package repository

import (
	"testing"
	"fmt"
	"sync"
)

func TestConcurrency(t *testing.T){
	store := NewMemory()

	var wg sync.WaitGroup

	total := 1000

	for i := range total{
  wg.Go(func() {
	store.Create(fmt.Sprintf("task %d",i))
  })
	}

	wg.Wait()

	tasks := store.GetAll()

	if len(tasks)!= total{
		fmt.Printf("expected %d tasks, and got %d",total,len(tasks))
	}

	seen := make(map[int64]bool, total)

	for _, task := range tasks {
		if seen[task.ID] {
			t.Fatalf("duplicate task ID: %d", task.ID)
		}

		seen[task.ID] = true
	}
}