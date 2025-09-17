package stores

import (
	"fmt"

	"grantjames.github.io/todo-app/types"
)

func NewInMemoryTodoStore() *InMemoryTodoStore {
	return &InMemoryTodoStore{map[int]types.Todo{}}
}

type InMemoryTodoStore struct {
	store map[int]types.Todo
}

func (i *InMemoryTodoStore) GetTodo(id int) (types.Todo, error) {
	todo, ok := i.store[id]
	if !ok {
		return types.Todo{}, fmt.Errorf("no todo with id %d found", id)
	}
	return todo, nil
}

func (i *InMemoryTodoStore) AddTodo(todo types.Todo) {
	i.store[len(i.store)] = todo
}

func (i *InMemoryTodoStore) UpdateTodoStatus(id int, status types.Status) error {
	todo, ok := i.store[id]
	if !ok {
		return fmt.Errorf("no todo with ID %d was found", id)
	}
	todo.SetStatus(status)
	i.store[id] = todo // Have to assign the value back since it's retrieved by value (not storing pointers)
	return nil
}

func (i *InMemoryTodoStore) GetTodosByStatus(status types.Status) map[int]types.Todo {
	results := map[int]types.Todo{}
	for key, value := range i.store {
		if value.Status() == status {
			results[key] = value
		}
	}
	return results
}

func (i *InMemoryTodoStore) GetAllTodos() map[int]types.Todo {
	return i.store
}
