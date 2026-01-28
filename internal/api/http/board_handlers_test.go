package httpapi

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ovk741/TasksStream/internal/domain"
	"github.com/ovk741/TasksStream/internal/service"
	"github.com/ovk741/TasksStream/internal/storage/memory"
)

func TestCreateBoardHandler(t *testing.T) {
	boardRepo := memory.NewBoardRepository()
	columnRepo := memory.NewColumnRepository()
	taskRepo := memory.NewTaskRepository()

	boardService := service.NewBoardService(boardRepo, columnRepo, taskRepo, func() string {
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
