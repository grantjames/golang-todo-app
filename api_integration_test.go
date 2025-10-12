package todoapp

// import (
// 	"net/http"
// 	"net/http/httptest"
// 	"testing"

// 	"grantjames.github.io/todo-app/stores"
// )

// func TestAddingTodosAndRetrievingThem(t *testing.T) {
// 	store := stores.NewInMemoryTodoStore()
// 	server := NewTodoAPIServer(store)
// 	id := "no-existent-id"

// 	server.ServeHTTP(httptest.NewRecorder(), newPostTodoRequest())
// 	server.ServeHTTP(httptest.NewRecorder(), newPostTodoRequest())
// 	server.ServeHTTP(httptest.NewRecorder(), newPostTodoRequest())

// 	response := httptest.NewRecorder()
// 	server.ServeHTTP(response, newGetTodoRequest(id))
// 	assertStatus(t, response.Code, http.StatusNotFound)

// 	//assertResponseBody(t, response.Body.String(), "3")
// }
