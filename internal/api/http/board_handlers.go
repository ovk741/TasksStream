package httpapi

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/ovk741/TasksStream/internal/domain"
	"github.com/ovk741/TasksStream/internal/storage"
)

func CreateBoardHandler(repo storage.BoardRepository, generateID func() string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		var input struct {
			Name string `json:"name"`
		}

		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		board := domain.Board{
			ID:        generateID(),
			Name:      input.Name,
			CreatedAt: time.Now(),
		}

		repo.Create(board)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(board)
	}
}

func GetBoardsHandler(repo storage.BoardRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		boards := repo.GetAll()

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(boards)
	}
}
