package httpapi

import (
	"encoding/json"
	"errors"
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
			if errors.Is(err, service.ErrInvalidInput) {
				SendError(w, http.StatusBadRequest, err)
				return
			}

			SendError(w, http.StatusInternalServerError, err)

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
			if errors.Is(err, service.ErrInvalidInput) {
				SendError(w, http.StatusBadRequest, err)
				return
			}
			SendError(w, http.StatusInternalServerError, err)

			return

		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(columns)
	}
}

func UpdateColumnHandler(columnService service.ColumnService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		columnID := r.URL.Query().Get("id")
		if columnID == "" {
			w.WriteHeader(http.StatusBadRequest)
			_ = json.NewEncoder(w).Encode(map[string]string{
				"error": "column_id is required",
			})
			return
		}

		var input struct {
			Title string `json:"title"`
		}

		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		column, err := columnService.Update(columnID, input.Title)
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
		json.NewEncoder(w).Encode(column)
	}
}

func DeleteColumnHandler(columnService service.ColumnService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		columnID := r.URL.Query().Get("id")
		if columnID == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		err := columnService.Delete(columnID)
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
