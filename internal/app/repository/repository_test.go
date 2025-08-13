package repository

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/supchaser/LO_test_task/internal/app/models"
	"github.com/supchaser/LO_test_task/internal/utils/errs"
)

func TestCreateTask_Success(t *testing.T) {
	repo := CreateTaskRepository()
	task := &models.Task{
		ID:     1,
		Title:  "Test Task",
		Status: models.StatusPending,
	}

	createdTask, err := repo.CreateTask(context.Background(), task)

	assert.NoError(t, err)
	assert.Equal(t, task.ID, createdTask.ID)
	assert.Equal(t, "Test Task", createdTask.Title)
	assert.Equal(t, models.StatusPending, createdTask.Status)
	assert.WithinDuration(t, time.Now(), createdTask.CreatedAt, time.Second)
	assert.WithinDuration(t, time.Now(), createdTask.UpdatedAt, time.Second)
}

func TestCreateTask_InvalidID(t *testing.T) {
	repo := CreateTaskRepository()
	task := &models.Task{
		ID:    -1,
		Title: "Invalid Task",
	}

	_, err := repo.CreateTask(context.Background(), task)

	assert.Error(t, err)
	assert.ErrorIs(t, err, errs.ErrInvalidID)
}

func TestGetTaskByID_Success(t *testing.T) {
	repo := CreateTaskRepository()
	task := &models.Task{ID: 1, Title: "Existing Task"}
	repo.tasks[task.ID] = task

	foundTask, err := repo.GetTaskByID(context.Background(), task.ID)

	assert.NoError(t, err)
	assert.Equal(t, task, foundTask)
}

func TestGetTaskByID_NotFound(t *testing.T) {
	repo := CreateTaskRepository()

	_, err := repo.GetTaskByID(context.Background(), 999)

	assert.Error(t, err)
	assert.ErrorIs(t, err, errs.ErrTaskNotFound)
}

func TestGetAllTasks_NoFilter(t *testing.T) {
	repo := CreateTaskRepository()
	tasks := []*models.Task{
		{ID: 1, Status: models.StatusPending},
		{ID: 2, Status: models.StatusInProgress},
		{ID: 3, Status: models.StatusCompleted},
	}
	for _, task := range tasks {
		repo.tasks[task.ID] = task
	}

	result, err := repo.GetAllTasks(context.Background(), "")

	assert.NoError(t, err)
	assert.Len(t, result, 3)
}

func TestGetAllTasks_WithFilter(t *testing.T) {
	repo := CreateTaskRepository()
	tasks := []*models.Task{
		{ID: 1, Status: models.StatusPending},
		{ID: 2, Status: models.StatusInProgress},
		{ID: 3, Status: models.StatusPending},
	}
	for _, task := range tasks {
		repo.tasks[task.ID] = task
	}

	result, err := repo.GetAllTasks(context.Background(), models.StatusPending)

	assert.NoError(t, err)
	assert.Len(t, result, 2)
	for _, task := range result {
		assert.Equal(t, models.StatusPending, task.Status)
	}
}

func TestUpdateTask_Success(t *testing.T) {
	repo := CreateTaskRepository()

	originalTask := &models.Task{
		ID:          1,
		Title:       "Original",
		Description: "Original Desc",
		Status:      models.StatusPending,
		CreatedAt:   time.Now().Add(-time.Hour),
		UpdatedAt:   time.Now().Add(-time.Hour),
	}
	repo.tasks[originalTask.ID] = originalTask

	updatedTask := &models.Task{
		ID:          1,
		Title:       "Updated",
		Description: "Updated Desc",
		Status:      models.StatusInProgress,
	}

	beforeUpdate := time.Now()
	result, err := repo.UpdateTask(context.Background(), updatedTask)
	assert.NoError(t, err)

	assert.Equal(t, "Updated", result.Title)
	assert.Equal(t, "Updated Desc", result.Description)
	assert.Equal(t, models.StatusInProgress, result.Status)

	assert.Equal(t, originalTask.CreatedAt, result.CreatedAt)

	assert.True(t, result.UpdatedAt.After(beforeUpdate) || result.UpdatedAt.Equal(beforeUpdate))
	assert.True(t, result.UpdatedAt.Before(time.Now().Add(time.Second)) || result.UpdatedAt.Equal(time.Now().Add(time.Second)))
}

func TestUpdateTask_NotFound(t *testing.T) {
	repo := CreateTaskRepository()
	task := &models.Task{ID: 999}

	_, err := repo.UpdateTask(context.Background(), task)

	assert.Error(t, err)
	assert.ErrorIs(t, err, errs.ErrTaskNotFound)
}

func TestDeleteTask_Success(t *testing.T) {
	repo := CreateTaskRepository()
	task := &models.Task{ID: 1}
	repo.tasks[task.ID] = task

	err := repo.DeleteTask(context.Background(), task.ID)

	assert.NoError(t, err)
	_, exists := repo.tasks[task.ID]
	assert.False(t, exists)
}

func TestDeleteTask_NotFound(t *testing.T) {
	repo := CreateTaskRepository()

	err := repo.DeleteTask(context.Background(), 999)

	assert.Error(t, err)
	assert.ErrorIs(t, err, errs.ErrTaskNotFound)
}

func TestConcurrentAccess(t *testing.T) {
	repo := CreateTaskRepository()
	count := 100

	var wg sync.WaitGroup
	wg.Add(count)

	for i := range count {
		go func(id int64) {
			defer wg.Done()
			_, err := repo.CreateTask(context.Background(), &models.Task{ID: id})
			assert.NoError(t, err)
		}(int64(i + 1))
	}

	wg.Wait()

	tasks, err := repo.GetAllTasks(context.Background(), "")
	assert.NoError(t, err)
	assert.Len(t, tasks, count)
}
