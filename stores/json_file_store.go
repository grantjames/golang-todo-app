package stores

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"os"

	"github.com/google/uuid"
	"grantjames.github.io/todo-app/types"
)

type JSONFileTodoStore struct {
	database *json.Encoder
	todos    map[string]types.Todo
}

func initialiseTodosDBFile(file *os.File) error {
	file.Seek(0, io.SeekStart)

	info, err := file.Stat()

	if err != nil {
		return fmt.Errorf("problem getting file info from file %s, %v", file.Name(), err)
	}

	if info.Size() == 0 {
		file.Write([]byte("{}"))
		file.Seek(0, io.SeekStart)
	}

	return nil
}

func NewJSONFileTodoStore(file *os.File) (*JSONFileTodoStore, error) {
	err := initialiseTodosDBFile(file)

	if err != nil {
		return nil, fmt.Errorf("problem initialising todo db file, %v", err)
	}

	var todos map[string]types.Todo
	err = json.NewDecoder(file).Decode(&todos)

	if err != nil {
		return nil, fmt.Errorf("problem parsing todo file store, %v", err)
	}

	return &JSONFileTodoStore{
		database: json.NewEncoder(&tape{file}),
		todos:    todos,
	}, nil
}

func (i *JSONFileTodoStore) GetTodo(ctx context.Context, id string) (types.Todo, error) {
	slog.InfoContext(ctx, "JSONFileTodoStore: GetTodo called", "todo_id", id)

	todo, ok := i.todos[id]
	if !ok {
		return types.Todo{}, fmt.Errorf("no todo with id %s found", id)
	}
	return todo, nil
}

func (i *JSONFileTodoStore) AddTodo(ctx context.Context, todo types.Todo) (string, error) {
	slog.InfoContext(ctx, "JSONFileTodoStore: AddTodo called")

	id := uuid.NewString()

	i.todos[id] = todo

	i.database.Encode(i.todos)
	return id, nil
}

func (i *JSONFileTodoStore) UpdateTodoStatus(ctx context.Context, id string, status types.Status) error {
	slog.InfoContext(ctx, "JSONFileTodoStore: UpdateTodoStatus called", "todo_id", id, "status", status)

	todo, ok := i.todos[id]
	if !ok {
		return fmt.Errorf("no todo with ID %s was found", id)
	}
	todo.SetStatus(status)
	i.todos[id] = todo // Have to assign the value back since it's retrieved by value (not storing pointers)

	i.database.Encode(i.todos)
	return nil
}

func (i *JSONFileTodoStore) GetTodosByStatus(ctx context.Context, status types.Status) map[string]types.Todo {
	slog.InfoContext(ctx, "JSONFileTodoStore: GetTodosByStatus called", "status", status)

	results := map[string]types.Todo{}
	for key, value := range i.todos {
		if value.Status == status {
			results[key] = value
		}
	}
	return results
}

func (i *JSONFileTodoStore) GetOverdueTodos(ctx context.Context) map[string]types.Todo {
	slog.InfoContext(ctx, "JSONFileTodoStore: GetOverdueTodos called")

	results := map[string]types.Todo{}
	for key, t := range i.todos {
		if t.IsOverdue() {
			results[key] = t
		}
	}
	return results
}

func (i *JSONFileTodoStore) GetAllTodos(ctx context.Context) map[string]types.Todo {
	slog.InfoContext(ctx, "JSONFileTodoStore: GetAllTodos called")

	results := map[string]types.Todo{}
	for key, t := range i.todos {
		if t.Status != types.Completed {
			results[key] = t
		}
	}
	return results
}
