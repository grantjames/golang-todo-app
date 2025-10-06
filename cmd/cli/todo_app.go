package main

import (
	"bufio"
	"fmt"
	"log/slog"
	"os"
	"strings"
	"time"

	"grantjames.github.io/todo-app/types"
)

type TodoStore interface {
	GetTodo(id string) (types.Todo, error)
	AddTodo(todo types.Todo) (string, error)
	UpdateTodoStatus(id string, status types.Status) error
	GetTodosByStatus(status types.Status) map[string]types.Todo
	GetOverdueTodos() map[string]types.Todo
	GetAllTodos() map[string]types.Todo
}

type TodoApp struct {
	store  TodoStore
	logger *slog.Logger
}

func (app *TodoApp) Start() {
	app.logger.Debug("Application started")
	for {
		greeting := []string{
			"What do you want to do?",
			"1. Show todos",
			"2. Show archived/completed todos",
			"3. Show overdue todos",
			"4. Add a new todo",
			"5. Update a todo status",
			"6. Quit",
		}

		for _, t := range greeting {
			fmt.Println(t)
		}

		var input string

		_, err := fmt.Scanln(&input)
		if err != nil {
			fmt.Println("Invalid input.")
			continue
		}

		app.logger.Debug("Main menu option chosen", "option", input)

		switch input {
		case "1":
			fmt.Println("*** Your todos are ***")
			app.showTodos(app.store.GetAllTodos())
			fmt.Println("Press enter to continue...")
			fmt.Scanln()
		case "2":
			fmt.Println("*** Your archived/completed todos are ***")
			app.showTodos(app.store.GetTodosByStatus(types.Completed))
			fmt.Println("Press enter to continue...")
			fmt.Scanln()
		case "3":
			fmt.Println("*** Your overdue todos are ***")
			app.showTodos(app.store.GetOverdueTodos())
			fmt.Println("Press enter to continue...")
			fmt.Scanln()
		case "4":
			app.addNewTodo()
		case "5":
			app.updateTodo()
		case "6":
			fmt.Println("So sad to see you... go.")
			os.Exit(0)
		default:
			fmt.Println("Invalid input.")
			continue
		}
	}
}

func (t *TodoApp) addNewTodo() {
	scanner := bufio.NewReader(os.Stdin)
	var desc string

	for desc == "" {
		fmt.Print("Description: ")
		desc, _ = scanner.ReadString('\n')
		desc = strings.TrimSpace(desc)
	}

	due := t.readDate()

	t.store.AddTodo(types.NewTodo(desc, due))

	fmt.Println("Todo successfully added!")
}

// Returning a pointer to time.Time, even though the docs say typically you should pass by value
// because I want it to be an optional datetime.
func (t *TodoApp) readDate() *time.Time {
	scanner := bufio.NewReader(os.Stdin)
	var input string

	for {
		fmt.Print("Due: (yyyy-mm-dd format, leave blank for no due date) ")
		input, _ = scanner.ReadString('\n')
		input = strings.TrimSpace(input)

		if input == "" {
			return nil
		}

		parsedDueDate, err := time.Parse("2006-01-02", input)
		if err != nil {
			fmt.Println("Could not parse date:", err)
		} else {
			return &parsedDueDate
		}
	}

}

func (t *TodoApp) updateTodo() {
	scanner := bufio.NewReader(os.Stdin)
	var input string
	var id string
	var err error

	t.showTodos(t.store.GetAllTodos())
	for {
		fmt.Print("State the ID of the todo you wish to update: ")
		input, _ = scanner.ReadString('\n')
		input = strings.TrimSpace(input)
		id = input

		// Check the todo exists
		_, err = t.store.GetTodo(id)
		if err != nil {
			fmt.Println(err.Error())
			continue
		}

		break
	}

	for {
		fmt.Println("What status do you want to updated your todo to? (Started or Completed): ")
		input, _ = scanner.ReadString('\n')
		input = strings.TrimSpace(input)
		input = strings.ToUpper(string(input[0])) + input[1:]
		if input != "Started" && input != "Completed" {
			fmt.Println("Status should be Started or Completed")
			continue
		}

		err = t.store.UpdateTodoStatus(id, types.Status(input))
		if err != nil {
			fmt.Println(err.Error())
			continue
		} else {
			fmt.Println("Todo status updated")
			break
		}
	}
}

func (t *TodoApp) showTodos(todos map[string]types.Todo) {
	t.logger.Debug("Method called", "method", "showTodos(todos map[int]types.Todo)")

	for key, t := range todos {
		fmt.Printf("%s: ", key)
		fmt.Println(t.String())
	}
}
