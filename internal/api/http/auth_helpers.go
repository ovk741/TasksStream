package httpapi

import (
	"net/http"

	"github.com/ovk741/TasksStream/internal/api/http/middleware"
)

func GetUserID(r *http.Request) (string, bool) {
	userID, ok := r.Context().Value(middleware.UserIDKey).(string)
	return userID, ok
}

func MustGetUserID(w http.ResponseWriter, r *http.Request) (string, bool) {
	userID, ok := GetUserID(r)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return "", false
	}
	return userID, true
}
