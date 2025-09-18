package main

import (
	"bufio"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"strconv"
	"strings"
	"time"

	"grantjames.github.io/todo-app/stores"
	"grantjames.github.io/todo-app/types"
)

var store *stores.InMemoryTodoStore
var logger *slog.Logger

func init() {
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

	logger = slog.New(slog.NewTextHandler(f, opts))

	store = stores.NewInMemoryTodoStore()
}

func main() {
	logger.Debug("Application started")

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

		logger.Debug("Main menu option chosen", "option", input)

		switch input {
		case "1":
			fmt.Println("*** Your todos are ***")
			showTodos(store.GetAllTodos())
			fmt.Println("Press enter to continue...")
			fmt.Scanln()
		case "2":
			fmt.Println("*** Your archived/completed todos are ***")
			showTodos(store.GetTodosByStatus(types.Completed))
			fmt.Println("Press enter to continue...")
			fmt.Scanln()
		case "3":
			fmt.Println("*** Your overdue todos are ***")
			showTodos(store.GetOverdueTodos())
			fmt.Println("Press enter to continue...")
			fmt.Scanln()
		case "4":
			addNewTodo()
		case "5":
			UpdateTodo()
		case "6":
			fmt.Println("So sad to see you... go.")
			os.Exit(0)
		default:
			fmt.Println("Invalid input.")
			continue
		}
	}
}

func showTodos(todos map[int]types.Todo) {
	logger.Debug("Method called", "method", "showTodos(todos map[int]types.Todo)")

	for index, t := range todos {
		fmt.Printf("%d: ", index)
		fmt.Println(t.String())
	}
}

func addNewTodo() {
	scanner := bufio.NewReader(os.Stdin)
	var desc string

	for desc == "" {
		fmt.Print("Description: ")
		desc, _ = scanner.ReadString('\n')
		desc = strings.TrimSpace(desc)
	}

	due := readDate()

	store.AddTodo(types.NewTodo(desc, due))

	fmt.Println("Todo successfully added!")
}

// Returning a pointer to time.Time, even though the docs say typically you should pass by value
// because I want it to be an optional datetime.
func readDate() *time.Time {
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

func UpdateTodo() {
	scanner := bufio.NewReader(os.Stdin)
	var input string
	var id int
	var err error

	showTodos(store.GetAllTodos())
	for {
		fmt.Print("State the ID of the todo you wish to update: ")
		input, _ = scanner.ReadString('\n')
		input = strings.TrimSpace(input)

		id, err = strconv.Atoi(input)
		if err != nil {
			fmt.Println("ID should be an integer")
			continue
		}

		// Check the todo exists
		_, err = store.GetTodo(id)
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

		err = store.UpdateTodoStatus(id, types.Status(input))
		if err != nil {
			fmt.Println(err.Error())
			continue
		} else {
			fmt.Println("Todo status updated")
			break
		}
	}
}
