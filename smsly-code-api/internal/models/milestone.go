package models

import "time"

type MilestoneState string

const (
	MilestoneOpen   MilestoneState = "open"
	MilestoneClosed MilestoneState = "closed"
)

type Milestone struct {
	ID          int64          `json:"id"`
	RepoID      int64          `json:"repo_id"`
	Title       string         `json:"title"`
	Description string         `json:"description"`
	State       MilestoneState `json:"state"`
	DueOn       *time.Time     `json:"due_on,omitempty"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
}
