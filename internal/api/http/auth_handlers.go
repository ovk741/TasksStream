package httpapi

import (
	"encoding/json"
	"net/http"

	"github.com/ovk741/TasksStream/internal/service"
)

func RegisterHandler(authService service.AuthService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		var input struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}

		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			SendError(w, http.StatusBadRequest, err)
			return
		}

		if err := authService.Register(input.Email, input.Password); err != nil {
			HandleError(w, err)
			return
		}

		w.WriteHeader(http.StatusCreated)
	}
}

func LoginHandler(authService service.AuthService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		var input struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}

		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			SendError(w, http.StatusBadRequest, err)
			return
		}

		tokens, err := authService.Login(input.Email, input.Password)
		if err != nil {
			HandleError(w, err)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]string{
			"access_token":  tokens.AccessToken,
			"refresh_token": tokens.RefreshToken,
		})
	}
}

func RefreshHandler(authService service.AuthService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		var input struct {
			RefreshToken string `json:"refresh_token"`
		}

		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			SendError(w, http.StatusBadRequest, err)
			return
		}

		tokens, err := authService.Refresh(input.RefreshToken)
		if err != nil {
			HandleError(w, err)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]string{
			"access_token":  tokens.AccessToken,
			"refresh_token": tokens.RefreshToken,
		})
	}
}
