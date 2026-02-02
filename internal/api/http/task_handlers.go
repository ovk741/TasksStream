package httpapi

import (
	"encoding/json"
	"net/http"

	"github.com/ovk741/TasksStream/internal/domain"
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
			HandleError(w, domain.ErrInvalidInput)
			return
		}

		task, err := taskService.Create(input.Title, input.Description, input.ColumnID)
		if err != nil {
			HandleError(w, err)
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
			HandleError(w, domain.ErrInvalidInput)
			return
		}

		tasks, err := taskService.GetByColumnID(columnID)
		if err != nil {
			HandleError(w, err)
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
			HandleError(w, domain.ErrInvalidInput)
			return
		}

		var input struct {
			Title       string `json:"title"`
			Description string `json:"description"`
		}

		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			HandleError(w, domain.ErrInvalidInput)
			return
		}

		task, err := taskService.Update(taskID, input.Title, input.Description)
		if err != nil {
			HandleError(w, err)
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
			HandleError(w, domain.ErrInvalidInput)
			return
		}
		err := taskService.Delete(taskID)
		if err != nil {
			HandleError(w, err)
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
			HandleError(w, domain.ErrInvalidInput)
			return
		}

		var input struct {
			ColumnID string `json:"column_id"`
			Position int    `json:"position"`
		}

		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			HandleError(w, domain.ErrInvalidInput)
			return
		}

		task, err := taskService.Move(taskID, input.ColumnID, input.Position)
		if err != nil {
			HandleError(w, err)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(task)
	}
}
