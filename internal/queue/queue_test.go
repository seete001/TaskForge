package queue

import (
	"fmt"
	"sync"
	"taskforge/internal/task"
	"testing"
)

func TestEnqueueAndDequeue(t *testing.T) {
	q := NewQueue()

	original := task.Task{
		ID:   "task-1",
		Name: "send-email",
	}

	q.Enqueue(original)

	result, ok := q.Dequeue()

	if !ok {
		t.Fatal("expected a task, got an empty queue")
	}

	if result.ID != original.ID {
		t.Errorf("expected %s, got %s, ", original.ID, result.ID)
	}
}

func TestDequeueEmptyQueue(t *testing.T) {
	q := NewQueue()

	_, ok := q.Dequeue()

	if ok {
		t.Fatal("expected empty queue")
	}
}

func TestConcurrentEnqueue(t *testing.T) {
	q := NewQueue()

	const workers = 100

	var wg sync.WaitGroup

	for i := 0; i < workers; i++ {
		wg.Add(1)

		go func(i int) {
			defer wg.Done()

			q.Enqueue(task.Task{
				ID:   fmt.Sprintf("task-%d", i),
				Name: "test-task",
			})
		}(i)
	}

	wg.Wait()

	if len(q.tasks) != workers {
		t.Fatalf("expected %d tasks, got %d", workers, len(q.tasks))
	}
}
