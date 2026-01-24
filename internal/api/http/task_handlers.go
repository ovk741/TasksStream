package httpapi

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/ovk741/TasksStream/internal/domain"
	"github.com/ovk741/TasksStream/internal/storage"
)

func CreateTaskHandler(
	repo storage.TaskRepository,
	generateID func() string,
) http.HandlerFunc {
	var input struct {
		ColumnID    string `json:"column_id"`
		Title       string `json:"title"`
		Description string `json:"description"`
	}
	return func(w http.ResponseWriter, r *http.Request) {

		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		task := domain.Task{
			ID:          generateID(),
			ColumnID:    input.ColumnID,
			Title:       input.Title,
			Description: input.Description,
			Position:    0,
			CreatedAt:   time.Now(),
		}

		repo.Create(task)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(task)
	}
}

func GetTasksByColumnHandler(repo storage.TaskRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		columnID := r.URL.Query().Get("column_id")
		if columnID == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		tasks := repo.GetByColumnID(columnID)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(tasks)
	}
}
