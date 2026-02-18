package service

import (
	"context"
	"log"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/ovk741/TasksStream/internal/domain"
	"github.com/ovk741/TasksStream/internal/storage/postgres"
)

func TestCreateColumn(t *testing.T) {
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

	board := domain.Board{
		ID:   "board-1",
		Name: "Board",
	}
	boardRepo.Create(board)

	service := NewColumnService(columnRepo, boardRepo, boardMemberRepo, taskRepo, func() string {
		return "column-1"
	})

	column, err := service.Create("1", "My column", board.ID)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if column.ID != "column-1" {
		t.Errorf("expected id board-1, got %s", column.ID)
	}

	if column.Title != "My column" {
		t.Errorf("expected name 'My board', got %s", column.Title)
	}
}

func TestColumnServiceCreateInvalidInput(t *testing.T) {
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

	boardRepo.Create(domain.Board{
		ID: "board-1",
	})

	service := NewColumnService(columnRepo, boardRepo, boardMemberRepo, taskRepo, func() string {
		return "column-1"
	})

	_, err = service.Create("", "board-1", "1")

	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if err != domain.ErrInvalidInput {
		t.Errorf("expected ErrInvalidInput, got %v", err)
	}
}

func TestColumnServiceCreateBoardNotFound(t *testing.T) {
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

	service := NewColumnService(columnRepo, boardRepo, boardMemberRepo, taskRepo, func() string {
		return "column-1"
	})

	_, err = service.Create("1", "Column", "unknown-board")

	if err != domain.ErrNotFound {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
}

func TestColumnServiceGetByBoardIDEmpty(t *testing.T) {
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

	boardRepo.Create(domain.Board{
		ID: "board-1",
	})

	service := NewColumnService(columnRepo, boardRepo, boardMemberRepo, taskRepo, func() string {
		return "id"
	})

	columns, err := service.GetByBoardID("1", "board-1")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(columns) != 0 {
		t.Errorf("expected 0 columns, got %d", len(columns))
	}
}

func TestColumnServiceGetByBoardIDWithData(t *testing.T) {
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

	boardRepo.Create(domain.Board{
		ID: "board-1",
	})

	service := NewColumnService(columnRepo, boardRepo, boardMemberRepo, taskRepo, func() string {
		return "column-id"
	})

	_, _ = service.Create("1", "Column 1", "board-1")
	_, _ = service.Create("1", "Column 2", "board-1")

	columns, err := service.GetByBoardID("1", "board-1")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(columns) != 2 {
		t.Errorf("expected 2 columns, got %d", len(columns))
	}
}
