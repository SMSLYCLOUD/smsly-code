package services

import (
	"context"
	"fmt"
	"time"

	"github.com/SMSLYCLOUD/smsly-code/smsly-code-api/internal/database"
	"github.com/SMSLYCLOUD/smsly-code/smsly-code-api/internal/models"
	"github.com/jackc/pgx/v5"
	"github.com/redis/go-redis/v9"
)

type IssueService struct {
	db    database.DBPool
	redis *redis.Client
}

func NewIssueService(db database.DBPool, redis *redis.Client) *IssueService {
	return &IssueService{db: db, redis: redis}
}

func (s *IssueService) Create(ctx context.Context, repoID int64, title, body string, authorID int64) (*models.Issue, error) {
	// Generate issue number (autoincrement per repo)
	// This is tricky with high concurrency. Usually requires a lock or a separate sequence table per repo.
	// For simplicity, we'll do a naive MAX(number) + 1 inside a transaction.

	tx, err := s.db.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback(ctx)

	var lastNumber int64
	err = tx.QueryRow(ctx, "SELECT COALESCE(MAX(number), 0) FROM issue WHERE repo_id = $1", repoID).Scan(&lastNumber)
	if err != nil {
		return nil, fmt.Errorf("get max number: %w", err)
	}
	newNumber := lastNumber + 1

	var issue models.Issue
	err = tx.QueryRow(ctx,
		`INSERT INTO issue (repo_id, number, title, body, author_id, state, created_at, updated_at)
		 VALUES ($1, $2, $3, $4, $5, $6, NOW(), NOW())
		 RETURNING id, repo_id, number, title, body, state, author_id, created_at, updated_at`,
		repoID, newNumber, title, body, authorID, models.IssueOpen,
	).Scan(
		&issue.ID, &issue.RepoID, &issue.Number, &issue.Title, &issue.Body, &issue.State, &issue.AuthorID, &issue.CreatedAt, &issue.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("insert issue: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("commit tx: %w", err)
	}

	return &issue, nil
}

func (s *IssueService) Get(ctx context.Context, repoID int64, number int64) (*models.Issue, error) {
	var issue models.Issue
	var milestoneID *int64
	var assigneeID *int64
	var closedAt *time.Time

	err := s.db.QueryRow(ctx,
		`SELECT id, repo_id, number, title, body, state, author_id, milestone_id, assignee_id, created_at, updated_at, closed_at
		 FROM issue
		 WHERE repo_id = $1 AND number = $2`,
		repoID, number,
	).Scan(
		&issue.ID, &issue.RepoID, &issue.Number, &issue.Title, &issue.Body, &issue.State, &issue.AuthorID, &milestoneID, &assigneeID, &issue.CreatedAt, &issue.UpdatedAt, &closedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil // Not found
		}
		return nil, fmt.Errorf("get issue: %w", err)
	}

	issue.MilestoneID = milestoneID
	issue.AssigneeID = assigneeID
	issue.ClosedAt = closedAt

	return &issue, nil
}

func (s *IssueService) List(ctx context.Context, repoID int64, page, perPage int) ([]models.Issue, int, error) {
	offset := (page - 1) * perPage

	rows, err := s.db.Query(ctx,
		`SELECT id, repo_id, number, title, body, state, author_id, created_at, updated_at
		 FROM issue
		 WHERE repo_id = $1
		 ORDER BY number DESC
		 LIMIT $2 OFFSET $3`,
		repoID, perPage, offset,
	)
	if err != nil {
		return nil, 0, fmt.Errorf("list issues: %w", err)
	}
	defer rows.Close()

	var issues []models.Issue
	for rows.Next() {
		var i models.Issue
		if err := rows.Scan(&i.ID, &i.RepoID, &i.Number, &i.Title, &i.Body, &i.State, &i.AuthorID, &i.CreatedAt, &i.UpdatedAt); err != nil {
			return nil, 0, fmt.Errorf("scan issue: %w", err)
		}
		issues = append(issues, i)
	}

	// Get total count
	var total int
	err = s.db.QueryRow(ctx, "SELECT COUNT(*) FROM issue WHERE repo_id = $1", repoID).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("count issues: %w", err)
	}

	return issues, total, nil
}

func (s *IssueService) Close(ctx context.Context, repoID, number int64) error {
	result, err := s.db.Exec(ctx,
		`UPDATE issue SET state = $1, closed_at = NOW(), updated_at = NOW()
		 WHERE repo_id = $2 AND number = $3`,
		models.IssueClosed, repoID, number,
	)
	if err != nil {
		return fmt.Errorf("close issue: %w", err)
	}
	if result.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}

func (s *IssueService) Reopen(ctx context.Context, repoID, number int64) error {
	result, err := s.db.Exec(ctx,
		`UPDATE issue SET state = $1, closed_at = NULL, updated_at = NOW()
		 WHERE repo_id = $2 AND number = $3`,
		models.IssueOpen, repoID, number,
	)
	if err != nil {
		return fmt.Errorf("reopen issue: %w", err)
	}
	if result.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}
