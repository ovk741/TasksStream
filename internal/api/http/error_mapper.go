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

	default:

		log.Printf("internal error: %v", err)

		SendError(
			w,
			http.StatusInternalServerError,
			errors.New("internal server error"),
		)
	}
}
