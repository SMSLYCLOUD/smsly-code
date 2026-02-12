package services

import (
	"context"
	"fmt"

	"github.com/SMSLYCLOUD/smsly-code/smsly-code-api/internal/database"
	"github.com/SMSLYCLOUD/smsly-code/smsly-code-api/internal/models"
	"github.com/jackc/pgx/v5"
	"github.com/redis/go-redis/v9"
)

type ProjectService struct {
	db    database.DBPool
	redis *redis.Client
}

func NewProjectService(db database.DBPool, redis *redis.Client) *ProjectService {
	return &ProjectService{db: db, redis: redis}
}

func (s *ProjectService) CreateProject(ctx context.Context, ownerID int64, name, description string) (*models.Project, error) {
	var p models.Project
	err := s.db.QueryRow(ctx,
		`INSERT INTO project (owner_id, name, description, state, created_at, updated_at)
		 VALUES ($1, $2, $3, $4, NOW(), NOW())
		 RETURNING id, owner_id, name, description, state, created_at, updated_at`,
		ownerID, name, description, models.ProjectOpen,
	).Scan(
		&p.ID, &p.OwnerID, &p.Name, &p.Description, &p.State, &p.CreatedAt, &p.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("create project: %w", err)
	}
	return &p, nil
}

func (s *ProjectService) GetProject(ctx context.Context, id int64) (*models.Project, error) {
	var p models.Project
	err := s.db.QueryRow(ctx,
		`SELECT id, owner_id, name, description, state, created_at, updated_at
		 FROM project
		 WHERE id = $1`,
		id,
	).Scan(
		&p.ID, &p.OwnerID, &p.Name, &p.Description, &p.State, &p.CreatedAt, &p.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil // Not found
		}
		return nil, fmt.Errorf("get project: %w", err)
	}

	// Fetch columns
	cols, err := s.ListColumns(ctx, p.ID)
	if err != nil {
		return nil, fmt.Errorf("list columns: %w", err)
	}
	p.Columns = cols

	return &p, nil
}

func (s *ProjectService) ListColumns(ctx context.Context, projectID int64) ([]models.ProjectColumn, error) {
	rows, err := s.db.Query(ctx,
		`SELECT id, project_id, name, position
		 FROM project_column
		 WHERE project_id = $1
		 ORDER BY position ASC`,
		projectID,
	)
	if err != nil {
		return nil, fmt.Errorf("list columns: %w", err)
	}
	defer rows.Close()

	var cols []models.ProjectColumn
	for rows.Next() {
		var c models.ProjectColumn
		if err := rows.Scan(&c.ID, &c.ProjectID, &c.Name, &c.Position); err != nil {
			return nil, fmt.Errorf("scan column: %w", err)
		}
		// Fetch cards for each column
		cards, err := s.ListCards(ctx, c.ID)
		if err != nil {
			return nil, fmt.Errorf("list cards: %w", err)
		}
		c.Cards = cards
		cols = append(cols, c)
	}
	return cols, nil
}

func (s *ProjectService) ListCards(ctx context.Context, columnID int64) ([]models.ProjectCard, error) {
	rows, err := s.db.Query(ctx,
		`SELECT id, column_id, content_url, note, position
		 FROM project_card
		 WHERE column_id = $1
		 ORDER BY position ASC`,
		columnID,
	)
	if err != nil {
		return nil, fmt.Errorf("list cards: %w", err)
	}
	defer rows.Close()

	var cards []models.ProjectCard
	for rows.Next() {
		var c models.ProjectCard
		if err := rows.Scan(&c.ID, &c.ColumnID, &c.ContentURL, &c.Note, &c.Position); err != nil {
			return nil, fmt.Errorf("scan card: %w", err)
		}
		cards = append(cards, c)
	}
	return cards, nil
}

func (s *ProjectService) CreateColumn(ctx context.Context, projectID int64, name string) (*models.ProjectColumn, error) {
	var pos int
	// Get max position
	err := s.db.QueryRow(ctx, "SELECT COALESCE(MAX(position), 0) FROM project_column WHERE project_id = $1", projectID).Scan(&pos)
	if err != nil {
		return nil, fmt.Errorf("get max pos: %w", err)
	}
	newPos := pos + 1

	var c models.ProjectColumn
	err = s.db.QueryRow(ctx,
		`INSERT INTO project_column (project_id, name, position, created_at, updated_at)
		 VALUES ($1, $2, $3, NOW(), NOW())
		 RETURNING id, project_id, name, position`,
		projectID, name, newPos,
	).Scan(&c.ID, &c.ProjectID, &c.Name, &c.Position)
	if err != nil {
		return nil, fmt.Errorf("create column: %w", err)
	}
	return &c, nil
}

func (s *ProjectService) CreateCard(ctx context.Context, columnID int64, contentURL, note string) (*models.ProjectCard, error) {
	var pos int
	err := s.db.QueryRow(ctx, "SELECT COALESCE(MAX(position), 0) FROM project_card WHERE column_id = $1", columnID).Scan(&pos)
	if err != nil {
		return nil, fmt.Errorf("get max pos: %w", err)
	}
	newPos := pos + 1

	var c models.ProjectCard
	err = s.db.QueryRow(ctx,
		`INSERT INTO project_card (column_id, content_url, note, position, created_at, updated_at)
		 VALUES ($1, $2, $3, $4, NOW(), NOW())
		 RETURNING id, column_id, content_url, note, position`,
		columnID, contentURL, note, newPos,
	).Scan(&c.ID, &c.ColumnID, &c.ContentURL, &c.Note, &c.Position)
	if err != nil {
		return nil, fmt.Errorf("create card: %w", err)
	}
	return &c, nil
}
