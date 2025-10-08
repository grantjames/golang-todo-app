package types

import "context"

type TodoStore interface {
	GetTodo(ctx context.Context, id string) (Todo, error)
	AddTodo(ctx context.Context, todo Todo) (string, error)
	UpdateTodoStatus(ctx context.Context, id string, status Status) error
	GetTodosByStatus(ctx context.Context, status Status) map[string]Todo
	GetOverdueTodos(ctx context.Context) map[string]Todo
	GetAllTodos(ctx context.Context) map[string]Todo
}
