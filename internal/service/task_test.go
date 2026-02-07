package service

import (
	"context"
	"log"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/ovk741/TasksStream/internal/domain"
	"github.com/ovk741/TasksStream/internal/storage/postgres"
)

func TestCreateTask(t *testing.T) {
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

	column := domain.Column{
		ID:      "Column-1",
		Title:   "Column",
		BoardID: board.ID,
	}

	columnRepo.Create(column)

	service := NewTaskService(taskRepo, columnRepo, boardMemberRepo, func() string {
		return "task-1"
	})

	task, err := service.Create("1", "My task", "New", column.ID)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if task.ID != "task-1" {
		t.Errorf("expected id task-1, got %s", task.ID)
	}

	if task.Title != "My task" {
		t.Errorf("expected name 'My task', got %s", task.Title)
	}

}

func TestCreateTaskInvalidInput(t *testing.T) {
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

	column := domain.Column{
		ID:      "Column-1",
		Title:   "Column",
		BoardID: board.ID,
	}

	columnRepo.Create(column)

	service := NewTaskService(taskRepo, columnRepo, boardMemberRepo, func() string {
		return "task-1"
	})

	_, err = service.Create("1", "", "New", column.ID)

	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if err != domain.ErrInvalidInput {
		t.Errorf("expected ErrInvalidInput, got %v", err)
	}
}

func TestTaskServiceCreateColumnNotFound(t *testing.T) {
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
	column := domain.Column{
		ID:      "Column-1",
		Title:   "Column",
		BoardID: board.ID,
	}

	columnRepo.Create(column)

	service := NewTaskService(taskRepo, columnRepo, boardMemberRepo, func() string {
		return "task-1"
	})

	_, err = service.Create("1", "Column", "New", "unknown-column")

	if err != domain.ErrNotFound {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
}

func TestTaskServiceGetByColumnIDEmpty(t *testing.T) {
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
	column := domain.Column{
		ID:      "Column-1",
		Title:   "Column",
		BoardID: board.ID,
	}

	columnRepo.Create(column)

	service := NewTaskService(taskRepo, columnRepo, boardMemberRepo, func() string {
		return "task-1"
	})

	tasks, err := service.GetByColumnID("1", "Column-1")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(tasks) != 0 {
		t.Errorf("expected 0 columns, got %d", len(tasks))
	}
}

func TestTaskServiceGetByColumnIDWithData(t *testing.T) {
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
	column := domain.Column{
		ID:      "Column-1",
		Title:   "Column",
		BoardID: board.ID,
	}

	columnRepo.Create(column)

	service := NewTaskService(taskRepo, columnRepo, boardMemberRepo, func() string {
		return "task-1"
	})

	_, _ = service.Create("1", "Task 1", "New", "Column-1")
	_, _ = service.Create("1", "Task 2", "Old", "Column-1")

	tasks, err := service.GetByColumnID("1", "Column-1")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(tasks) != 2 {
		t.Errorf("expected 2 columns, got %d", len(tasks))
	}
}
