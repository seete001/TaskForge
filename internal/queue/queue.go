package queue

import (
	"sync"
	"taskforge/internal/task"
)

// Defining a struct called Queue for Task queue
type Queue struct {
	mu    sync.Mutex
	tasks []task.Task
}

// This will create a new Queue and returns the pointer to it
// Pointer because of shared memory
func NewQueue() *Queue {
	return &Queue{
		tasks: make([]task.Task, 0),
	}
}

// Enqueue method to push a task to the end of slice
func (q *Queue) Enqueue(t task.Task) {
	q.mu.Lock()
	defer q.mu.Unlock()

	q.tasks = append(q.tasks, t)
}

// Dequeue the task slice and return false and an empty task if there is no task
func (q *Queue) Dequeue() (task.Task, bool) {
	q.mu.Lock()
	defer q.mu.Unlock()

	if len(q.tasks) == 0 {
		return task.Task{}, false
	}

	t := q.tasks[0]
	q.tasks = q.tasks[1:]

	return t, true
}

func (q *Queue) GetAll() []task.Task {
	q.mu.Lock()
	defer q.mu.Unlock()

	result := make([]task.Task, len(q.tasks))
	copy(result, q.tasks)

	return result
}

func (q *Queue) GetById(id string) (task.Task, bool) {
	q.mu.Lock()
	defer q.mu.Unlock()

	for _, t := range q.tasks {
		if t.ID == id {
			return t, true
		}
	}

	return task.Task{}, false
}
