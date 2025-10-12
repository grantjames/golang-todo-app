package todoapp

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"grantjames.github.io/todo-app/stores"
	"grantjames.github.io/todo-app/types"
)

func TestGETTodos(t *testing.T) {
	store := StubTodoStore{
		map[string]types.Todo{
			"0": types.NewTodo("First todo", nil),
			"1": types.NewTodo("Second todo", nil),
		},
		[]types.Todo{},
		[]struct {
			id     string
			status types.Status
		}{},
		[]types.Status{},
		0,
		0,
	}
	server := NewTodoServer(stores.NewTodoStoreActor(&store))

	t.Run("Returns a todo", func(t *testing.T) {
		request := newGetTodoRequest("0")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatus(t, response.Code, http.StatusOK)

		var got types.Todo
		err := json.NewDecoder(response.Body).Decode(&got)

		if err != nil {
			t.Fatalf("Unable to parse response from server %q into a todo, '%v'", response.Body, err)
		}
	})

	t.Run("Returns another todo", func(t *testing.T) {
		request := newGetTodoRequest("1")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatus(t, response.Code, http.StatusOK)
	})

	t.Run("returns 404 on missing todo", func(t *testing.T) {
		request := newGetTodoRequest("non-existent-id")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatus(t, response.Code, http.StatusNotFound)
	})
}

func TestStoreTodos(t *testing.T) {
	store := StubTodoStore{
		map[string]types.Todo{},
		[]types.Todo{},
		[]struct {
			id     string
			status types.Status
		}{},
		[]types.Status{},
		0,
		0,
	}
	server := NewTodoServer(stores.NewTodoStoreActor(&store))

	t.Run("it adds a todo on POST", func(t *testing.T) {
		request := newPostTodoRequest()
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatus(t, response.Code, http.StatusAccepted)

		if len(store.addCalls) != 1 {
			t.Errorf("got %d calls to AddTodo want %d", len(store.addCalls), 1)
		}
	})
}

func TestPUTTodos(t *testing.T) {
	store := StubTodoStore{
		map[string]types.Todo{
			"stub-id": types.NewTodo("A todo", nil),
		},
		[]types.Todo{},
		[]struct {
			id     string
			status types.Status
		}{},
		[]types.Status{},
		0,
		0,
	}
	server := NewTodoServer(stores.NewTodoStoreActor(&store))

	t.Run("it updates a todo on PUT", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodPut, "/api/todos/stub-id", bytes.NewBuffer([]byte(`{"status":"Completed"}`)))
		response := httptest.NewRecorder()

		server.ServeHTTP(response, req)

		assertStatus(t, response.Code, http.StatusAccepted)

		if len(store.updateCalls) != 1 {
			t.Errorf("got %d calls to UpdateTodoStatus want %d", len(store.updateCalls), 1)
		}

		if store.updateCalls[0].id != "stub-id" {
			t.Errorf("got id %q want %q", store.updateCalls[0].id, "stub-id")
		}

		if store.updateCalls[0].status != types.Completed {
			t.Errorf("got status %q want %q", store.updateCalls[0].status, types.Completed)
		}
	})
}

func newGetTodoRequest(id string) *http.Request {
	req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/api/todos/%s", id), nil)
	return req
}

func newPostTodoRequest() *http.Request {
	req, _ := http.NewRequest(http.MethodPost, "/api/todos/", bytes.NewBuffer([]byte(`{"description":"New todo"}`)))
	return req
}

func assertStatus(t testing.TB, got, want int) {
	t.Helper()
	if got != want {
		t.Errorf("did not get correct status, got %d, want %d", got, want)
	}
}

type StubTodoStore struct {
	todos       map[string]types.Todo
	addCalls    []types.Todo
	updateCalls []struct {
		id     string
		status types.Status
	}
	statusCalls  []types.Status
	overdueCalls int
	allCalls     int
}

func (s *StubTodoStore) AddTodo(ctx context.Context, todo types.Todo) (string, error) {
	s.addCalls = append(s.addCalls, todo)
	return "stub-id", nil
}

func (s *StubTodoStore) GetTodo(ctx context.Context, id string) (types.Todo, error) {
	if todo, ok := s.todos[id]; ok {
		return todo, nil
	} else {
		return types.Todo{}, fmt.Errorf("todo not found")
	}
}

func (s *StubTodoStore) UpdateTodoStatus(ctx context.Context, id string, status types.Status) error {
	s.updateCalls = append(s.updateCalls, struct {
		id     string
		status types.Status
	}{id, status})
	return nil
}

func (s *StubTodoStore) GetTodosByStatus(ctx context.Context, status types.Status) map[string]types.Todo {
	s.statusCalls = append(s.statusCalls, status)
	return map[string]types.Todo{}
}

func (s *StubTodoStore) GetOverdueTodos(ctx context.Context) map[string]types.Todo {
	s.overdueCalls++
	return map[string]types.Todo{}
}

func (s *StubTodoStore) GetAllTodos(ctx context.Context) map[string]types.Todo {
	s.allCalls++
	return map[string]types.Todo{}
}
