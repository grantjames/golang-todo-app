package todoapp

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"grantjames.github.io/todo-app/types"
)

type ServerTodoStore interface {
	GetTodo(id string) (types.Todo, error)
	AddTodo(todo types.Todo) (string, error)
}

type TodoServer struct {
	store ServerTodoStore
	http.Handler
}

func NewTodoServer(store ServerTodoStore) *TodoServer {
	s := new(TodoServer)

	s.store = store

	router := http.NewServeMux()
	router.Handle("/v1/todos/", http.HandlerFunc(s.todosHandler))

	s.Handler = router

	return s
}

func (s *TodoServer) todosHandler(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/v1/todos/")

	switch r.Method {
	case http.MethodGet:
		s.GetTodo(w, id)
	case http.MethodPost:
		s.AddTodo(w, r)
	}
}

func (s *TodoServer) GetTodo(w http.ResponseWriter, id string) {
	todo, err := s.store.GetTodo(id)

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(todo)
}

func (s *TodoServer) AddTodo(w http.ResponseWriter, r *http.Request) {
	var todo types.Todo
	err := json.NewDecoder(r.Body).Decode(&todo)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusInternalServerError)
		return
	}

	id, _ := s.store.AddTodo(types.NewTodo(todo.Description, nil))
	w.WriteHeader(http.StatusAccepted)
	fmt.Fprintf(w, "%s", id)
}
