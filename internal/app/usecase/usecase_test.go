package usecase

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	mock_app "github.com/supchaser/LO_test_task/internal/app/mocks"
	"github.com/supchaser/LO_test_task/internal/app/models"
	"github.com/supchaser/LO_test_task/internal/utils/errs"
	"github.com/supchaser/LO_test_task/internal/utils/validate"
)

func TestTaskUsecase_CreateTask(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	now := time.Now()
	mockTask := &models.Task{
		ID:          now.UnixNano() / int64(time.Millisecond),
		Title:       "Valid Title",
		Description: "Valid Description",
		Status:      models.StatusPending,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	tests := []struct {
		name          string
		title         string
		description   string
		mockSetup     func(*mock_app.MockTaskRepository)
		expectedTask  *models.Task
		expectedError error
	}{
		{
			name:        "Success",
			title:       "Valid Title",
			description: "Valid Description",
			mockSetup: func(mockRepo *mock_app.MockTaskRepository) {
				mockRepo.EXPECT().
					CreateTask(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, task *models.Task) (*models.Task, error) {
						return mockTask, nil
					})
			},
			expectedTask:  mockTask,
			expectedError: nil,
		},
		{
			name:          "Invalid Title",
			title:         "",
			description:   "Valid Description",
			mockSetup:     func(mockRepo *mock_app.MockTaskRepository) {},
			expectedTask:  nil,
			expectedError: fmt.Errorf("%w: task title cannot be empty", errs.ErrValidation),
		},
		{
			name:          "Invalid Description",
			title:         "Valid Title",
			description:   "This description is way too long and exceeds the maximum allowed length of 5000 characters. " + strings.Repeat("a", 5000),
			mockSetup:     func(mockRepo *mock_app.MockTaskRepository) {},
			expectedTask:  nil,
			expectedError: fmt.Errorf("%w: task description cannot be longer than %d characters", errs.ErrValidation, validate.MaxTaskDescriptionLength),
		},
		{
			name:        "Repository Error",
			title:       "Valid Title",
			description: "Valid Description",
			mockSetup: func(mockRepo *mock_app.MockTaskRepository) {
				mockRepo.EXPECT().
					CreateTask(gomock.Any(), gomock.Any()).
					Return(nil, errors.New("repository error"))
			},
			expectedTask:  nil,
			expectedError: errors.New("repository error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := mock_app.NewMockTaskRepository(ctrl)
			if tt.mockSetup != nil {
				tt.mockSetup(mockRepo)
			}

			uc := CreateTaskUsecase(mockRepo)
			result, err := uc.CreateTask(context.Background(), tt.title, tt.description)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedTask.ID, result.ID)
				assert.Equal(t, tt.expectedTask.Title, result.Title)
				assert.Equal(t, tt.expectedTask.Description, result.Description)
				assert.Equal(t, tt.expectedTask.Status, result.Status)
				assert.WithinDuration(t, tt.expectedTask.CreatedAt, result.CreatedAt, time.Second)
				assert.WithinDuration(t, tt.expectedTask.UpdatedAt, result.UpdatedAt, time.Second)
			}
		})
	}
}

func TestTaskUsecase_GetTask(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTask := &models.Task{
		ID:          1,
		Title:       "Test Task",
		Description: "Test Description",
		Status:      models.StatusPending,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	tests := []struct {
		name          string
		taskID        int64
		mockSetup     func(*mock_app.MockTaskRepository)
		expectedTask  *models.Task
		expectedError error
	}{
		{
			name:   "Success",
			taskID: 1,
			mockSetup: func(mockRepo *mock_app.MockTaskRepository) {
				mockRepo.EXPECT().
					GetTaskByID(gomock.Any(), int64(1)).
					Return(mockTask, nil)
			},
			expectedTask:  mockTask,
			expectedError: nil,
		},
		{
			name:   "Task Not Found",
			taskID: 2,
			mockSetup: func(mockRepo *mock_app.MockTaskRepository) {
				mockRepo.EXPECT().
					GetTaskByID(gomock.Any(), int64(2)).
					Return(nil, errors.New("task not found"))
			},
			expectedTask:  nil,
			expectedError: errors.New("task not found"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := mock_app.NewMockTaskRepository(ctrl)
			if tt.mockSetup != nil {
				tt.mockSetup(mockRepo)
			}

			uc := CreateTaskUsecase(mockRepo)
			result, err := uc.GetTask(context.Background(), tt.taskID)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedTask, result)
			}
		})
	}
}

func TestTaskUsecase_ListTasks(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTasks := []*models.Task{
		{
			ID:          1,
			Title:       "Task 1",
			Description: "Description 1",
			Status:      models.StatusPending,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			ID:          2,
			Title:       "Task 2",
			Description: "Description 2",
			Status:      models.StatusPending,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
	}

	tests := []struct {
		name          string
		statusFilter  models.TaskStatus
		mockSetup     func(*mock_app.MockTaskRepository)
		expectedTasks []*models.Task
		expectedError error
	}{
		{
			name:         "Success - No Filter",
			statusFilter: "",
			mockSetup: func(mockRepo *mock_app.MockTaskRepository) {
				mockRepo.EXPECT().
					GetAllTasks(gomock.Any(), models.TaskStatus("")).
					Return(mockTasks, nil)
			},
			expectedTasks: mockTasks,
			expectedError: nil,
		},
		{
			name:         "Success - With Filter",
			statusFilter: models.StatusPending,
			mockSetup: func(mockRepo *mock_app.MockTaskRepository) {
				mockRepo.EXPECT().
					GetAllTasks(gomock.Any(), models.StatusPending).
					Return(mockTasks, nil)
			},
			expectedTasks: mockTasks,
			expectedError: nil,
		},
		{
			name:         "Repository Error",
			statusFilter: "",
			mockSetup: func(mockRepo *mock_app.MockTaskRepository) {
				mockRepo.EXPECT().
					GetAllTasks(gomock.Any(), models.TaskStatus("")).
					Return(nil, errors.New("repository error"))
			},
			expectedTasks: nil,
			expectedError: errors.New("repository error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := mock_app.NewMockTaskRepository(ctrl)
			if tt.mockSetup != nil {
				tt.mockSetup(mockRepo)
			}

			uc := CreateTaskUsecase(mockRepo)
			result, err := uc.ListTasks(context.Background(), tt.statusFilter)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedTasks, result)
			}
		})
	}
}

