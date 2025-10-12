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
	Description string     `json:"description"`
	Status      Status     `json:"status"`
	Due         *time.Time `json:"due"`
	Updated     time.Time  `json:"updated"`
}

func NewTodo(desc string, due *time.Time) Todo {
	return Todo{
		Description: desc,
		Due:         due,
		Status:      NotStarted,
		Updated:     time.Now(),
	}
}

func (t *Todo) SetStatus(s Status) {
	t.Status = s
	t.Updated = time.Now()
}

func (t *Todo) IsOverdue() bool {
	today := time.Now().Truncate(24 * time.Hour)
	return t.Due != nil && t.Due.Before(today) && t.Status != Completed
}

func (t *Todo) String() string {
	due := "No due date set"
	if t.Due != nil {
		due = t.Due.Format("02/01/2006")
	}
	return fmt.Sprintf(`%s
  Status: %s
  Due: %s
  Updated: %s
	`, t.Description, t.Status, due, t.Updated.Format("02/01/2006 at 15:04:05"))
}
