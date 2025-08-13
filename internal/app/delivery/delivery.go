package delivery

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/supchaser/LO_test_task/internal/app"
	"github.com/supchaser/LO_test_task/internal/app/models"
	"github.com/supchaser/LO_test_task/internal/utils/errs"
	"github.com/supchaser/LO_test_task/internal/utils/logger"
)

type TaskDelivery struct {
	taskUsecase app.TaskUsecase
}

func CreateTaskDelivery(taskUsecase app.TaskUsecase) *TaskDelivery {
	return &TaskDelivery{
		taskUsecase: taskUsecase,
	}
}

func (d *TaskDelivery) CreateTask(w http.ResponseWriter, r *http.Request) {
	const funcName = "Delivery.CreateTask"

	var req models.CreateTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Error("failed to decode request", err, map[string]any{
			"method": funcName,
		})
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	task, err := d.taskUsecase.CreateTask(r.Context(), req.Title, req.Description)
	if err != nil {
		logger.Error("failed to create task", err, map[string]any{
			"method": funcName,
			"title":  req.Title,
		})
		respondWithError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(task); err != nil {
		logger.Error("failed to encode response", err, map[string]any{
			"method":  funcName,
			"task_id": task.ID,
		})
	}
}

func (d *TaskDelivery) GetTask(w http.ResponseWriter, r *http.Request) {
	const funcName = "Delivery.GetTask"

	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		logger.Error("invalid task ID", err, map[string]any{
			"method": funcName,
			"id":     idStr,
		})
		http.Error(w, "invalid task ID", http.StatusBadRequest)
		return
	}

	task, err := d.taskUsecase.GetTask(r.Context(), id)
	if err != nil {
		logger.Error("failed to get task", err, map[string]any{
			"method": funcName,
			"id":     id,
		})
		respondWithError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(task)
}

func (d *TaskDelivery) ListTasks(w http.ResponseWriter, r *http.Request) {
	const funcName = "Delivery.ListTasks"

	statusFilter := models.TaskStatus(r.URL.Query().Get("status"))

	tasks, err := d.taskUsecase.ListTasks(r.Context(), statusFilter)
	if err != nil {
		logger.Error("failed to list tasks", err, map[string]any{
			"method": funcName,
			"status": statusFilter,
		})
		respondWithError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)
}

func (d *TaskDelivery) UpdateTask(w http.ResponseWriter, r *http.Request) {
	const funcName = "Delivery.UpdateTask"

	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		logger.Error("invalid task ID", err, map[string]any{
			"method": funcName,
			"id":     idStr,
		})
		http.Error(w, "invalid task ID", http.StatusBadRequest)
		return
	}

	req := models.UpdateTaskRequest{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Error("failed to decode request", err, map[string]any{
			"method": funcName,
			"id":     id,
		})
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	task, err := d.taskUsecase.UpdateTask(r.Context(), id, req.Title, req.Description, req.Status)
	if err != nil {
		logger.Error("failed to update task", err, map[string]any{
			"method": funcName,
			"id":     id,
		})
		respondWithError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(task)
}

func (d *TaskDelivery) DeleteTask(w http.ResponseWriter, r *http.Request) {
	const funcName = "Delivery.DeleteTask"

	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		logger.Error("invalid task ID", err, map[string]any{
			"method": funcName,
			"id":     idStr,
		})
		http.Error(w, "invalid task ID", http.StatusBadRequest)
		return
	}

	if err := d.taskUsecase.DeleteTask(r.Context(), id); err != nil {
		logger.Error("failed to delete task", err, map[string]any{
			"method": funcName,
			"id":     id,
		})
		respondWithError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func respondWithError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, errs.ErrTaskNotFound):
		http.Error(w, err.Error(), http.StatusNotFound)
	case errors.Is(err, errs.ErrInvalidID):
		http.Error(w, err.Error(), http.StatusBadRequest)
	case errors.Is(err, errs.ErrValidation):
		http.Error(w, err.Error(), http.StatusBadRequest)
	default:
		logger.Error("unhandled error", err, nil)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
}
