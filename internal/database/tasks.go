package database

import (
	"database/sql"
	"errors"
	"fmt"
	"restapi/internal/models"
	"time"

	"github.com/jmoiron/sqlx"
)

type TaskStore struct {
	db *sqlx.DB
}

func NewTaskStore(db *sqlx.DB) *TaskStore {
	return &TaskStore{db: db}
}

func (s *TaskStore) GetAll() ([]models.Task, error) {
	var tasks []models.Task

	query := `
SELECT id, title, description, completed, created_at, updated_at
FROM tasks
order by created_at desc;
`

	err := s.db.Select(&tasks, query)
	if err != nil {
		return nil, fmt.Errorf("error during getting all tasks from db: %w", err)
	}

	return tasks, nil
}

func (s *TaskStore) GetByID(id int) (*models.Task, error) {
	var task models.Task

	query := `
SELECT id, title, description, completed, created_at, updated_at
FROM tasks
WHERE id = $1;
`

	err := s.db.Get(&task, query, id)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("task with id %d not found: %w", id, err)
	}
	if err != nil {
		return nil, fmt.Errorf("error during getting task by id from db: %w", err)
	}

	return &task, nil
}

func (s *TaskStore) Create(newTask *models.CreateTaskInput) (*models.Task, error) {
	var task models.Task

	query := `
INSERT INTO tasks (title, description, completed, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5)
RETURNING id, title, description, completed, created_at, updated_at;
`

	now := time.Now()

	err := s.db.QueryRowx(query, newTask.Title, newTask.Description, newTask.Completed, now, now).StructScan(&task)
	if err != nil {
		return nil, fmt.Errorf("error during creating task: %w", err)
	}

	return &task, nil
}

func (s *TaskStore) Update(id int, newTask *models.UpdateTaskInput) (*models.Task, error) {
	task, getErr := s.GetByID(id)
	if getErr != nil {
		return nil, fmt.Errorf("error during updating task: %w", getErr)
	}

	if newTask.Title != nil {
		task.Title = *newTask.Title
	}
	if newTask.Description != nil {
		task.Description = *newTask.Description
	}
	if newTask.Completed != nil {
		task.Completed = *newTask.Completed
	}

	query := `
UPDATE tasks
SET title = $1, description = $2, completed = $3, updated_at = $4
WHERE id = $5
RETURNING id, title, description, completed, created_at, updated_at;
`

	var updatedTask models.Task

	err := s.db.QueryRowx(query, task.Title, task.Description, task.Completed, task.UpdatedAt, id).StructScan(&updatedTask)
	if err != nil {
		return nil, fmt.Errorf("error during updating task: %w", err)
	}

	return task, nil
}

func (s *TaskStore) Delete(id int) error {
	query := `DELETE FROM tasks WHERE id = $1;`

	result, err := s.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("error during deleting task: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error during deleting task: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("task with id %d not found", id)
	}

	return nil
}
