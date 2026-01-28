package httpapi

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/ovk741/TasksStream/internal/service"
)

func CreateBoardHandler(boardService service.BoardService) http.HandlerFunc {
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

		board, err := boardService.Create(input.Name)
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
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(board)
	}
}

func GetBoardsHandler(boardService service.BoardService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		boards, err := boardService.GetAll()
		if err != nil {
			if errors.Is(err, service.ErrInvalidInput) {
				SendError(w, http.StatusInternalServerError, err)
				return
			}
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(boards)
	}
}

func UpdateBoardHandler(boardService service.BoardService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		boardID := r.URL.Query().Get("id")
		if boardID == "" {
			w.WriteHeader(http.StatusBadRequest)
			_ = json.NewEncoder(w).Encode(map[string]string{
				"error": "board_id is required",
			})
			return
		}

		var input struct {
			Name string `json:"name"`
		}

		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		board, err := boardService.Update(boardID, input.Name)
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
		json.NewEncoder(w).Encode(board)
	}
}

func DeleteBoardHandler(boardService service.BoardService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		boardID := r.URL.Query().Get("id")
		if boardID == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		err := boardService.Delete(boardID)
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
