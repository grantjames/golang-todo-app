package main

import (
	"flag"
	"log"
	"log/slog"
	"net/http"
	"os"

	todoapp "grantjames.github.io/todo-app"
	"grantjames.github.io/todo-app/stores"
	"grantjames.github.io/todo-app/types"
)

const dbFileName = "db.json"

func main() {
	var storageFlag = flag.Int("f", 0, "Specify which data store to use. 0 = File, 1 = Memory. Default = 0")
	flag.Parse()

	base := slog.NewTextHandler(os.Stdout, nil)
	logger := slog.New(&traceIdContextHandler{h: base})
	slog.SetDefault(logger)

	var store types.TodoStore
	if *storageFlag == 0 {
		slog.Info("Using File Todo Store")

		db, err := os.OpenFile(dbFileName, os.O_RDWR|os.O_CREATE, 0666)

		if err != nil {
			log.Fatalf("problem opening %s %v", dbFileName, err)
		}

		store, err = stores.NewJSONFileTodoStore(db)
		if err != nil {
			slog.Error("problem creating file todo store, %v", "error", err.Error())
		}
	} else {
		slog.Info("Using Memory Todo Store")
		store = stores.NewInMemoryTodoStore()
	}

	a := stores.NewTodoStoreActor(store)
	server := todoapp.NewTodoServer(a)
	log.Fatal(http.ListenAndServe(":5000", todoapp.LoggingMiddleware(server)))
}
