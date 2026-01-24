package main

import (
	"net/http"
	"strconv"
	"time"

	httpapi "github.com/ovk741/TasksStream/internal/api/http"
	"github.com/ovk741/TasksStream/internal/storage/memory"
)

func main() {
	repo := memory.NewBoardRepository()
	columnRepo := memory.NewColumnRepository()
	taskRepo := memory.NewTaskRepository()

	http.HandleFunc("/boards", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {

		case http.MethodPost:
			handler := httpapi.CreateBoardHandler(repo, generateID)
			handler(w, r)

		case http.MethodGet:
			handler := httpapi.GetBoardsHandler(repo)
			handler(w, r)

		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/columns", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {

		case http.MethodPost:
			handler := httpapi.CreateColumnHandler(columnRepo, generateID)
			handler(w, r)

		case http.MethodGet:
			handler := httpapi.GetColumnsByBoardHandler(columnRepo)
			handler(w, r)

		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/tasks", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {

		case http.MethodPost:
			handler := httpapi.CreateTaskHandler(taskRepo, generateID)
			handler(w, r)

		case http.MethodGet:
			handler := httpapi.GetTasksByColumnHandler(taskRepo)
			handler(w, r)

		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	http.ListenAndServe(":8080", nil)
}

func generateID() string {
	return strconv.FormatInt(time.Now().UnixNano(), 10)
}
