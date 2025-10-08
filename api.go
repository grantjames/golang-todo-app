package todoapp

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"grantjames.github.io/todo-app/types"
)

type TraceIdKey struct{}

type TodoAPIServer struct {
	store types.TodoStore
	http.Handler
}

func NewTodoAPIServer(store types.TodoStore) *TodoAPIServer {
	s := new(TodoAPIServer)

	s.store = store

	router := http.NewServeMux()
	router.Handle("/v1/todos/", http.HandlerFunc(s.todosHandler))

	s.Handler = router

	return s
}

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), TraceIdKey{}, uuid.NewString())

		slog.InfoContext(ctx, "HTTP Request:", slog.String("method", r.Method), slog.String("path", r.URL.Path))
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (s *TodoAPIServer) todosHandler(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/v1/todos/")

	switch r.Method {
	case http.MethodGet:
		if id == "" {
			status := r.URL.Query().Get("status")
			if status != "" {
				s.GetTodosByStatus(w, r, types.Status(status))
				return
			}

			if r.URL.Query().Has("overdue") {
				s.GetOverdueTodos(w, r)
				return
			}

			s.GetAllTodos(w, r)
			return
		}
		s.GetTodo(w, r, id)
	case http.MethodPost:
		s.AddTodo(w, r)
	case http.MethodPut:
		s.UpdateTodoStatus(w, r, id)
	}
}

func (s *TodoAPIServer) GetTodo(w http.ResponseWriter, r *http.Request, id string) {
	logEndpointCall(r, "GetTodo", map[string]string{"todo_id": id})

	todo, err := s.store.GetTodo(r.Context(), id)

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(todo)
}

func (s *TodoAPIServer) AddTodo(w http.ResponseWriter, r *http.Request) {
	logEndpointCall(r, "AddTodo", nil)

	var todo types.Todo
	err := json.NewDecoder(r.Body).Decode(&todo)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusInternalServerError)
		return
	}

	id, _ := s.store.AddTodo(r.Context(), types.NewTodo(todo.Description, nil))
	w.WriteHeader(http.StatusAccepted)
	fmt.Fprintf(w, "%s", id)
}

func (s *TodoAPIServer) UpdateTodoStatus(w http.ResponseWriter, r *http.Request, id string) {
	logEndpointCall(r, "UpdateTodoStatus", map[string]string{"todo_id": id})

	var req struct {
		Status types.Status `json:"status"`
	}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		slog.Error("Failed to decode request body", slog.String("error", err.Error()))
		http.Error(w, "Failed to read request body", http.StatusInternalServerError)
		return
	}

	err = s.store.UpdateTodoStatus(r.Context(), id, req.Status)
	if err != nil {
		slog.Error("Failed to update todo status", slog.String("error", err.Error()))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusAccepted)
}

func (s *TodoAPIServer) GetTodosByStatus(w http.ResponseWriter, r *http.Request, status types.Status) {
	logEndpointCall(r, "GetTodosByStatus", map[string]string{"status": string(status)})

	todos := s.store.GetTodosByStatus(r.Context(), status)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(todos)
}

func (s *TodoAPIServer) GetOverdueTodos(w http.ResponseWriter, r *http.Request) {
	logEndpointCall(r, "GetOverdueTodos", nil)

	todos := s.store.GetOverdueTodos(r.Context())

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(todos)
}

func (s *TodoAPIServer) GetAllTodos(w http.ResponseWriter, r *http.Request) {
	logEndpointCall(r, "GetAllTodos", nil)
	todos := s.store.GetAllTodos(r.Context())

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(todos)
}

func logEndpointCall(r *http.Request, endpoint string, params map[string]string) {
	logParams := []any{slog.String("endpoint", endpoint)}
	for k, v := range params {
		logParams = append(logParams, slog.String(k, v))
	}
	slog.InfoContext(r.Context(), "Endpoint called", logParams...)
}
