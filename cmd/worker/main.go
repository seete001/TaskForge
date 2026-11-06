package main

import (
	"context"
	"log"
	"time"

	"taskforge/internal/database"
	"taskforge/internal/queue"
	"taskforge/internal/repository"
	"taskforge/internal/task"
)

func ProcessTask(t task.Task) error {
	log.Printf("Processing task: %s | ID:%s", t.Name, t.ID)

	time.Sleep(2 * time.Second)

	log.Printf("Task completed: %s", t.ID)

	return nil
}

func main() {
	ctx := context.Background()

	// Connect to PostgreSQL
	db, err := database.NewPostgresPool(ctx)
	if err != nil {
		log.Fatal("PostgreSQL connection failed: ", err)
	}

	defer db.Close()

	repo := repository.NewTaskRepository(db)

	// Connect to redis
	client := queue.NewRedisClient()
	defer client.Close()

	q := queue.NewRedisQueue(client)

	log.Println("Taskforge Worker started")
	log.Println("Waiting for tasks...")

	for {
		log.Println("Attempting to dequeue...")
		t, err := q.Dequeue(ctx)

		if err != nil {
			log.Printf("Failed to dequeue task: %v", err)
			time.Sleep(time.Second)
			continue
		}

		log.Printf("Successfully received task: %+v", t)

		// Mark as Processing
		if err := repo.UpdateTaskStatus(ctx, t.ID, task.StatusProcessing); err != nil {
			log.Printf("Failed to update task status: %v", err)
			continue
		}

		if err := ProcessTask(t); err != nil {
			log.Printf("Task failed: %s | %v", t.ID, err)

			if updateErr := repo.UpdateTaskStatus(ctx, t.ID, task.StatusFailed); updateErr != nil {
				log.Printf("Failed to mark task as FAILED: %v", updateErr)
			}
			continue
		}

		if err := repo.UpdateTaskStatus(ctx, t.ID, task.StatusCompleted); err != nil {
			log.Printf("Failed to mark as completed: %v", err)
			continue
		}

		log.Printf("Task %s successfully marked as COMPLETED", t.ID)
	}
}
