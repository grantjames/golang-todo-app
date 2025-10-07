package todoapp

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"grantjames.github.io/todo-app/types"
)

type TodoAPIServer struct {
	store types.TodoStore
	http.Handler
}

func NewTodoAPIServer(store types.TodoStore) *TodoAPIServer {
	s := new(TodoAPIServer)

	s.store = store

	router := http.NewServeMux()
	//router.Handle("/v1/todos/", http.HandlerFunc(s.todosHandler))
	router.Handle("/v1/todos/", http.HandlerFunc(s.todosHandler))

	s.Handler = router

	return s
}

func (s *TodoAPIServer) todosHandler(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/v1/todos/")

	switch r.Method {
	case http.MethodGet:
		s.GetTodo(w, id)
	case http.MethodPost:
		s.AddTodo(w, r)
	case http.MethodPut:
		s.UpdateTodoStatus(w, r, id)
	}
}

func (s *TodoAPIServer) GetTodo(w http.ResponseWriter, id string) {
	todo, err := s.store.GetTodo(id)

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(todo)
}

func (s *TodoAPIServer) AddTodo(w http.ResponseWriter, r *http.Request) {
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

func (s *TodoAPIServer) UpdateTodoStatus(w http.ResponseWriter, r *http.Request, id string) {
	var req struct {
		Status types.Status `json:"status"`
	}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusInternalServerError)
		return
	}

	err = s.store.UpdateTodoStatus(id, req.Status)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusAccepted)
}

func (s *TodoAPIServer) GetTodosByStatus(w http.ResponseWriter, r *http.Request, status types.Status) {
	todos := s.store.GetTodosByStatus(status)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(todos)
}

func (s *TodoAPIServer) GetOverdueTodos(w http.ResponseWriter, r *http.Request) {
	todos := s.store.GetOverdueTodos()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(todos)
}

func (s *TodoAPIServer) GetAllTodos(w http.ResponseWriter, r *http.Request) {
	todos := s.store.GetAllTodos()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(todos)
}
