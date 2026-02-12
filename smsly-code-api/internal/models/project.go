package models

import "time"

type ProjectState string

const (
	ProjectOpen   ProjectState = "open"
	ProjectClosed ProjectState = "closed"
)

type Project struct {
	ID          int64        `json:"id"`
	OwnerID     int64        `json:"owner_id"` // User or Org
	Name        string       `json:"name"`
	Description string       `json:"description"`
	State       ProjectState `json:"state"`
	CreatedAt   time.Time    `json:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at"`

	Columns []ProjectColumn `json:"columns,omitempty"`
}

type ProjectColumn struct {
	ID        int64         `json:"id"`
	ProjectID int64         `json:"project_id"`
	Name      string        `json:"name"`
	Position  int           `json:"position"`
	Cards     []ProjectCard `json:"cards,omitempty"`
}

type ProjectCard struct {
	ID         int64  `json:"id"`
	ColumnID   int64  `json:"column_id"`
	ContentURL string `json:"content_url,omitempty"` // issue/pr link
	Note       string `json:"note,omitempty"`
	Position   int    `json:"position"`
}
