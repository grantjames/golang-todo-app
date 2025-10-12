package todoapp

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"grantjames.github.io/todo-app/stores"
)

var actor = stores.NewTodoStoreActor(stores.NewInMemoryTodoStore())
var apiServer = NewTodoServer(actor)

func FuzzPOSTTodo(f *testing.F) {
	f.Add(`{"description": "test todo", "due_date": "2023-12-31T23:59:59Z"}`)
	f.Add(`{"description": "another test todo", "due_date": null}`)

	f.Fuzz(func(t *testing.T, todoText string) {
		req := httptest.NewRequest(http.MethodPost, "/api/todos/", strings.NewReader(todoText))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		apiServer.ServeHTTP(w, req)

		res := w.Result()

		if res.StatusCode != http.StatusAccepted && res.StatusCode != http.StatusBadRequest {
			t.Errorf("Unexpected status code: %d", res.StatusCode)
		}
	})
}
