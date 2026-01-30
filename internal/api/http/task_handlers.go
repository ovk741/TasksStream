package httpapi

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/ovk741/TasksStream/internal/service"
)

func CreateTaskHandler(taskService service.TaskService) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		var input struct {
			Title       string `json:"title"`
			ColumnID    string `json:"column_id"`
			Description string `json:"description"`
		}

		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			SendError(w, http.StatusBadRequest, err)
			return
		}

		task, err := taskService.Create(input.Title, input.Description, input.ColumnID)
		if err != nil {
			if errors.Is(err, service.ErrNotFound) {
				SendError(w, http.StatusNotFound, err)
				return
			}
			if errors.Is(err, service.ErrInvalidInput) {
				SendError(w, http.StatusBadRequest, err)
				return
			}

			SendError(w, http.StatusInternalServerError, err)

			return

		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(task)
	}
}

func GetTasksByColumnHandler(taskService service.TaskService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		columnID := r.URL.Query().Get("column_id")
		if columnID == "" {
			SendError(w, http.StatusBadRequest, errors.New("missing column_id"))
			return
		}

		tasks, err := taskService.GetByColumnID(columnID)
		if err != nil {
			if errors.Is(err, service.ErrInvalidInput) {
				SendError(w, http.StatusBadRequest, err)
				return
			}
			if errors.Is(err, service.ErrNotFound) {

				SendError(w, http.StatusNotFound, err)

				return
			}

			SendError(w, http.StatusInternalServerError, err)
			return

		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(tasks)
	}
}

func UpdateTaskHandler(taskService service.TaskService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		taskID := r.URL.Query().Get("id")
		if taskID == "" {
			w.WriteHeader(http.StatusBadRequest)
			_ = json.NewEncoder(w).Encode(map[string]string{
				"error": "task_id is required",
			})
			return
		}

		var input struct {
			Title       string `json:"title"`
			Description string `json:"description"`
		}

		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		task, err := taskService.Update(taskID, input.Title, input.Description)
		if err != nil {
			switch {
			case errors.Is(err, service.ErrNotFound):
				w.WriteHeader(http.StatusNotFound)

			case errors.Is(err, service.ErrInvalidInput):
				w.WriteHeader(http.StatusBadRequest)

			default:
				w.WriteHeader(http.StatusInternalServerError)
			}

			_ = json.NewEncoder(w).Encode(map[string]string{
				"error": err.Error(),
			})
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(task)
	}
}

func DeleteTaskHandler(taskService service.TaskService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		taskID := r.URL.Query().Get("id")
		if taskID == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		err := taskService.Delete(taskID)
		if err != nil {
			if errors.Is(err, service.ErrNotFound) {
				w.WriteHeader(http.StatusNotFound)
				return
			}
			if errors.Is(err, service.ErrInvalidInput) {
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}

}

func MoveTaskHandler(taskService service.TaskService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		taskID := r.URL.Query().Get("id")
		if taskID == "" {
			w.WriteHeader(http.StatusBadRequest)
			_ = json.NewEncoder(w).Encode(map[string]string{
				"error": "task id is required",
			})
			return
		}

		var input struct {
			ColumnID string `json:"column_id"`
			Position int    `json:"position"`
		}

		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		task, err := taskService.Move(taskID, input.ColumnID, input.Position)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			switch {
			case errors.Is(err, service.ErrNotFound):
				w.WriteHeader(http.StatusNotFound)
			case errors.Is(err, service.ErrInvalidInput):
				w.WriteHeader(http.StatusBadRequest)
			default:
				w.WriteHeader(http.StatusInternalServerError)
			}
			_ = json.NewEncoder(w).Encode(map[string]string{
				"error": err.Error(),
			})
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(task)
	}
}
