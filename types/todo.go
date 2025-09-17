package types

import (
	"fmt"
	"time"
)

type Status string

const (
	NotStarted Status = "Not Started"
	Started    Status = "Started"
	Completed  Status = "Completed"
)

type Todo struct {
	description string
	status      Status
	due         *time.Time
	updated     time.Time
}

func NewTodo(desc string, due *time.Time) Todo {
	return Todo{
		description: desc,
		due:         due,
		status:      NotStarted,
		updated:     time.Now(),
	}
}

func (t *Todo) SetStatus(s Status) {
	t.status = s
	t.updated = time.Now()
}

//
// Getters to expose todo internals
//

func (t *Todo) Description() string {
	return t.description
}

func (t *Todo) Due() *time.Time {
	return t.due
}

func (t *Todo) Status() Status {
	return t.status
}

func (t *Todo) Updated() time.Time {
	return t.updated
}

func (t *Todo) String() string {
	due := "No due date set"
	if t.due != nil {
		due = t.due.Format("02/01/2006")
	}
	return fmt.Sprintf(`%s
  Status: %s
  Due: %s
  Updated: %s
	`, t.description, t.status, due, t.updated.Format("02/01/2006 at 15:04:05"))
}
