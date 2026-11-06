package main

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"taskforge/internal/database"
	"taskforge/internal/queue"
	"taskforge/internal/repository"
	"taskforge/internal/task"
)

func main() {
	ctx := context.Background()

	// PostgreSQL
	db, err := database.NewPostgresPool(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	repo := repository.NewTaskRepository(db)

	// Redis
	client := queue.NewRedisClient()
	defer client.Close()

	redisQueue := queue.NewRedisQueue(client)

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

		// Save to PostgreSQL
		if err := repo.CreateTask(r.Context(), newTask); err != nil {
			log.Printf("Database error: %v", err)
			http.Error(w, "Failed to save task", http.StatusInternalServerError)
			return
		}
		// Publish to redis
		if err := redisQueue.Enqueue(r.Context(), newTask); err != nil {
			log.Printf("Failed to enqueue task: %v", err)
			http.Error(w, "Failed to enqueue task", http.StatusServiceUnavailable)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)

		if err := json.NewEncoder(w).Encode(newTask); err != nil {
			log.Printf("Response encoding error: %v", err)
		}

		log.Printf("Task created and queued:%s", newTask.ID)
	})

	mux.HandleFunc("GET /tasks", func(w http.ResponseWriter, r *http.Request) {
		tasks, err := repo.GetAllTasks(r.Context())

		if err != nil {
			log.Printf("Failed to retrieve tasks from repository: %v", err)
			http.Error(w, "Failed to retrieve tasks from repository", http.StatusServiceUnavailable)
			return
		}

		w.Header().Set("Content-Type", "application/json")

		if err := json.NewEncoder(w).Encode(tasks); err != nil {
			log.Printf("Response encoding error: %v", err)
		}
	})

	mux.HandleFunc("GET /tasks/{id}", func(w http.ResponseWriter, r *http.Request) {

		id := r.PathValue("id")
		task, err := repo.GetTaskByID(r.Context(), id)

		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				http.Error(w, "Task not found", http.StatusNotFound)
				return
			}
			log.Printf("Failed to retrieve task: %v", err)
			http.Error(w, "Failed to retrieve task", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")

		json.NewEncoder(w).Encode(task)
	})

	log.Println("TaskForge API running on :8080")

	log.Fatal(http.ListenAndServe(":8080", mux))

}
