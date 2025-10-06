package main

import (
	"log"
	"net/http"

	"grantjames.github.io/todo-app/stores"
)

func main() {
	server := NewTodoServer(stores.NewInMemoryTodoStore())
	log.Fatal(http.ListenAndServe(":5000", server))
}
