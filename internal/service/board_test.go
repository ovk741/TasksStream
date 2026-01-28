package service

import (
	"testing"

	"github.com/ovk741/TasksStream/internal/storage/memory"
)

func TestCreateBoard(t *testing.T) {
	boardRepo := memory.NewBoardRepository()
	columnRepo := memory.NewColumnRepository()
	taskRepo := memory.NewTaskRepository()

	generateID := func() string {
		return "board-1"
	}

	service := NewBoardService(boardRepo, columnRepo, taskRepo, generateID)

	board, err := service.Create("My board")

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
	boardRepo := memory.NewBoardRepository()
	columnRepo := memory.NewColumnRepository()
	taskRepo := memory.NewTaskRepository()

	service := NewBoardService(boardRepo, columnRepo, taskRepo, func() string {
		return "board-1"
	})

	_, err := service.Create("")

	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if err != ErrInvalidInput {
		t.Errorf("expected ErrInvalidInput, got %v", err)
	}
}

func TestBoardServiceGetAllEmpty(t *testing.T) {
	boardRepo := memory.NewBoardRepository()
	columnRepo := memory.NewColumnRepository()
	taskRepo := memory.NewTaskRepository()

	service := NewBoardService(boardRepo, columnRepo, taskRepo, func() string {
		return "id"
	})

	boards, err := service.GetAll()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(boards) != 0 {
		t.Errorf("expected 0 boards, got %d", len(boards))
	}
}

func TestBoardServiceGetAllWithData(t *testing.T) {
	boardRepo := memory.NewBoardRepository()
	columnRepo := memory.NewColumnRepository()
	taskRepo := memory.NewTaskRepository()

	service := NewBoardService(boardRepo, columnRepo, taskRepo, func() string {
		return "id"
	})

	_, _ = service.Create("Board 1")
	_, _ = service.Create("Board 2")

	boards, err := service.GetAll()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(boards) != 2 {
		t.Errorf("expected 2 boards, got %d", len(boards))
	}
}
