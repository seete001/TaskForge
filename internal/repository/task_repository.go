package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	"taskforge/internal/task"
)

type TaskRepository struct {
	db *pgxpool.Pool
}

func NewTaskRepository(db *pgxpool.Pool) *TaskRepository {
	return &TaskRepository{
		db: db,
	}
}

func (r *TaskRepository) CreateTask(ctx context.Context, t task.Task) error {
	query := `INSERT INTO tasks (id, name, status, created_at, updated_at)
			  VALUES ($1, $2, $3, $4, $5)
			  `
	_, err := r.db.Exec(ctx, query, t.ID, t.Name, t.Status, t.CreatedAt, t.CreatedAt)

	return err
}

func (r *TaskRepository) GetAllTasks(ctx context.Context) ([]task.Task, error) {
	query := `SELECT id, name, status, created_at
			  FROM tasks
			  ORDER BY created_at DESC
			  `
	rows, err := r.db.Query(ctx, query)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	tasks := make([]task.Task, 0)

	for rows.Next() {
		var t task.Task

		err := rows.Scan(&t.ID, &t.Name, &t.Status, &t.CreatedAt)

		if err != nil {
			return nil, err
		}

		tasks = append(tasks, t)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return tasks, nil
}

func (r *TaskRepository) GetTaskByID(ctx context.Context, id string) (task.Task, error) {
	query := `SELECT id, name, status, created_at
			  FROM tasks
			  WHERE id = $1
	`

	var t task.Task

	err := r.db.QueryRow(ctx, query, id).Scan(&t.ID, &t.Name, &t.Status, &t.CreatedAt)

	return t, err
}

func (r *TaskRepository) UpdateTaskStatus(ctx context.Context, id, status string) error {
	query := `UPDATE tasks
	 		  SET status = $1, updated_at = NOW()
			  WHERE id = $2
			`
	_, err := r.db.Exec(ctx, query, status, id)

	return err
}
