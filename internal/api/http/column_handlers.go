package httpapi

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/ovk741/TasksStream/internal/domain"
	"github.com/ovk741/TasksStream/internal/storage/memory"
)

func CreateColumnHandler(
	repo *memory.ColumnRepository,
	generateID func() string,
) http.HandlerFunc {
	var input struct {
		BoardID string `json:"board_id"`
		Name    string `json:"name"`
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

		column := domain.Column{
			ID:        generateID(),
			BoardID:   input.BoardID,
			Title:     input.Name,
			Position:  0,
			CreatedAt: time.Now(),
		}

		repo.Create(column)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(column)
	}
}

func GetColumnsByBoardHandler(repo *memory.ColumnRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		boardID := r.URL.Query().Get("board_id")
		if boardID == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		columns := repo.GetByBoardID(boardID)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(columns)
	}
}
