package httpapi

import (
	"encoding/json"

	"net/http"

	"github.com/ovk741/TasksStream/internal/domain"
	"github.com/ovk741/TasksStream/internal/service"
)

func CreateColumnHandler(columnService service.ColumnService) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		userID, ok := MustGetUserID(w, r)
		if !ok {
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

		column, err := columnService.Create(userID, input.Title, input.BoardID)
		if err != nil {
			HandleError(w, err)
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

		userID, ok := MustGetUserID(w, r)
		if !ok {
			return
		}

		boardID := r.URL.Query().Get("board_id")
		if boardID == "" {
			HandleError(w, domain.ErrInvalidInput)
			return
		}

		columns, err := columnService.GetByBoardID(userID, boardID)
		if err != nil {
			HandleError(w, err)
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

		userID, ok := MustGetUserID(w, r)
		if !ok {
			return
		}

		columnID := r.URL.Query().Get("id")
		if columnID == "" {
			HandleError(w, domain.ErrInvalidInput)
			return
		}

		var input struct {
			Title string `json:"title"`
		}

		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		column, err := columnService.Update(userID, columnID, input.Title)
		if err != nil {
			HandleError(w, err)
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

		userID, ok := MustGetUserID(w, r)
		if !ok {
			return
		}

		columnID := r.URL.Query().Get("id")
		if columnID == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		err := columnService.Delete(userID, columnID)
		if err != nil {
			HandleError(w, err)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}

}

func MoveColumnHandler(columnService service.ColumnService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		userID, ok := MustGetUserID(w, r)
		if !ok {
			return
		}

		columnID := r.URL.Query().Get("id")
		if columnID == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		var input struct {
			Position int `json:"position"`
		}

		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		column, err := columnService.Move(userID, columnID, input.Position)
		if err != nil {
			HandleError(w, err)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(column)
	}
}
