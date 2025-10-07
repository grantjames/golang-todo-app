package types

type TodoStore interface {
	GetTodo(id string) (Todo, error)
	AddTodo(todo Todo) (string, error)
	UpdateTodoStatus(id string, status Status) error
	GetTodosByStatus(status Status) map[string]Todo
	GetOverdueTodos() map[string]Todo
	GetAllTodos() map[string]Todo
}
