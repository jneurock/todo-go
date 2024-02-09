package todo

import "time"

type Todo struct {
	Complete    bool
	ID          int64
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Description string
}
