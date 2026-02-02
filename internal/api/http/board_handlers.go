package httpapi

import (
	"encoding/json"
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
			HandleError(w, err)
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
			HandleError(w, err)
			return
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
			HandleError(w, err)
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
			HandleError(w, err)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}

}
