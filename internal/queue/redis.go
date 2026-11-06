package queue

import (
	"context"
	"encoding/json"
	"os"

	"github.com/redis/go-redis/v9"

	"taskforge/internal/task"
)

const taskQueueKey = "taskforge:tasks"

type RedisQueue struct {
	client *redis.Client
}

func NewRedisClient() *redis.Client {

	addr := os.Getenv("REDIS_ADDR")

	if addr == "" {
		addr = "localhost:6379"
	}
	return redis.NewClient(&redis.Options{
		Addr: addr,
	})
}

func NewRedisQueue(client *redis.Client) *RedisQueue {
	return &RedisQueue{
		client: client,
	}
}

func (q *RedisQueue) Enqueue(ctx context.Context, t task.Task) error {
	data, err := json.Marshal(t)
	if err != nil {
		return err
	}

	return q.client.LPush(ctx, taskQueueKey, data).Err()
}

func (q *RedisQueue) Dequeue(ctx context.Context) (task.Task, error) {
	data, err := q.client.BRPop(ctx, 0, taskQueueKey).Result()

	if err != nil {
		return task.Task{}, err
	}

	var t task.Task

	err = json.Unmarshal([]byte(data[1]), &t)

	if err != nil {
		return task.Task{}, err
	}

	return t, nil
}

func CheckRedis(client *redis.Client) error {
	return client.Ping(context.Background()).Err()
}
