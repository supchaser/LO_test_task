package app

import (
	"context"

	"github.com/supchaser/LO_test_task/internal/app/models"
)

//go:generate mockgen -source=interfaces.go -destination=mocks/mock.go

type TaskRepository interface {
	CreateTask(ctx context.Context, task *models.Task) (*models.Task, error)
	GetTaskByID(ctx context.Context, id int64) (*models.Task, error)
	GetAllTasks(ctx context.Context, statusFilter models.TaskStatus) ([]*models.Task, error)
	UpdateTask(ctx context.Context, task *models.Task) (*models.Task, error)
	DeleteTask(ctx context.Context, id int64) error
}

type TaskUsecase interface {
	CreateTask(ctx context.Context, title, description string) (*models.Task, error)
	GetTask(ctx context.Context, id int64) (*models.Task, error)
	ListTasks(ctx context.Context, statusFilter models.TaskStatus) ([]*models.Task, error)
	UpdateTask(ctx context.Context, id int64, newTitle, newDescription string, status models.TaskStatus) (*models.Task, error)
	DeleteTask(ctx context.Context, id int64) error
}
