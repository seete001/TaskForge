package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"

	"taskforge/internal/queue"
	"taskforge/internal/task"
)

func main() {
	q := queue.NewQueue()

	mux := http.NewServeMux()

	mux.HandleFunc("POST /tasks", func(w http.ResponseWriter, r *http.Request) {
		var request struct {
			Name string `json:"name"`
		}

		err := json.NewDecoder(r.Body).Decode(&request)
		if err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		if request.Name == "" {
			http.Error(w, "Task name is required", http.StatusBadRequest)
			return
		}

		newTask := task.Task{
			ID:        uuid.NewString(),
			Name:      request.Name,
			Status:    task.StatusPending,
			CreatedAt: time.Now(),
		}

		q.Enqueue(newTask)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)

		json.NewEncoder(w).Encode(newTask)
	})

	mux.HandleFunc("GET /tasks", func(w http.ResponseWriter, r *http.Request) {
		tasks := q.GetAll()

		w.Header().Set("Content-Type", "application/json")

		json.NewEncoder(w).Encode(tasks)
	})

	mux.HandleFunc("GET /tasks/{id}", func(w http.ResponseWriter, r *http.Request) {

		id := r.PathValue("id")
		task, ok := q.GetById(id)

		if !ok {
			http.Error(w, "Task not found", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")

		json.NewEncoder(w).Encode(task)
	})

	log.Println("TaskForge API running on :8080")

	log.Fatal(http.ListenAndServe(":8080", mux))
}
