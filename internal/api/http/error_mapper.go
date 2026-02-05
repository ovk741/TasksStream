package httpapi

import (
	"errors"
	"log"
	"net/http"

	"github.com/ovk741/TasksStream/internal/domain"
)

func HandleError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, domain.ErrInvalidInput):
		SendError(w, http.StatusBadRequest, err)

	case errors.Is(err, domain.ErrNotFound):
		SendError(w, http.StatusNotFound, err)

	case errors.Is(err, domain.ErrInvalidCredentials):
		SendError(w, http.StatusUnauthorized, err)

	case errors.Is(err, domain.ErrForbidden):
		SendError(w, http.StatusForbidden, err)

	case errors.Is(err, domain.ErrUserAlreadyExists),
		errors.Is(err, domain.ErrConflict):
		SendError(w, http.StatusConflict, err)

	default:

		log.Printf("internal error: %v", err)

		SendError(
			w,
			http.StatusInternalServerError,
			errors.New("internal server error"),
		)
	}
}
