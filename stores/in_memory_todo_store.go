package stores

import (
	"fmt"
	"log/slog"
	"sync"

	"grantjames.github.io/todo-app/types"
)

func NewInMemoryTodoStore(logger *slog.Logger) *InMemoryTodoStore {
	return &InMemoryTodoStore{
		map[int]types.Todo{},
		logger,
		sync.RWMutex{},
	}
}

type InMemoryTodoStore struct {
	store  map[int]types.Todo
	logger *slog.Logger
	lock   sync.RWMutex
}

func (i *InMemoryTodoStore) GetTodo(id int) (types.Todo, error) {
	i.logger.Debug("Retrieving todo from in memory store")

	i.lock.RLock()
	defer i.lock.RUnlock()

	todo, ok := i.store[id]
	if !ok {
		return types.Todo{}, fmt.Errorf("no todo with id %d found", id)
	}
	return todo, nil
}

func (i *InMemoryTodoStore) AddTodo(todo types.Todo) {
	i.lock.Lock()
	defer i.lock.Unlock()

	i.store[len(i.store)] = todo
}

func (i *InMemoryTodoStore) UpdateTodoStatus(id int, status types.Status) error {
	i.lock.Lock()
	defer i.lock.Unlock()

	todo, ok := i.store[id]
	if !ok {
		return fmt.Errorf("no todo with ID %d was found", id)
	}
	todo.SetStatus(status)
	i.store[id] = todo // Have to assign the value back since it's retrieved by value (not storing pointers)
	return nil
}

func (i *InMemoryTodoStore) GetTodosByStatus(status types.Status) map[int]types.Todo {
	i.lock.RLock()
	defer i.lock.RUnlock()

	results := map[int]types.Todo{}
	for key, value := range i.store {
		if value.Status() == status {
			results[key] = value
		}
	}
	return results
}

func (i *InMemoryTodoStore) GetOverdueTodos() map[int]types.Todo {
	i.lock.RLock()
	defer i.lock.RUnlock()

	results := map[int]types.Todo{}
	for key, t := range i.store {
		if t.IsOverdue() {
			results[key] = t
		}
	}
	return results
}

func (i *InMemoryTodoStore) GetAllTodos() map[int]types.Todo {
	i.lock.RLock()
	defer i.lock.RUnlock()

	results := map[int]types.Todo{}
	for key, t := range i.store {
		if t.Status() != types.Completed {
			results[key] = t
		}
	}
	return results
}
