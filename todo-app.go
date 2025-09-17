package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"grantjames.github.io/todo-app/stores"
	"grantjames.github.io/todo-app/types"
)

var store = stores.NewInMemoryTodoStore()

func main() {
	time1, _ := time.Parse(time.DateOnly, "2025-10-01")
	store.AddTodo(types.NewTodo("First todo", &time1))
	store.AddTodo(types.NewTodo("Second todo", nil))

	for {
		greeting := []string{
			"What do you want to do?",
			"1. Show todos",
			"2. Add a new todo",
			"3. Update a todo status",
			"4. Quit",
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

		switch input {
		case "1":
			fmt.Println("*** Your todos are ***")
			showAllTodos()
			fmt.Println("Press enter to continue...")
			fmt.Scanln()
		case "2":
			addNewTodo()
		case "3":
			UpdateTodo()
		case "4":
			fmt.Println("So sad to see you... go.")
			os.Exit(0)
		default:
			fmt.Println("Invalid input.")
			continue
		}
	}
}

func showAllTodos() {
	for index, t := range store.GetAllTodos() {
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

	showAllTodos()
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
