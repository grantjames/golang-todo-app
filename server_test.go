package todoapp

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"grantjames.github.io/todo-app/types"
)

func TestGETTodos(t *testing.T) {
	store := StubTodoStore{
		map[string]types.Todo{
			"0": types.NewTodo("First todo", nil),
			"1": types.NewTodo("Second todo", nil),
		},
		nil,
		sync.RWMutex{},
	}
	server := NewTodoServer(&store)

	t.Run("Returns a todo", func(t *testing.T) {
		request := newGetTodoRequest("0")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatus(t, response.Code, http.StatusOK)

		var got types.Todo

		err := json.NewDecoder(response.Body).Decode(&got)

		if err != nil {
			t.Fatalf("Unable to parse response from server %q into slice of Player, '%v'", response.Body, err)
		}
	})

	t.Run("Returns another todo", func(t *testing.T) {
		request := newGetTodoRequest("1")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatus(t, response.Code, http.StatusOK)
		//assertTodoResponseBody(t, response.Body.String(), "10")
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
		sync.RWMutex{},
	}
	server := NewTodoServer(&store)

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

// func TestLeague(t *testing.T) {
// 	store := StubTodoStore{}
// 	server := NewTodoServer(&store)

// 	t.Run("it returns 200 on /league", func(t *testing.T) {
// 		request, _ := http.NewRequest(http.MethodGet, "/league", nil)
// 		response := httptest.NewRecorder()

// 		server.ServeHTTP(response, request)

// 		assertStatus(t, response.Code, http.StatusOK)
// 	})
// }

func newGetTodoRequest(id string) *http.Request {
	req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/v1/todos/%s", id), nil)
	return req
}

func newPostTodoRequest() *http.Request {
	req, _ := http.NewRequest(http.MethodPost, "/v1/todos/", nil)
	return req
}

func assertStatus(t testing.TB, got, want int) {
	t.Helper()
	if got != want {
		t.Errorf("did not get correct status, got %d, want %d", got, want)
	}
}

// func assertTodoResponseBody(t testing.TB, got, want string) {
// 	t.Helper()
// 	if got != want {
// 		t.Errorf("response body is wrong, got %q want %q", got, want)
// 	}
// }

type StubTodoStore struct {
	todos    map[string]types.Todo
	addCalls []types.Todo
	lock     sync.RWMutex
}

func (s *StubTodoStore) AddTodo(todo types.Todo) (string, error) {
	s.lock.Lock()
	defer s.lock.Unlock()

	s.addCalls = append(s.addCalls, todo)
	return "stub-id", nil
}

func (s *StubTodoStore) GetTodo(id string) (types.Todo, error) {
	s.lock.RLock()
	defer s.lock.RUnlock()

	if todo, ok := s.todos[id]; ok {
		return todo, nil
	} else {
		return types.Todo{}, fmt.Errorf("todo not found")
	}
}
