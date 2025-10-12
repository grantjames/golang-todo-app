package todoapp

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"grantjames.github.io/todo-app/stores"
)

func TestAddingTodosAndRetrievingThem(t *testing.T) {
	server := NewTodoServer(stores.NewTodoStoreActor(stores.NewInMemoryTodoStore()))
	id := "none-existent-id"

	server.ServeHTTP(httptest.NewRecorder(), newPostTodoRequest())
	server.ServeHTTP(httptest.NewRecorder(), newPostTodoRequest())
	server.ServeHTTP(httptest.NewRecorder(), newPostTodoRequest())

	response := httptest.NewRecorder()
	server.ServeHTTP(response, newGetTodoRequest(id))
	assertStatus(t, response.Code, http.StatusNotFound)
}