func TestTaskUsecase_UpdateTask(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	now := time.Now()
	existingTask := &models.Task{
		ID:          1,
		Title:       "Old Title",
		Description: "Old Description",
		Status:      models.StatusPending,
		CreatedAt:   now.Add(-time.Hour),
		UpdatedAt:   now.Add(-time.Hour),
	}

	updatedTask := &models.Task{
		ID:          1,
		Title:       "New Title",
		Description: "New Description",
		Status:      models.StatusCompleted,
		CreatedAt:   now.Add(-time.Hour),
		UpdatedAt:   now,
	}

	tests := []struct {
		name           string
		taskID         int64
		newTitle       string
		newDescription string
		newStatus      models.TaskStatus
		mockSetup      func(*mock_app.MockTaskRepository)
		expectedTask   *models.Task
		expectedError  error
	}{
		{
			name:           "Success - Full Update",
			taskID:         1,
			newTitle:       "New Title",
			newDescription: "New Description",
			newStatus:      models.StatusCompleted,
			mockSetup: func(mockRepo *mock_app.MockTaskRepository) {
				mockRepo.EXPECT().
					GetTaskByID(gomock.Any(), int64(1)).
					Return(existingTask, nil)
				mockRepo.EXPECT().
					UpdateTask(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, task *models.Task) (*models.Task, error) {
						return updatedTask, nil
					})
			},
			expectedTask:  updatedTask,
			expectedError: nil,
		},
		{
			name:           "Task Not Found",
			taskID:         2,
			newTitle:       "New Title",
			newDescription: "New Description",
			newStatus:      models.StatusCompleted,
			mockSetup: func(mockRepo *mock_app.MockTaskRepository) {
				mockRepo.EXPECT().
					GetTaskByID(gomock.Any(), int64(2)).
					Return(nil, errors.New("task not found"))
			},
			expectedTask:  nil,
			expectedError: errors.New("task not found"),
		},
		{
			name:           "Repository Update Error",
			taskID:         1,
			newTitle:       "New Title",
			newDescription: "New Description",
			newStatus:      models.StatusCompleted,
			mockSetup: func(mockRepo *mock_app.MockTaskRepository) {
				mockRepo.EXPECT().
					GetTaskByID(gomock.Any(), int64(1)).
					Return(existingTask, nil)
				mockRepo.EXPECT().
					UpdateTask(gomock.Any(), gomock.Any()).
					Return(nil, errors.New("update error"))
			},
			expectedTask:  nil,
			expectedError: errors.New("update error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := mock_app.NewMockTaskRepository(ctrl)
			if tt.mockSetup != nil {
				tt.mockSetup(mockRepo)
			}

			uc := CreateTaskUsecase(mockRepo)
			result, err := uc.UpdateTask(
				context.Background(),
				tt.taskID,
				tt.newTitle,
				tt.newDescription,
				tt.newStatus,
			)

			if tt.expectedError != nil {
				assert.Error(t, err)
				if tt.expectedError != nil {
					assert.Contains(t, err.Error(), tt.expectedError.Error())
				}
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedTask.ID, result.ID)
				assert.Equal(t, tt.expectedTask.Title, result.Title)
				assert.Equal(t, tt.expectedTask.Description, result.Description)
				assert.Equal(t, tt.expectedTask.Status, result.Status)
				assert.Equal(t, tt.expectedTask.CreatedAt, result.CreatedAt)
				assert.WithinDuration(t, time.Now(), result.UpdatedAt, time.Second)
			}
		})
	}
}

func TestTaskUsecase_DeleteTask(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name          string
		taskID        int64
		mockSetup     func(*mock_app.MockTaskRepository)
		expectedError error
	}{
		{
			name:   "Success",
			taskID: 1,
			mockSetup: func(mockRepo *mock_app.MockTaskRepository) {
				mockRepo.EXPECT().
					DeleteTask(gomock.Any(), int64(1)).
					Return(nil)
			},
			expectedError: nil,
		},
		{
			name:   "Task Not Found",
			taskID: 2,
			mockSetup: func(mockRepo *mock_app.MockTaskRepository) {
				mockRepo.EXPECT().
					DeleteTask(gomock.Any(), int64(2)).
					Return(errors.New("task not found"))
			},
			expectedError: errors.New("task not found"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := mock_app.NewMockTaskRepository(ctrl)
			if tt.mockSetup != nil {
				tt.mockSetup(mockRepo)
			}

			uc := CreateTaskUsecase(mockRepo)
			err := uc.DeleteTask(context.Background(), tt.taskID)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
