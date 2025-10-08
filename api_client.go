package todoapp

import (
	"encoding/json"
	"fmt"
	"net/http"

	"grantjames.github.io/todo-app/types"
)

func NewTodoAPIClient(apiBaseUrl string) *TodoAPIClient {
	return &TodoAPIClient{
		apiBaseUrl: apiBaseUrl,
		client:     &http.Client{},
	}
}

type TodoAPIClient struct {
	apiBaseUrl string
	client     *http.Client
}

func (c *TodoAPIClient) GetTodo(id string) (*types.Todo, error) {
	url := fmt.Sprintf("%s/todos/%s", c.apiBaseUrl, id)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get todo: status code %d", resp.StatusCode)
	}

	var todo types.Todo
	if err := json.NewDecoder(resp.Body).Decode(&todo); err != nil {
		return nil, err
	}

	return &todo, nil
}

// AddTodo(todo Todo) (string, error)
// UpdateTodoStatus(id string, status Status) error
// GetTodosByStatus(status Status) map[string]Todo
// GetOverdueTodos() map[string]Todo
// GetAllTodos() map[string]Todo
