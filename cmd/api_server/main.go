package main

import (
	"log"
	"log/slog"
	"net/http"
	"os"

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

	base := slog.NewTextHandler(f, nil)
	logger := slog.New(&traceIdContextHandler{h: base})
	slog.SetDefault(logger)

	server := todoapp.NewTodoAPIServer(stores.NewInMemoryTodoStore())
	log.Fatal(http.ListenAndServe(":5000", todoapp.LoggingMiddleware(server)))
}
