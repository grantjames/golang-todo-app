package stores

import (
	"log/slog"
	"os"
	"testing"

	"grantjames.github.io/todo-app/types"
)

type MockLogger struct{}

func (m *MockLogger) Log(args ...any)   {}
func (m *MockLogger) Debug(args ...any) {}

func CreateTestStore() *InMemoryTodoStore {
	store := NewInMemoryTodoStore(slog.New(slog.NewTextHandler(os.Stderr, nil)))

	todo1 := types.NewTodo("Todo 1", nil)
	todo2 := types.NewTodo("Todo 2", nil)

	store.AddTodo(todo1)
	store.AddTodo(todo2)

	return store
}

func TestInMemoryStore(t *testing.T) {
	store := CreateTestStore()

	t.Run("Retrieve existing todo", func(t *testing.T) {
		todo, err := store.GetTodo(0)

		if err != nil {
			t.Fatalf("Expected to retrieve todo with ID 0, got error: %v", err)
		}
		if todo.Description() != "Todo 1" {
			t.Errorf("Expected description 'Todo 1', got '%s'", todo.Description())
		}
	})

	t.Run("Add new todo", func(t *testing.T) {
		newTodo := types.NewTodo("New Todo", nil)
		store.AddTodo(newTodo)

		todo, err := store.GetTodo(2)
		if err != nil {
			t.Fatalf("Expected to retrieve newly added todo with ID 2, got error: %v", err)
		}
		if todo.Description() != "New Todo" {
			t.Errorf("Expected description 'New Todo', got '%s'", todo.Description())
		}
	})

	t.Run("Update todo status", func(t *testing.T) {
		err := store.UpdateTodoStatus(1, types.Completed)
		if err != nil {
			t.Fatalf("Expected to update status of todo with ID 1, got error: %v", err)
		}

		todo, err := store.GetTodo(1)
		if err != nil {
			t.Fatalf("Expected to retrieve todo with ID 1, got error: %v", err)
		}
		if todo.Status() != types.Completed {
			t.Errorf("Expected status 'Completed', got '%s'", todo.Status())
		}
	})

	t.Run("Get todos by status", func(t *testing.T) {
		todos := store.GetTodosByStatus(types.NotStarted)
		if len(todos) != 2 {
			t.Errorf("Expected 2 not started todos, got %d", len(todos))
		}
	})

	t.Run("Get all todos", func(t *testing.T) {
		todos := store.GetAllTodos()
		if len(todos) != 2 {
			t.Errorf("Expected 2 todos (excluding completed), got %d", len(todos))
		}
	})

	t.Run("Retrieve non-existent todo", func(t *testing.T) {
		_, err := store.GetTodo(999999)
		if err == nil {
			t.Fatalf("Expected error when retrieving non-existent todo, got nil")
		}
	})
}
