package stores

import (
	"context"
	"testing"

	"grantjames.github.io/todo-app/types"
)

var id1 string
var ctx = context.Background()

func CreateTestStore() *InMemoryTodoStore {
	store := NewInMemoryTodoStore()

	todo1 := types.NewTodo("Todo 1", nil)
	todo2 := types.NewTodo("Todo 2", nil)

	id1, _ = store.AddTodo(ctx, todo1)
	store.AddTodo(ctx, todo2)

	return store
}

func TestInMemoryStore(t *testing.T) {
	store := CreateTestStore()

	t.Run("Retrieve existing todo", func(t *testing.T) {
		todo, err := store.GetTodo(ctx, id1)

		if err != nil {
			t.Fatalf("Expected to retrieve todo with ID %s, got error: %v", id1, err)
		}
		if todo.Description != "Todo 1" {
			t.Errorf("Expected description 'Todo 1', got '%s'", todo.Description)
		}
	})

	t.Run("Add new todo", func(t *testing.T) {
		newTodo := types.NewTodo("New Todo", nil)
		id, _ := store.AddTodo(ctx, newTodo)

		todo, err := store.GetTodo(ctx, id)
		if err != nil {
			t.Fatalf("Expected to retrieve newly added todo with ID %s, got error: %v", id, err)
		}
		if todo.Description != "New Todo" {
			t.Errorf("Expected description 'New Todo', got '%s'", todo.Description)
		}
	})

	t.Run("Update todo status", func(t *testing.T) {
		err := store.UpdateTodoStatus(ctx, id1, types.Completed)
		if err != nil {
			t.Fatalf("Expected to update status of todo with ID %s, got error: %v", id1, err)
		}

		todo, err := store.GetTodo(ctx, id1)
		if err != nil {
			t.Fatalf("Expected to retrieve todo with ID %s, got error: %v", id1, err)
		}
		if todo.Status != types.Completed {
			t.Errorf("Expected status 'Completed', got '%s'", todo.Status)
		}
	})

	t.Run("Get todos by status", func(t *testing.T) {
		todos := store.GetTodosByStatus(ctx, types.NotStarted)
		if len(todos) != 2 {
			t.Errorf("Expected 2 not started todos, got %d", len(todos))
		}
	})

	t.Run("Get all todos", func(t *testing.T) {
		todos := store.GetAllTodos(ctx)
		if len(todos) != 2 {
			t.Errorf("Expected 2 todos (excluding completed), got %d", len(todos))
		}
	})

	t.Run("Retrieve non-existent todo", func(t *testing.T) {
		_, err := store.GetTodo(ctx, "non-existent-id")
		if err == nil {
			t.Fatalf("Expected error when retrieving non-existent todo, got nil")
		}
	})
}
