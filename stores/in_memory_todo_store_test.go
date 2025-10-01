package stores

import (
	"testing"

	"grantjames.github.io/todo-app/types"
)

func TestStore(t *testing.T) {
	t.Run("Add and Get Todo", func(t *testing.T) {})

	store := NewInMemoryTodoStore(nil)

	todo1 := types.NewTodo("Test Todo 1", nil)
	todo2 := types.NewTodo("Test Todo 2", nil)

	store.AddTodo(todo1)
	store.AddTodo(todo2)

	retrievedTodo, err := store.GetTodo(0)
	if err != nil {
		t.Fatalf("Expected to retrieve todo with ID 0, got error: %v", err)
	}
	if retrievedTodo.Description != "Test Todo 1" {
		t.Errorf("Expected description 'Test Todo 1', got '%s'", retrievedTodo.Description)
	}
}
