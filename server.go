package todoapp

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"strings"
	"text/template"

	"github.com/google/uuid"
	"grantjames.github.io/todo-app/stores"
	"grantjames.github.io/todo-app/types"
)

type TraceIdKey struct{}

type TodoServer struct {
	actor *stores.TodoStoreActor
	http.Handler
}

func NewTodoServer(actor *stores.TodoStoreActor) *TodoServer {
	s := new(TodoServer)

	s.actor = actor
	go s.actor.Run(context.Background())

	router := http.NewServeMux()
	router.Handle("/api/todos/", http.HandlerFunc(s.todosHandler))

	static := http.FileServer(http.Dir("./static"))
	router.Handle("/about/", http.StripPrefix("/about/", static))
	router.HandleFunc("/list", s.handleListPage)

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

func (s *TodoServer) handleListPage(w http.ResponseWriter, r *http.Request) {
	resp := make(chan types.GetAllTodosResponse)
	s.actor.Send(types.GetAllTodosRequest{Ctx: r.Context(), Resp: resp})

	tmpl, err := template.ParseFiles("templates/list.html")
	if err != nil {
		slog.InfoContext(r.Context(), "template parse error", "error", err.Error())
		http.Error(w, "template error", http.StatusInternalServerError)
		return
	}

	select {
	case res := <-resp:
		slog.InfoContext(r.Context(), "Received response from actor")
		if res.Err != nil {
			http.Error(w, res.Err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		if err := tmpl.Execute(w, res.Todos); err != nil {
			slog.InfoContext(r.Context(), "template execute error", "error", err.Error())
			http.Error(w, "template error", http.StatusInternalServerError)
			return
		}
	case <-r.Context().Done():
		http.Error(w, "request canceled", http.StatusRequestTimeout)
	}
}

func (s *TodoServer) todosHandler(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/api/todos/")

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

func (s *TodoServer) GetTodo(w http.ResponseWriter, r *http.Request, id string) {
	logEndpointCall(r, "GetTodo", map[string]string{"todo_id": id})

	resp := make(chan types.GetTodoResponse)
	s.actor.Send(types.GetTodoRequest{Ctx: r.Context(), Id: id, Resp: resp})

	select {
	case res := <-resp:
		slog.InfoContext(r.Context(), "Received response from actor", slog.String("todo_id", id))
		if res.Err != nil {
			http.Error(w, res.Err.Error(), http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(res.Todo)
	case <-r.Context().Done():
		http.Error(w, "request canceled", http.StatusRequestTimeout)
	}
}

func (s *TodoServer) AddTodo(w http.ResponseWriter, r *http.Request) {
	logEndpointCall(r, "AddTodo", nil)

	var todo types.Todo
	err := json.NewDecoder(r.Body).Decode(&todo)
	if err != nil {
		slog.ErrorContext(r.Context(), "Failed to decode request body", slog.String("error", err.Error()))
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}

	resp := make(chan types.AddTodoResponse)
	s.actor.Send(types.AddTodoRequest{Ctx: r.Context(), Todo: todo, Resp: resp})

	select {
	case res := <-resp:
		slog.InfoContext(r.Context(), "Received response from actor", slog.String("todo_id", res.Id))
		if res.Err != nil {
			http.Error(w, res.Err.Error(), http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusAccepted)
	case <-r.Context().Done():
		http.Error(w, "request canceled", http.StatusRequestTimeout)
	}
}

func (s *TodoServer) UpdateTodoStatus(w http.ResponseWriter, r *http.Request, id string) {
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

	resp := make(chan types.UpdateTodoStatusResponse)
	s.actor.Send(types.UpdateTodoStatusRequest{Ctx: r.Context(), Id: id, Status: req.Status, Resp: resp})

	select {
	case res := <-resp:
		if res.Err != nil {
			http.Error(w, res.Err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusAccepted)
	case <-r.Context().Done():
		http.Error(w, "request canceled", http.StatusRequestTimeout)
	}
}

func (s *TodoServer) GetTodosByStatus(w http.ResponseWriter, r *http.Request, status types.Status) {
	logEndpointCall(r, "GetTodosByStatus", map[string]string{"status": string(status)})

	resp := make(chan types.GetTodosByStatusResponse)
	s.actor.Send(types.GetTodosByStatusRequest{Ctx: r.Context(), Status: status, Resp: resp})

	select {
	case res := <-resp:
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(res.Todos)
	case <-r.Context().Done():
		http.Error(w, "request canceled", http.StatusRequestTimeout)
	}
}

func (s *TodoServer) GetOverdueTodos(w http.ResponseWriter, r *http.Request) {
	logEndpointCall(r, "GetOverdueTodos", nil)

	resp := make(chan types.GetOverDueTodosResponse)
	s.actor.Send(types.GetOverDueTodosRequest{Ctx: r.Context(), Resp: resp})

	select {
	case res := <-resp:
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(res.Todos)
	case <-r.Context().Done():
		http.Error(w, "request canceled", http.StatusRequestTimeout)
	}
}

func (s *TodoServer) GetAllTodos(w http.ResponseWriter, r *http.Request) {
	logEndpointCall(r, "GetAllTodos", nil)

	resp := make(chan types.GetAllTodosResponse)
	s.actor.Send(types.GetAllTodosRequest{Ctx: r.Context(), Resp: resp})

	select {
	case res := <-resp:
		if res.Err != nil {
			http.Error(w, res.Err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(res.Todos)
	case <-r.Context().Done():
		http.Error(w, "request canceled", http.StatusRequestTimeout)
	}
}

func logEndpointCall(r *http.Request, endpoint string, params map[string]string) {
	logParams := []any{slog.String("endpoint", endpoint)}
	for k, v := range params {
		logParams = append(logParams, slog.String(k, v))
	}
	slog.InfoContext(r.Context(), "Endpoint called", logParams...)
}
