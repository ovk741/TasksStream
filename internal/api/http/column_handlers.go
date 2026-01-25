package httpapi

import (
	"encoding/json"
	"net/http"

	"github.com/ovk741/TasksStream/internal/service"
)

func CreateColumnHandler(columnService service.ColumnService) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		var input struct {
			BoardID string `json:"board_id"`
			Title   string `json:"title"`
		}

		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		column, err := columnService.Create(input.Title, input.BoardID)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			_ = json.NewEncoder(w).Encode(map[string]string{
				"error": err.Error(),
			})
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(column)
	}
}

func GetColumnsByBoardHandler(columnService service.ColumnService) http.HandlerFunc {
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

		columns, err := columnService.GetByBoardID(boardID)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			_ = json.NewEncoder(w).Encode(map[string]string{
				"error": err.Error(),
			})
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(columns)
	}
}
