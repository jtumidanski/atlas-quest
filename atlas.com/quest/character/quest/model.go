package quest

import "time"

const (
	StatusUndefined  = "UNDEFINED"
	StatusNotStarted = "NOT_STARTED"
	StatusStarted    = "STARTED"
	StatusCompleted  = "COMPLETED"
)

type Model struct {
	id         uint16
	status     string
	completion time.Time
}

func (m Model) Id() uint16 {
	return m.id
}

func (m Model) Status() string {
	return m.status
}

func (m Model) Completion() time.Time {
	return m.completion
}
