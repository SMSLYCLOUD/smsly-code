package services

import (
	"context"
	"testing"
	"time"

	"github.com/SMSLYCLOUD/smsly-code/smsly-code-api/internal/models"
	"github.com/jackc/pgx/v5"
	"github.com/pashagolub/pgxmock/v4"
	"github.com/stretchr/testify/assert"
)

func TestIssueService_Create(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mock.Close()

	service := NewIssueService(mock, nil)

	ctx := context.Background()
	repoID := int64(1)
	authorID := int64(10)
	title := "Test Issue"
	body := "This is a test issue"

	// Expectation for Create
	mock.ExpectBegin()

	// Expect QueryRow for MAX(number)
	mock.ExpectQuery("SELECT COALESCE").
		WithArgs(repoID).
		WillReturnRows(pgxmock.NewRows([]string{"coalesce"}).AddRow(int64(0)))

	// Expect QueryRow for INSERT
	mock.ExpectQuery("INSERT INTO issue").
		WithArgs(repoID, int64(1), title, body, authorID, models.IssueOpen).
		WillReturnRows(pgxmock.NewRows([]string{
			"id", "repo_id", "number", "title", "body", "state", "author_id", "created_at", "updated_at",
		}).AddRow(int64(100), repoID, int64(1), title, body, models.IssueOpen, authorID, time.Now(), time.Now()))

	mock.ExpectCommit()

	issue, err := service.Create(ctx, repoID, title, body, authorID)

	assert.NoError(t, err)
	assert.NotNil(t, issue)
	assert.Equal(t, int64(100), issue.ID)
	assert.Equal(t, int64(1), issue.Number)
	assert.Equal(t, title, issue.Title)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestIssueService_Get(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mock.Close()

	service := NewIssueService(mock, nil)

	ctx := context.Background()
	repoID := int64(1)
	issueID := int64(100)
	number := int64(1)
	title := "Test Issue"
	authorID := int64(10)

	// Expect QueryRow
	mock.ExpectQuery("SELECT id, repo_id, number").
		WithArgs(repoID, number).
		WillReturnRows(pgxmock.NewRows([]string{
			"id", "repo_id", "number", "title", "body", "state", "author_id", "milestone_id", "assignee_id", "created_at", "updated_at", "closed_at",
		}).AddRow(issueID, repoID, number, title, "body", models.IssueOpen, authorID, nil, nil, time.Now(), time.Now(), nil))

	issue, err := service.Get(ctx, repoID, number)

	assert.NoError(t, err)
	assert.NotNil(t, issue)
	assert.Equal(t, issueID, issue.ID)
	assert.Equal(t, title, issue.Title)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestIssueService_Get_NotFound(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mock.Close()

	service := NewIssueService(mock, nil)

	ctx := context.Background()
	repoID := int64(1)
	number := int64(999)

	// Expect QueryRow returns NoRows
	mock.ExpectQuery("SELECT id, repo_id, number").
		WithArgs(repoID, number).
		WillReturnError(pgx.ErrNoRows)

	issue, err := service.Get(ctx, repoID, number)

	assert.NoError(t, err)
	assert.Nil(t, issue)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
