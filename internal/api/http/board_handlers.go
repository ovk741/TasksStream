package httpapi

import (
	"encoding/json"
	"net/http"

	"github.com/ovk741/TasksStream/internal/domain"
	"github.com/ovk741/TasksStream/internal/service"
)

func CreateBoardHandler(boardService service.BoardService) http.HandlerFunc {
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
			Name string `json:"name"`
		}

		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		board, err := boardService.Create(userID, input.Name)
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

		userID, ok := MustGetUserID(w, r)
		if !ok {
			return
		}

		boards, err := boardService.GetAll(userID)
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
		userID, ok := MustGetUserID(w, r)
		if !ok {
			return
		}

		board, err := boardService.Update(userID, boardID, input.Name)
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
		userID, ok := MustGetUserID(w, r)
		if !ok {
			return
		}

		err := boardService.Delete(userID, boardID)
		if err != nil {
			HandleError(w, err)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}

}

func InviteToBoardHandler(boardService service.BoardService) http.HandlerFunc {
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
			BoardID string           `json:"board_id"`
			UserID  string           `json:"user_id"`
			Role    domain.BoardRole `json:"role"`
		}

		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if err := boardService.InviteUser(
			userID,
			input.BoardID,
			input.UserID,
			input.Role,
		); err != nil {
			HandleError(w, err)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

func GetBoardMembersHandler(boardService service.BoardService) http.HandlerFunc {
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
			w.WriteHeader(http.StatusBadRequest)
			_ = json.NewEncoder(w).Encode(map[string]string{
				"error": "board_id is required",
			})
			return
		}

		members, err := boardService.GetMembers(userID, boardID)
		if err != nil {
			HandleError(w, err)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(members)
	}
}

func RemoveBoardMemberHandler(boardService service.BoardService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		if r.Method != http.MethodDelete {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		requesterID, ok := MustGetUserID(w, r)
		if !ok {
			return
		}

		var input struct {
			BoardID string `json:"board_id"`
			UserID  string `json:"user_id"`
		}

		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if input.BoardID == "" || input.UserID == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if err := boardService.RemoveUser(requesterID, input.BoardID, input.UserID); err != nil {
			HandleError(w, err)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
