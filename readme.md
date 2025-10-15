# Go Todo App

This is an example of a command line todo app, that reads todos from a REST API. The API can use an in-memory store or a JSON file store.

## Building the application

There are 2 executables. The CLI is in `cmd/cli` and the web server is in `cmd/server`.

To build them run `go build` in each folder, and then use `./cli` and `./server` in each respective folder.

The server will start on port 5000 and the CLI looks for the API at `http://localhost:5000`.

## Using the CLI

The CLI is a REPL (Read-Evaluate-Print-Loop). When it first starts, it asks the user to specify what they want to do. The user chooses their option by giving a number.

The following options are available:
* Show all todos (or just completed/archived, and overdue)
* Add a new todo
* Update a todo's status
* Quit the application

## Design Considerations

### Reading and writing todos
The API reads and returns JSON objects, so the CLI uses the `api_client.go` file to interact with it. This way the CLI isn't concerned with how the todos are represented or stored. This means another client could be created if the server served XML instead, leaving the internal workings of the CLI itself unchanged.

### Handling concurrency

Initilly, my solution used locks to ensure the stores could be read concurrently, but the final solution uses the Actor pattern. The Actor "owns" access to the store and communication is done with the actor via messages (using channels).