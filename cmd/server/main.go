package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	httpapi "github.com/ovk741/TasksStream/internal/api/http"
	"github.com/ovk741/TasksStream/internal/api/http/middleware"
	"github.com/ovk741/TasksStream/internal/infra/auth"
	"github.com/ovk741/TasksStream/internal/infra/security"
	"github.com/ovk741/TasksStream/internal/service"
	"github.com/ovk741/TasksStream/internal/storage/postgres"
	"golang.org/x/crypto/bcrypt"
)

func main() {

	_ = godotenv.Load()

	dsn := os.Getenv("DATABASE_DSN")
	accessSecret := os.Getenv("ACCESS_SECRET")
	refreshSecret := os.Getenv("REFRESH_SECRET")

	accessTTLMinutes, _ := strconv.Atoi(os.Getenv("ACCESS_TTL_MINUTES"))
	refreshTTLHours, _ := strconv.Atoi(os.Getenv("REFRESH_TTL_HOURS"))

	accessTTL := time.Duration(accessTTLMinutes) * time.Minute
	refreshTTL := time.Duration(refreshTTLHours) * time.Hour

	pool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	boardRepo := postgres.NewBoardRepository(pool)
	columnRepo := postgres.NewColumnRepository(pool)
	taskRepo := postgres.NewTaskRepository(pool)
	userRepo := postgres.NewUserRepository(pool)
	boardMemberRepo := postgres.NewBoardMemberRepository(pool)

	jwtManager := auth.NewJWTManager(accessSecret, refreshSecret, accessTTL, refreshTTL)

	hasher := security.NewBcryptHasher(bcrypt.DefaultCost)

	authService := service.NewAuthService(userRepo, hasher, jwtManager, generateID)

	boardService := service.NewBoardService(boardRepo, columnRepo, taskRepo, boardMemberRepo, generateID)
	columnService := service.NewColumnService(columnRepo, boardRepo, boardMemberRepo, taskRepo, generateID)
	taskService := service.NewTaskService(taskRepo, columnRepo, boardMemberRepo, generateID)

	mux := http.NewServeMux()
	authMW := middleware.AuthMiddleware(jwtManager)

	mux.Handle("/auth/register", httpapi.RegisterHandler(authService))
	mux.Handle("/auth/login", httpapi.LoginHandler(authService))
	mux.Handle("/auth/refresh", httpapi.RefreshHandler(authService))

	mux.Handle("/boards", authMW(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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

	})))

	mux.Handle("/boards/invite", authMW(httpapi.InviteToBoardHandler(boardService)))
	mux.Handle("/boards/members", authMW(httpapi.GetBoardMembersHandler(boardService)))
	mux.Handle("/boards/members/remove", authMW(httpapi.RemoveBoardMemberHandler(boardService)))

	mux.Handle("/columns", authMW(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			httpapi.CreateColumnHandler(columnService)(w, r)
		case http.MethodGet:
			httpapi.GetColumnsByBoardHandler(columnService)(w, r)
		case http.MethodPut:
			httpapi.UpdateColumnHandler(columnService)(w, r)
		case http.MethodDelete:
			httpapi.DeleteColumnHandler(columnService)(w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})))

	mux.Handle("/columns/move", authMW(httpapi.MoveColumnHandler(columnService)))

	mux.Handle("/tasks", authMW(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			httpapi.CreateTaskHandler(taskService)(w, r)
		case http.MethodGet:
			httpapi.GetTasksByColumnHandler(taskService)(w, r)
		case http.MethodPut:
			httpapi.UpdateTaskHandler(taskService)(w, r)
		case http.MethodDelete:
			httpapi.DeleteTaskHandler(taskService)(w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})))

	mux.Handle("/tasks/move", authMW(httpapi.MoveTaskHandler(taskService)))

	port := os.Getenv("PORT")
	log.Fatal(http.ListenAndServe(":"+port, mux))
}

func generateID() string {
	return strconv.FormatInt(time.Now().UnixNano(), 10)
}
