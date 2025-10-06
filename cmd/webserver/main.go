package main

import (
	"log"
	"net/http"

	todoapp "grantjames.github.io/todo-app"
	"grantjames.github.io/todo-app/stores"
)

func main() {
	server := todoapp.NewTodoServer(stores.NewInMemoryTodoStore())
	log.Fatal(http.ListenAndServe(":5000", server))
}
