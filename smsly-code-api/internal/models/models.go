package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	Username  string         `gorm:"uniqueIndex;not null" json:"username"`
	Email     string         `gorm:"uniqueIndex;not null" json:"email"`
	Password  string         `json:"-"` // Never return password in JSON
}

type Repository struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
	Name        string         `gorm:"uniqueIndex;not null" json:"name"`
	Description string         `json:"description"`
	OwnerID     uint           `json:"owner_id"`
	Owner       User           `json:"owner"`
	IsPrivate   bool           `json:"is_private"`
}

type Issue struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
	RepoID      uint           `gorm:"index;not null" json:"repo_id"`
	Repository  Repository     `json:"-"`
	Title       string         `gorm:"not null" json:"title"`
	Body        string         `json:"body"`
	CreatorID   uint           `gorm:"index;not null" json:"creator_id"`
	Creator     User           `json:"creator"`
	AssigneeID  *uint          `gorm:"index" json:"assignee_id"`
	Assignee    *User          `json:"assignee"`
	State       string         `gorm:"default:'open'" json:"state"` // open, closed
}

type Comment struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	IssueID   uint           `gorm:"index;not null" json:"issue_id"`
	Issue     Issue          `json:"-"`
	UserID    uint           `gorm:"index;not null" json:"user_id"`
	User      User           `json:"user"`
	Body      string         `gorm:"not null" json:"body"`
}
