package main

import (
	"net/http"
	"strconv"
	"time"

	httpapi "github.com/ovk741/TasksStream/internal/api/http"
	"github.com/ovk741/TasksStream/internal/service"
	"github.com/ovk741/TasksStream/internal/storage"
	"github.com/ovk741/TasksStream/internal/storage/memory"
)

func main() {
	var boardRepo storage.BoardRepository = memory.NewBoardRepository()
	var columnRepo storage.ColumnRepository = memory.NewColumnRepository()
	var taskRepo storage.TaskRepository = memory.NewTaskRepository()

	boardService := service.NewBoardService(boardRepo, columnRepo, taskRepo, generateID)
	columnService := service.NewColumnService(columnRepo, boardRepo, taskRepo, generateID)
	taskService := service.NewTaskService(taskRepo, columnRepo, generateID)

	http.HandleFunc("/boards", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {

		case http.MethodPost:

			handler := httpapi.CreateBoardHandler(boardService)
			handler(w, r)

		case http.MethodGet:
			handler := httpapi.GetBoardsHandler(boardService)
			handler(w, r)

		case http.MethodPut:
			handler := httpapi.UpdateBoardHandler(boardService)
			handler(w, r)

		case http.MethodDelete:
			handler := httpapi.DeleteBoardHandler(boardService)
			handler(w, r)

		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}

	})

	http.HandleFunc("/columns", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {

		case http.MethodPost:
			handler := httpapi.CreateColumnHandler(columnService)
			handler(w, r)

		case http.MethodGet:
			handler := httpapi.GetColumnsByBoardHandler(columnService)
			handler(w, r)
		case http.MethodPut:
			handler := httpapi.UpdateColumnHandler(columnService)
			handler(w, r)

		case http.MethodDelete:
			handler := httpapi.DeleteColumnHandler(columnService)
			handler(w, r)

		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/columns/move", httpapi.MoveColumnHandler(columnService))

	http.HandleFunc("/tasks", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {

		case http.MethodPost:
			handler := httpapi.CreateTaskHandler(taskService)
			handler(w, r)

		case http.MethodGet:
			handler := httpapi.GetTasksByColumnHandler(taskService)
			handler(w, r)

		case http.MethodPut:
			handler := httpapi.UpdateTaskHandler(taskService)
			handler(w, r)

		case http.MethodDelete:
			handler := httpapi.DeleteTaskHandler(taskService)
			handler(w, r)

		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/tasks/move", httpapi.MoveTaskHandler(taskService))

	http.ListenAndServe(":8080", nil)
}

func generateID() string {
	return strconv.FormatInt(time.Now().UnixNano(), 10)
}
