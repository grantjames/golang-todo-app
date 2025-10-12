package types

import (
	"context"
)

type Cmd interface{ isCmd() }

type GetTodoRequest struct {
	Ctx  context.Context
	Id   string
	Resp chan GetTodoResponse
}

func (GetTodoRequest) isCmd() {}

type GetTodoResponse struct {
	Todo Todo
	Err  error
}

type GetAllTodosRequest struct {
	Ctx  context.Context
	Resp chan GetAllTodosResponse
}

func (GetAllTodosRequest) isCmd() {}

type GetAllTodosResponse struct {
	Todos map[string]Todo
	Err   error
}

type AddTodoRequest struct {
	Ctx  context.Context
	Todo Todo
	Resp chan AddTodoResponse
}

func (AddTodoRequest) isCmd() {}

type AddTodoResponse struct {
	Id  string
	Err error
}

type UpdateTodoStatusRequest struct {
	Ctx    context.Context
	Id     string
	Status Status
	Resp   chan UpdateTodoStatusResponse
}

func (UpdateTodoStatusRequest) isCmd() {}

type UpdateTodoStatusResponse struct {
	Err error
}

type GetTodosByStatusRequest struct {
	Ctx    context.Context
	Status Status
	Resp   chan GetTodosByStatusResponse
}

func (GetTodosByStatusRequest) isCmd() {}

type GetTodosByStatusResponse struct {
	Todos map[string]Todo
}

type GetOverDueTodosRequest struct {
	Ctx  context.Context
	Resp chan GetOverDueTodosResponse
}

func (GetOverDueTodosRequest) isCmd() {}

type GetOverDueTodosResponse struct {
	Todos map[string]Todo
}
