package types

import (
	"testing"
	"time"
)

func TestNewTodo(t *testing.T) {
	t.Run("Creates a new todo with required properties", func(t *testing.T) {
		now := time.Now()
		want := Todo{
			description: "Test todo",
			status:      NotStarted,
			due:         &now,
			updated:     time.Now(),
		}

		got := NewTodo(want.description, want.due)

		if want.description != got.description || want.status != got.status || want.due != got.due {
			t.Errorf("got %v, want %v", got, want)
		}
	})
}

func TestTodoSetters(t *testing.T) {
	t.Run("SetStatus updates the status and updated timestamp", func(t *testing.T) {
		todo := NewTodo("Test todo", nil)
		initialUpdated := todo.updated

		time.Sleep(10 * time.Millisecond)
		newStatus := Started
		todo.SetStatus(newStatus)

		if todo.status != newStatus {
			t.Errorf("got status %v, want %v", todo.status, newStatus)
		}
		if !todo.updated.After(initialUpdated) {
			t.Errorf("updated timestamp was not updated")
		}
	})
}

func TestOverdue(t *testing.T) {
	t.Run("IsOverdue returns true for overdue todos", func(t *testing.T) {
		yesterday := time.Now().AddDate(0, 0, -1)
		todo := NewTodo("Overdue todo", &yesterday)

		if !todo.IsOverdue() {
			t.Errorf("expected todo to be overdue")
		}
	})

	t.Run("IsOverdue returns false for non-overdue todos", func(t *testing.T) {
		today := time.Now()
		todo := NewTodo("Not overdue todo", &today)

		if todo.IsOverdue() {
			t.Errorf("expected todo to not be overdue")
		}
	})

	t.Run("IsOverdue returns false for completed todos", func(t *testing.T) {
		yesterday := time.Now().AddDate(0, 0, -1)
		todo := NewTodo("Completed todo", &yesterday)
		todo.SetStatus(Completed)

		if todo.IsOverdue() {
			t.Errorf("expected completed todo to not be overdue")
		}
	})

	t.Run("IsOverdue returns false for todos without a due date", func(t *testing.T) {
		todo := NewTodo("No due date todo", nil)

		if todo.IsOverdue() {
			t.Errorf("expected todo without due date to not be overdue")
		}
	})
}
