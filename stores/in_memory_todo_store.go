package stores

import (
	"context"
	"fmt"
	"log/slog"
	"sync"

	"github.com/google/uuid"
	"grantjames.github.io/todo-app/types"
)

func NewInMemoryTodoStore() *InMemoryTodoStore {
	return &InMemoryTodoStore{
		map[string]types.Todo{},
		sync.RWMutex{},
	}
}

type InMemoryTodoStore struct {
	store map[string]types.Todo
	lock  sync.RWMutex
}

func (i *InMemoryTodoStore) GetTodo(ctx context.Context, id string) (types.Todo, error) {
	slog.InfoContext(ctx, "InMemoryTodoStore: GetTodo called", "todo_id", id)

	i.lock.RLock()
	defer i.lock.RUnlock()

	todo, ok := i.store[id]
	if !ok {
		return types.Todo{}, fmt.Errorf("no todo with id %s found", id)
	}
	return todo, nil
}

func (i *InMemoryTodoStore) AddTodo(ctx context.Context, todo types.Todo) (string, error) {
	slog.InfoContext(ctx, "InMemoryTodoStore: AddTodo called")

	i.lock.Lock()
	defer i.lock.Unlock()

	id := uuid.NewString()

	i.store[id] = todo

	return id, nil
}

func (i *InMemoryTodoStore) UpdateTodoStatus(ctx context.Context, id string, status types.Status) error {
	slog.InfoContext(ctx, "InMemoryTodoStore: UpdateTodoStatus called", "todo_id", id, "status", status)

	i.lock.Lock()
	defer i.lock.Unlock()

	todo, ok := i.store[id]
	if !ok {
		return fmt.Errorf("no todo with ID %s was found", id)
	}
	todo.SetStatus(status)
	i.store[id] = todo // Have to assign the value back since it's retrieved by value (not storing pointers)
	return nil
}

func (i *InMemoryTodoStore) GetTodosByStatus(ctx context.Context, status types.Status) map[string]types.Todo {
	slog.InfoContext(ctx, "InMemoryTodoStore: GetTodosByStatus called", "status", status)

	i.lock.RLock()
	defer i.lock.RUnlock()

	results := map[string]types.Todo{}
	for key, value := range i.store {
		if value.Status == status {
			results[key] = value
		}
	}
	return results
}

func (i *InMemoryTodoStore) GetOverdueTodos(ctx context.Context) map[string]types.Todo {
	slog.InfoContext(ctx, "InMemoryTodoStore: GetOverdueTodos called")

	i.lock.RLock()
	defer i.lock.RUnlock()

	results := map[string]types.Todo{}
	for key, t := range i.store {
		if t.IsOverdue() {
			results[key] = t
		}
	}
	return results
}

func (i *InMemoryTodoStore) GetAllTodos(ctx context.Context) map[string]types.Todo {
	slog.InfoContext(ctx, "InMemoryTodoStore: GetAllTodos called")

	i.lock.RLock()
	defer i.lock.RUnlock()

	results := map[string]types.Todo{}
	for key, t := range i.store {
		if t.Status != types.Completed {
			results[key] = t
		}
	}
	return results
}
