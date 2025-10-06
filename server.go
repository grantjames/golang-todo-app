package todoapp

import (
	"encoding/json"
	"net/http"
	"strings"

	"grantjames.github.io/todo-app/types"
)

type TodoStore interface {
	GetTodo(id string) (types.Todo, error)
	AddTodo(todo types.Todo) (string, error)
}

type TodoServer struct {
	store TodoStore
	http.Handler
}

func NewTodoServer(store TodoStore) *TodoServer {
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
		//s.AddTodo(w, r.Body)
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

func (s *TodoServer) AddTodo(w http.ResponseWriter, todo types.Todo) {
	s.store.AddTodo(types.NewTodo("New Todo", nil))
	w.WriteHeader(http.StatusAccepted)
}
