package models

import "time"

type Repository struct {
	ID            int64     `json:"id"`
	OwnerID       int64     `json:"owner_id"`
	Name          string    `json:"name"`
	Description   string    `json:"description"`
	IsPrivate     bool      `json:"is_private"`
	IsFork        bool      `json:"is_fork"`
	DefaultBranch string    `json:"default_branch"`
	Stars         int       `json:"stars"`
	Forks         int       `json:"forks"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`

	Owner *User `json:"owner,omitempty"`
}
