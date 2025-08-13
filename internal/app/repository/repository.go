package repository

import (
	"context"
	"sync"
	"time"

	"github.com/supchaser/LO_test_task/internal/app/models"
	"github.com/supchaser/LO_test_task/internal/utils/errs"
	"github.com/supchaser/LO_test_task/internal/utils/logger"
)

type TaskRepository struct {
	tasks map[int64]*models.Task
	mu    sync.RWMutex
}

func CreateTaskRepository() *TaskRepository {
	return &TaskRepository{
		tasks: make(map[int64]*models.Task),
	}
}

func (r *TaskRepository) CreateTask(ctx context.Context, task *models.Task) (*models.Task, error) {
	const funcName = "Repository.CreateTask"

	r.mu.Lock()
	defer r.mu.Unlock()

	if task.ID < 0 {
		logger.Error("invalid task ID", errs.ErrInvalidID, map[string]any{
			"task_id": task.ID,
			"method":  funcName,
		})
		return nil, errs.ErrInvalidID
	}

	now := time.Now()
	task.CreatedAt = now
	task.UpdatedAt = now

	r.tasks[task.ID] = task

	logger.Info("task created", map[string]any{
		"task_id": task.ID,
		"method":  funcName,
	})

	return task, nil
}

func (r *TaskRepository) GetTaskByID(ctx context.Context, id int64) (*models.Task, error) {
	const funcName = "Repository.GetTaskByID"

	r.mu.RLock()
	defer r.mu.RUnlock()

	task, exists := r.tasks[id]
	if !exists {
		logger.Error("task not found", errs.ErrTaskNotFound, map[string]any{
			"task_id": id,
			"method":  funcName,
		})
		return nil, errs.ErrTaskNotFound
	}

	logger.Info("task retrieved", map[string]any{
		"task_id": id,
		"method":  funcName,
	})

	return task, nil
}

func (r *TaskRepository) GetAllTasks(ctx context.Context, statusFilter models.TaskStatus) ([]*models.Task, error) {
	const funcName = "Repository.GetAllTasks"

	r.mu.RLock()
	defer r.mu.RUnlock()

	tasks := []*models.Task{}
	for _, task := range r.tasks {
		if statusFilter == "" || task.Status == statusFilter {
			tasks = append(tasks, task)
		}
	}

	logger.Info("tasks list retrieved", map[string]any{
		"count":         len(tasks),
		"status_filter": statusFilter,
		"method":        funcName,
	})

	return tasks, nil
}

func (r *TaskRepository) UpdateTask(ctx context.Context, task *models.Task) (*models.Task, error) {
	const funcName = "Repository.UpdateTask"

	r.mu.Lock()
	defer r.mu.Unlock()

	existingTask, exists := r.tasks[task.ID]
	if !exists {
		logger.Error("task not found for update", errs.ErrTaskNotFound, map[string]any{
			"task_id": task.ID,
			"method":  funcName,
		})
		return nil, errs.ErrTaskNotFound
	}

	existingTask.Title = task.Title
	existingTask.Description = task.Description
	existingTask.Status = task.Status
	existingTask.UpdatedAt = time.Now()

	logger.Info("task updated", map[string]any{
		"task_id": task.ID,
		"method":  funcName,
	})

	return existingTask, nil
}

func (r *TaskRepository) DeleteTask(ctx context.Context, id int64) error {
	const funcName = "Repository.DeleteTask"

	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.tasks[id]; !exists {
		logger.Error("task not found for deletion", errs.ErrTaskNotFound, map[string]any{
			"task_id": id,
			"method":  funcName,
		})
		return errs.ErrTaskNotFound
	}

	delete(r.tasks, id)

	logger.Info("task deleted", map[string]any{
		"task_id": id,
		"method":  funcName,
	})

	return nil
}
