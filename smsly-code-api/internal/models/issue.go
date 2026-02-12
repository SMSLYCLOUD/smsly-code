package models

import "time"

type IssueState string

const (
	IssueOpen   IssueState = "open"
	IssueClosed IssueState = "closed"
)

type Issue struct {
	ID          int64      `json:"id"`
	RepoID      int64      `json:"repo_id"`
	Number      int64      `json:"number"`
	Title       string     `json:"title"`
	Body        string     `json:"body"`
	State       IssueState `json:"state"`
	AuthorID    int64      `json:"author_id"`
	MilestoneID *int64     `json:"milestone_id,omitempty"`
	AssigneeID  *int64     `json:"assignee_id,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	ClosedAt    *time.Time `json:"closed_at,omitempty"`

	// Associations
	Repository *Repository `json:"repository,omitempty"`
	Author     *User       `json:"author,omitempty"`
	Assignee   *User       `json:"assignee,omitempty"`
	Milestone  *Milestone  `json:"milestone,omitempty"`
	Labels     []Label     `json:"labels,omitempty"`
}

type Comment struct {
	ID        int64     `json:"id"`
	IssueID   int64     `json:"issue_id"`
	AuthorID  int64     `json:"author_id"`
	Body      string    `json:"body"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	Author *User `json:"author,omitempty"`
}
