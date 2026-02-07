package service

import (
	"context"
	"log"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/ovk741/TasksStream/internal/domain"
	"github.com/ovk741/TasksStream/internal/storage/postgres"
)

func TestCreateBoard(t *testing.T) {
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

	generateID := func() string {
		return "board-1"
	}

	service := NewBoardService(boardRepo, columnRepo, taskRepo, boardMemberRepo, generateID)

	board, err := service.Create("1", "My board")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if board.ID != "board-1" {
		t.Errorf("expected id board-1, got %s", board.ID)
	}

	if board.Name != "My board" {
		t.Errorf("expected name 'My board', got %s", board.Name)
	}
}

func TestBoardServiceCreateInvalidInput(t *testing.T) {
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

	service := NewBoardService(boardRepo, columnRepo, taskRepo, boardMemberRepo, func() string {
		return "board-1"
	})

	_, err = service.Create("1", "")

	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if err != domain.ErrInvalidInput {
		t.Errorf("expected ErrInvalidInput, got %v", err)
	}
}

func TestBoardServiceGetAllEmpty(t *testing.T) {
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

	service := NewBoardService(boardRepo, columnRepo, taskRepo, boardMemberRepo, func() string {
		return "id"
	})

	boards, err := service.GetAll("1")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(boards) != 0 {
		t.Errorf("expected 0 boards, got %d", len(boards))
	}
}

func TestBoardServiceGetAllWithData(t *testing.T) {
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

	service := NewBoardService(boardRepo, columnRepo, taskRepo, boardMemberRepo, func() string {
		return "id"
	})

	_, _ = service.Create("1", "Board 1")
	_, _ = service.Create("1", "Board 2")

	boards, err := service.GetAll("1")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(boards) != 2 {
		t.Errorf("expected 2 boards, got %d", len(boards))
	}
}
