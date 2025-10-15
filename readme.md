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

### Logging
The server uses middleware to automatically add a trace ID to a context, and then uses a custom log handler that will add the trace ID from the context to any logs. This context is passed through the system so any calls from the API, to the server, through to the actor, and then underlying store can be linked via the trace ID. These logs are printed to `stdout`. A future improvement would be for the CLI to generate the trace ID and pass it via a header to the API. Then, the server could use this rather than generating its own.

The CLI prints logs to a file, and the log level can be set via an "l" flag when starting the application. Currently the application only prints messages at the `Info` level, so there's room for improvement here.

### Store
When starting the server, the default store to use is the file store. This ensures persistence between server restarts. However, the in memory store can be used by passing an "f" flag to the application with a value of 1.

### Tests
There are various tests demonstrating various techniques.

* `api_integrations_test.go` contains an integration test that calls the API, backed via the in-memory store, via the actor. It adds some todos and then verifies a 404 is returned when a non-existant ID is queried.
* `todo_store_actor_test.go` uses `t.Parallel()` to verify that the actor ensures safe concurrent read and write to the store by concurrently adding todos and then verifying that the number of todos returned are what was added.
* `server_test.go` tests adding and retrieving todos on the server. To ensure the server is tested in isolation, a "stub" todo store is created that verifies the server calls the expected methods on the store, without depending on a concrete implementation of the store.