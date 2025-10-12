package stores

import (
	"context"
	"log/slog"

	"grantjames.github.io/todo-app/types"
)

type TodoStoreActor struct {
	cmds  chan types.Cmd
	store types.TodoStore
}

func NewTodoStoreActor(store types.TodoStore) *TodoStoreActor {
	return &TodoStoreActor{
		cmds:  make(chan types.Cmd, 1),
		store: store,
	}
}

func (a *TodoStoreActor) Run(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case msg := <-a.cmds:
			switch m := msg.(type) {

			case types.GetTodoRequest:
				slog.InfoContext(ctx, "Actor received GetTodoRequest", slog.String("todo_id", m.Id))
				t, err := a.store.GetTodo(m.Ctx, m.Id)
				if err == nil {
					m.Resp <- types.GetTodoResponse{Todo: t}
				} else {
					m.Resp <- types.GetTodoResponse{Err: err}
				}

			case types.GetAllTodosRequest:
				slog.InfoContext(ctx, "Actor received GetAllTodosRequest")
				todos := a.store.GetAllTodos(m.Ctx)
				m.Resp <- types.GetAllTodosResponse{Todos: todos}

			case types.AddTodoRequest:
				slog.InfoContext(ctx, "Actor received AddTodoRequest")
				id, err := a.store.AddTodo(m.Ctx, m.Todo)
				if err == nil {
					m.Resp <- types.AddTodoResponse{Id: id}
				} else {
					m.Resp <- types.AddTodoResponse{Err: err}
				}

			case types.UpdateTodoStatusRequest:
				slog.InfoContext(ctx, "Actor received UpdateTodoStatusRequest")
				err := a.store.UpdateTodoStatus(m.Ctx, m.Id, m.Status)
				m.Resp <- types.UpdateTodoStatusResponse{Err: err}

			case types.GetOverDueTodosRequest:
				slog.InfoContext(ctx, "Actor received GetOverDueTodosRequest")
				todos := a.store.GetOverdueTodos(m.Ctx)
				m.Resp <- types.GetOverDueTodosResponse{Todos: todos}

			case types.GetTodosByStatusRequest:
				slog.InfoContext(ctx, "Actor received GetOverDueTodosRequest")
				todos := a.store.GetTodosByStatus(m.Ctx, m.Status)
				m.Resp <- types.GetTodosByStatusResponse{Todos: todos}
			}
		}
	}
}

func (a *TodoStoreActor) Send(cmd types.Cmd) {
	a.cmds <- cmd
}
