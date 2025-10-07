package main

import (
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/google/uuid"
	todoapp "grantjames.github.io/todo-app"
	"grantjames.github.io/todo-app/stores"
)

func main() {
	f, err := os.OpenFile("api.log",
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	logger := slog.New(slog.NewTextHandler(f, nil))

	server := todoapp.NewTodoAPIServer(stores.NewInMemoryTodoStore())
	log.Fatal(http.ListenAndServe(":5000", LoggingMiddleware(server, logger)))
}

func LoggingMiddleware(next http.Handler, logger *slog.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		trace_id := uuid.NewString()
		logger.Info("HTTP Request:", slog.String("trace_id", trace_id), slog.String("method", r.Method), slog.String("path", r.URL.Path))
		next.ServeHTTP(w, r)
	})
}
