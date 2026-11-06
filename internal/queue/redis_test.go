package queue

import (
	"context"
	"testing"
	"time"

	"taskforge/internal/task"

	"github.com/redis/go-redis/v9"
)

func TestRedisQueue(t *testing.T) {
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	defer client.Close()

	ctx := context.Background()

	if err := client.Ping(ctx).Err(); err != nil {
		t.Skip("Redis is not available:", err)
	}

	q := NewRedisQueue(client)

	original := task.Task{
		ID:     "test-task-1",
		Name:   "send-email",
		Status: task.StatusPending,
	}

	key := taskQueueKey
	client.Del(ctx, key)

	if err := q.Enqueue(ctx, original); err != nil {
		t.Fatal(err)
	}

	result, err := q.Dequeue(ctx)
	if err != nil {
		t.Fatal(err)
	}

	if result.ID != original.ID {
		t.Fatalf("expected %s, got %s", original.ID, result.ID)
	}

	client.Del(ctx, key)

	_ = time.Second
}
