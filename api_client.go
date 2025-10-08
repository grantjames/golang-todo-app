package todoapp

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

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

func (c *TodoAPIClient) AddTodo(todo types.Todo) (string, error) {
	url := fmt.Sprintf("%s/todos/", c.apiBaseUrl)
	todoData, err := json.Marshal(todo)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest(http.MethodPost, url, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Body = io.NopCloser(strings.NewReader(string(todoData)))

	resp, err := c.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return "", fmt.Errorf("failed to add todo: status code %d", resp.StatusCode)
	}

	var result struct {
		ID string `json:"id"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	return result.ID, nil
}

func (c *TodoAPIClient) UpdateTodoStatus(id string, status types.Status) error {
	url := fmt.Sprintf("%s/todos/%s", c.apiBaseUrl, id)
	updateData := struct {
		Status types.Status `json:"status"`
	}{
		Status: status,
	}
	data, err := json.Marshal(updateData)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPut, url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Body = io.NopCloser(strings.NewReader(string(data)))

	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to update todo status: status code %d", resp.StatusCode)
	}

	return nil
}

func (c *TodoAPIClient) GetTodosByStatus(status types.Status) (map[string]types.Todo, error) {
	url := fmt.Sprintf("%s/todos/?status=%s", c.apiBaseUrl, status)
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
		return nil, fmt.Errorf("failed to get todos by status: status code %d", resp.StatusCode)
	}

	var todos map[string]types.Todo
	if err := json.NewDecoder(resp.Body).Decode(&todos); err != nil {
		return nil, err
	}

	return todos, nil
}

func (c *TodoAPIClient) GetOverdueTodos() (map[string]types.Todo, error) {
	url := fmt.Sprintf("%s/todos/?overdue", c.apiBaseUrl)
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
		return nil, fmt.Errorf("failed to get overdue todos: status code %d", resp.StatusCode)
	}

	var todos map[string]types.Todo
	if err := json.NewDecoder(resp.Body).Decode(&todos); err != nil {
		return nil, err
	}

	return todos, nil
}

func (c *TodoAPIClient) GetAllTodos() (map[string]types.Todo, error) {
	url := fmt.Sprintf("%s/todos/", c.apiBaseUrl)
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
		return nil, fmt.Errorf("failed to get all todos: status code %d", resp.StatusCode)
	}

	var todos map[string]types.Todo
	if err := json.NewDecoder(resp.Body).Decode(&todos); err != nil {
		return nil, err
	}

	return todos, nil
}
