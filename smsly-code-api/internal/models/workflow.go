package models

import (
	"time"

	"github.com/google/uuid"
)

type WorkflowRun struct {
	ID        uuid.UUID `json:"id" db:"id"`
	RepoID    int64     `json:"repo_id" db:"repo_id"`
	CommitSHA string    `json:"commit_sha" db:"commit_sha"`
	Branch    string    `json:"branch" db:"branch"`
	Status    string    `json:"status" db:"status"`
	Event     string    `json:"event" db:"event"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

type WorkflowJob struct {
	ID          uuid.UUID  `json:"id" db:"id"`
	RunID       uuid.UUID  `json:"run_id" db:"run_id"`
	Name        string     `json:"name" db:"name"`
	Status      string     `json:"status" db:"status"`
	StartedAt   *time.Time `json:"started_at" db:"started_at"`
	CompletedAt *time.Time `json:"completed_at" db:"completed_at"`
	LogsURL     string     `json:"logs_url" db:"logs_url"`
}

type WorkflowStep struct {
	ID          uuid.UUID  `json:"id" db:"id"`
	JobID       uuid.UUID  `json:"job_id" db:"job_id"`
	Name        string     `json:"name" db:"name"`
	Status      string     `json:"status" db:"status"`
	StartedAt   *time.Time `json:"started_at" db:"started_at"`
	CompletedAt *time.Time `json:"completed_at" db:"completed_at"`
}
