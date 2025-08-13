package usecase

import (
	"context"
	"time"

	"github.com/supchaser/LO_test_task/internal/app"
	"github.com/supchaser/LO_test_task/internal/app/models"
	"github.com/supchaser/LO_test_task/internal/utils/logger"
	"github.com/supchaser/LO_test_task/internal/utils/validate"
)

type TaskUsecase struct {
	taskRepository app.TaskRepository
}

func CreateTaskUsecase(taskRepository app.TaskRepository) *TaskUsecase {
	return &TaskUsecase{
		taskRepository: taskRepository,
	}
}

func (u *TaskUsecase) CreateTask(ctx context.Context, title, description string) (*models.Task, error) {
	const funcName = "Usecase.CreateTask"

	if err := validate.CheckTaskTitle(title); err != nil {
		logger.Error("invalid task title", err, map[string]any{
			"method": funcName,
			"title":  title,
		})
		return nil, err
	}

	if err := validate.CheckTaskDescription(description); err != nil {
		logger.Error("invalid task description", err, map[string]any{
			"method":      funcName,
			"description": description,
		})
		return nil, err
	}

	task := &models.Task{
		ID:          time.Now().UnixNano() / int64(time.Millisecond),
		Title:       title,
		Description: description,
		Status:      models.StatusPending,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	createdTask, err := u.taskRepository.CreateTask(ctx, task)
	if err != nil {
		logger.Error("failed to create task in repository", err, map[string]any{
			"task_id": task.ID,
			"method":  funcName,
		})
		return nil, err
	}

	logger.Info("task created successfully", map[string]any{
		"task_id": createdTask.ID,
		"method":  funcName,
	})

	return createdTask, nil
}

func (u *TaskUsecase) GetTask(ctx context.Context, id int64) (*models.Task, error) {
	const funcName = "Usecase.GetTask"

	task, err := u.taskRepository.GetTaskByID(ctx, id)
	if err != nil {
		logger.Error("failed to get task", err, map[string]any{
			"task_id": id,
			"method":  funcName,
		})
		return nil, err
	}

	logger.Info("task retrieved", map[string]any{
		"task_id": task.ID,
		"method":  funcName,
	})

	return task, nil
}

func (u *TaskUsecase) ListTasks(ctx context.Context, statusFilter models.TaskStatus) ([]*models.Task, error) {
	const funcName = "Usecase.ListTasks"

	tasks, err := u.taskRepository.GetAllTasks(ctx, statusFilter)
	if err != nil {
		logger.Error("failed to list tasks", err, map[string]any{
			"method":        funcName,
			"status_filter": statusFilter,
		})
		return nil, err
	}

	logger.Info("tasks listed", map[string]any{
		"count":         len(tasks),
		"status_filter": statusFilter,
		"method":        funcName,
	})

	return tasks, nil
}

func (u *TaskUsecase) UpdateTask(ctx context.Context, id int64, newTitle, newDescription string, status models.TaskStatus) (*models.Task, error) {
	const funcName = "Usecase.UpdateTask"

	existingTask, err := u.taskRepository.GetTaskByID(ctx, id)
	if err != nil {
		logger.Error("task not found for update", err, map[string]any{
			"task_id": id,
			"method":  funcName,
		})
		return nil, err
	}

	if newTitle != "" {
		if err := validate.CheckTaskTitle(newTitle); err != nil {
			logger.Error("invalid new task title", err, map[string]any{
				"method":  funcName,
				"task_id": id,
				"title":   newTitle,
			})
			return nil, err
		}
		existingTask.Title = newTitle
	}

	if newDescription != "" {
		if err := validate.CheckTaskDescription(newDescription); err != nil {
			logger.Error("invalid new task description", err, map[string]any{
				"method":      funcName,
				"task_id":     id,
				"description": newDescription,
			})
			return nil, err
		}
		existingTask.Description = newDescription
	}

	if status != "" {
		existingTask.Status = status
	}

	existingTask.UpdatedAt = time.Now()

	updatedTask, err := u.taskRepository.UpdateTask(ctx, existingTask)
	if err != nil {
		logger.Error("failed to update task", err, map[string]any{
			"task_id": id,
			"method":  funcName,
		})
		return nil, err
	}

	logger.Info("task updated", map[string]any{
		"task_id": updatedTask.ID,
		"method":  funcName,
	})

	return updatedTask, nil
}

func (u *TaskUsecase) DeleteTask(ctx context.Context, id int64) error {
	const funcName = "Usecase.DeleteTask"

	if err := u.taskRepository.DeleteTask(ctx, id); err != nil {
		logger.Error("failed to delete task", err, map[string]any{
			"task_id": id,
			"method":  funcName,
		})
		return err
	}

	logger.Info("task deleted", map[string]any{
		"task_id": id,
		"method":  funcName,
	})

	return nil
}
