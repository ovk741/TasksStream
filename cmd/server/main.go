package main

import (
	"context"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	httpapi "github.com/ovk741/TasksStream/internal/api/http"
	"github.com/ovk741/TasksStream/internal/service"
	"github.com/ovk741/TasksStream/internal/storage/postgres"
)

func main() {

	dsn := "postgres://user:password@localhost:5432/tasks?sslmode=disable"

	pool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	boardRepo := postgres.NewBoardRepository(pool)
	columnRepo := postgres.NewColumnRepository(pool)
	taskRepo := postgres.NewTaskRepository(pool)

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
