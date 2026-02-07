package httpapi

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/ovk741/TasksStream/internal/domain"
	"github.com/ovk741/TasksStream/internal/service"
	"github.com/ovk741/TasksStream/internal/storage/postgres"
)

func TestCreateBoardHandler(t *testing.T) {
	dsn := "postgres://user:password@localhost:5432/tasks?sslmode=disable"

	pool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	boardRepo := postgres.NewBoardRepository(pool)
	columnRepo := postgres.NewColumnRepository(pool)
	taskRepo := postgres.NewTaskRepository(pool)
	boardMemberRepo := postgres.NewBoardMemberRepository(pool)

	boardService := service.NewBoardService(boardRepo, columnRepo, taskRepo, boardMemberRepo, func() string {
		return "new-id"
	})

	handler := CreateBoardHandler(boardService)

	body := bytes.NewBufferString(`{"name":"My board"}`)

	req, err := http.NewRequest(http.MethodPost, "/boards", body)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d", http.StatusCreated, rr.Code)
	}

	var response domain.Board
	if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
		t.Fatal(err)
	}

	if response.ID != "new-id" {
		t.Errorf("expected id new-id, got %s", response.ID)
	}

	if response.Name != "My board" {
		t.Errorf("expected name 'My board', got %s", response.Name)
	}
}
