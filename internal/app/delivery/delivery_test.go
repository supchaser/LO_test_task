package delivery

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	mock_app "github.com/supchaser/LO_test_task/internal/app/mocks"
	"github.com/supchaser/LO_test_task/internal/app/models"
	"github.com/supchaser/LO_test_task/internal/utils/errs"
)

func TestTaskDelivery_CreateTask(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUsecase := mock_app.NewMockTaskUsecase(ctrl)
	delivery := CreateTaskDelivery(mockUsecase)

	tests := []struct {
		name           string
		requestBody    interface{}
		mockSetup      func()
		expectedStatus int
	}{
		{
			name: "Success",
			requestBody: models.CreateTaskRequest{
				Title:       "Test Task",
				Description: "Test Description",
			},
			mockSetup: func() {
				mockUsecase.EXPECT().
					CreateTask(gomock.Any(), "Test Task", "Test Description").
					Return(&models.Task{
						ID:          1,
						Title:       "Test Task",
						Description: "Test Description",
						Status:      models.StatusPending,
						CreatedAt:   time.Now(),
						UpdatedAt:   time.Now(),
					}, nil)
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name:           "Invalid Request Body",
			requestBody:    "invalid",
			mockSetup:      func() {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "Validation Error",
			requestBody: models.CreateTaskRequest{
				Title:       "",
				Description: "Test Description",
			},
			mockSetup: func() {
				mockUsecase.EXPECT().
					CreateTask(gomock.Any(), "", "Test Description").
					Return(nil, errs.ErrValidation)
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "Internal Server Error",
			requestBody: models.CreateTaskRequest{
				Title:       "Test Task",
				Description: "Test Description",
			},
			mockSetup: func() {
				mockUsecase.EXPECT().
					CreateTask(gomock.Any(), "Test Task", "Test Description").
					Return(nil, errors.New("internal error"))
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			body, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest("POST", "/tasks", bytes.NewBuffer(body))
			w := httptest.NewRecorder()

			delivery.CreateTask(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			if tt.expectedStatus == http.StatusCreated {
				var task models.Task
				err := json.NewDecoder(w.Body).Decode(&task)
				assert.NoError(t, err)
				assert.Equal(t, "Test Task", task.Title)
			}
		})
	}
}

func TestTaskDelivery_GetTask(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUsecase := mock_app.NewMockTaskUsecase(ctrl)
	delivery := CreateTaskDelivery(mockUsecase)

	tests := []struct {
		name           string
		taskID         string
		mockSetup      func()
		expectedStatus int
	}{
		{
			name:   "Success",
			taskID: "1",
			mockSetup: func() {
				mockUsecase.EXPECT().
					GetTask(gomock.Any(), int64(1)).
					Return(&models.Task{
						ID:          1,
						Title:       "Test Task",
						Description: "Test Description",
						Status:      models.StatusPending,
					}, nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Invalid ID",
			taskID:         "invalid",
			mockSetup:      func() {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:   "Task Not Found",
			taskID: "1",
			mockSetup: func() {
				mockUsecase.EXPECT().
					GetTask(gomock.Any(), int64(1)).
					Return(nil, errs.ErrTaskNotFound)
			},
			expectedStatus: http.StatusNotFound,
		},
		{
			name:   "Internal Server Error",
			taskID: "1",
			mockSetup: func() {
				mockUsecase.EXPECT().
					GetTask(gomock.Any(), int64(1)).
					Return(nil, errors.New("internal error"))
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			req := httptest.NewRequest("GET", "/tasks/"+tt.taskID, nil)
			w := httptest.NewRecorder()

			req.SetPathValue("id", tt.taskID)

			delivery.GetTask(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}

func TestTaskDelivery_ListTasks(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUsecase := mock_app.NewMockTaskUsecase(ctrl)
	delivery := CreateTaskDelivery(mockUsecase)

	tests := []struct {
		name           string
		statusFilter   string
		mockSetup      func(statusFilter models.TaskStatus)
		expectedStatus int
	}{
		{
			name:         "Success - No Filter",
			statusFilter: "",
			mockSetup: func(statusFilter models.TaskStatus) {
				mockUsecase.EXPECT().
					ListTasks(gomock.Any(), statusFilter).
					Return([]*models.Task{
						{
							ID:          1,
							Title:       "Task 1",
							Description: "Description 1",
							Status:      models.StatusPending,
						},
					}, nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:         "Success - With Filter",
			statusFilter: "pending",
			mockSetup: func(statusFilter models.TaskStatus) {
				mockUsecase.EXPECT().
					ListTasks(gomock.Any(), statusFilter).
					Return([]*models.Task{
						{
							ID:          1,
							Title:       "Task 1",
							Description: "Description 1",
							Status:      models.StatusPending,
						},
					}, nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:         "Internal Server Error",
			statusFilter: "",
			mockSetup: func(statusFilter models.TaskStatus) {
				mockUsecase.EXPECT().
					ListTasks(gomock.Any(), statusFilter).
					Return(nil, errors.New("internal error"))
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			status := models.TaskStatus(tt.statusFilter)
			tt.mockSetup(status)

			req := httptest.NewRequest("GET", "/tasks?status="+tt.statusFilter, nil)
			w := httptest.NewRecorder()

			delivery.ListTasks(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}

func TestTaskDelivery_UpdateTask(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUsecase := mock_app.NewMockTaskUsecase(ctrl)
	delivery := CreateTaskDelivery(mockUsecase)

	tests := []struct {
		name           string
		taskID         string
		requestBody    interface{}
		mockSetup      func()
		expectedStatus int
	}{
		{
			name:   "Success",
			taskID: "1",
			requestBody: models.UpdateTaskRequest{
				Title:       "Updated Title",
				Description: "Updated Description",
				Status:      models.StatusCompleted,
			},
			mockSetup: func() {
				mockUsecase.EXPECT().
					UpdateTask(gomock.Any(), int64(1), "Updated Title", "Updated Description", models.StatusCompleted).
					Return(&models.Task{
						ID:          1,
						Title:       "Updated Title",
						Description: "Updated Description",
						Status:      models.StatusCompleted,
					}, nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Invalid ID",
			taskID:         "invalid",
			requestBody:    models.UpdateTaskRequest{},
			mockSetup:      func() {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Invalid Request Body",
			taskID:         "1",
			requestBody:    "invalid",
			mockSetup:      func() {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:   "Validation Error",
			taskID: "1",
			requestBody: models.UpdateTaskRequest{
				Title:       "",
				Description: "Updated Description",
				Status:      models.StatusCompleted,
			},
			mockSetup: func() {
				mockUsecase.EXPECT().
					UpdateTask(gomock.Any(), int64(1), "", "Updated Description", models.StatusCompleted).
					Return(nil, errs.ErrValidation)
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:   "Task Not Found",
			taskID: "1",
			requestBody: models.UpdateTaskRequest{
				Title:       "Updated Title",
				Description: "Updated Description",
				Status:      models.StatusCompleted,
			},
			mockSetup: func() {
				mockUsecase.EXPECT().
					UpdateTask(gomock.Any(), int64(1), "Updated Title", "Updated Description", models.StatusCompleted).
					Return(nil, errs.ErrTaskNotFound)
			},
			expectedStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			body, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest("PUT", "/tasks/"+tt.taskID, bytes.NewBuffer(body))
			w := httptest.NewRecorder()

			req.SetPathValue("id", tt.taskID)

			delivery.UpdateTask(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}

func TestTaskDelivery_DeleteTask(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUsecase := mock_app.NewMockTaskUsecase(ctrl)
	delivery := CreateTaskDelivery(mockUsecase)

	tests := []struct {
		name           string
		taskID         string
		mockSetup      func()
		expectedStatus int
	}{
		{
			name:   "Success",
			taskID: "1",
			mockSetup: func() {
				mockUsecase.EXPECT().
					DeleteTask(gomock.Any(), int64(1)).
					Return(nil)
			},
			expectedStatus: http.StatusNoContent,
		},
		{
			name:           "Invalid ID",
			taskID:         "invalid",
			mockSetup:      func() {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:   "Task Not Found",
			taskID: "1",
			mockSetup: func() {
				mockUsecase.EXPECT().
					DeleteTask(gomock.Any(), int64(1)).
					Return(errs.ErrTaskNotFound)
			},
			expectedStatus: http.StatusNotFound,
		},
		{
			name:   "Internal Server Error",
			taskID: "1",
			mockSetup: func() {
				mockUsecase.EXPECT().
					DeleteTask(gomock.Any(), int64(1)).
					Return(errors.New("internal error"))
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			req := httptest.NewRequest("DELETE", "/tasks/"+tt.taskID, nil)
			w := httptest.NewRecorder()

			req.SetPathValue("id", tt.taskID)

			delivery.DeleteTask(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}

func TestRespondWithError(t *testing.T) {
	tests := []struct {
		name           string
		err            error
		expectedStatus int
	}{
		{
			name:           "Task Not Found",
			err:            errs.ErrTaskNotFound,
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "Validation Error",
			err:            errs.ErrValidation,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Internal Server Error",
			err:            errors.New("internal error"),
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			respondWithError(w, tt.err)
			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}
