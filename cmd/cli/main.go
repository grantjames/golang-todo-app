package main

import (
	"flag"
	"log/slog"
	"os"

	"grantjames.github.io/todo-app/stores"
)

func main() {
	var lFlag = flag.Int("l", 0, "Specify the logging level. DEBUG, INFO, WARN, ERROR")
	flag.Parse()

	f, err := os.OpenFile("app.log",
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	opts := &slog.HandlerOptions{
		Level: slog.Level(*lFlag),
	}

	logger := slog.New(slog.NewTextHandler(f, opts))

	app := TodoApp{
		stores.NewInMemoryTodoStore(),
		logger,
	}

	app.Start()
}
