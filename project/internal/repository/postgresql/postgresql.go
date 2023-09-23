package postgresql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"

	"github.com/emzola/venato/project/internal/repository"
	"github.com/emzola/venato/project/pkg/model"
	_ "github.com/lib/pq"
)

// Repository defines a PostgreSQL-based project repository.
type Repository struct {
	db *sql.DB
}

// New creates a new PostgreSQL-based repository.
func New() (*Repository, error) {
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		return nil, err
	}
	return &Repository{db}, nil
}

// Create adds a new project record.
func (r *Repository) Create(ctx context.Context, project *model.Project) error {
	query := `
			INSERT INTO project (name, description, start_date, target_end_date, created_by, modified_by)
			VALUES ($1, $2, $3, $4, $5, $6)
		  	RETURNING id, created_on, version`
	args := []interface{}{project.Name, project.Description, project.StartDate, project.TargetEndDate, project.CreatedBy, project.ModifiedBy}
	return r.db.QueryRowContext(ctx, query, args...).Scan(&project.ID, &project.CreatedOn, &project.Version)
}

// Get retrieves a project record by its id.
func (r *Repository) Get(ctx context.Context, id int64) (*model.Project, error) {
	if id < 1 {
		return nil, repository.ErrNotFound
	}
	query := `
		SELECT id, name, description, start_date, target_end_date, actual_end_date, created_on, created_by, modified_on, modified_by, version
		FROM project
		WHERE id = $1`
	var project model.Project
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&project.ID,
		&project.Name,
		&project.Description,
		&project.StartDate,
		&project.TargetEndDate,
		&project.ActualEndDate,
		&project.CreatedOn,
		&project.CreatedBy,
		&project.ModifiedOn,
		&project.ModifiedBy,
		&project.Version,
	)
	if err != nil {
		switch {
		case err.Error() == "pq: canceling statement due to user request":
			return nil, fmt.Errorf("%v: %w", err, ctx.Err())
		case errors.Is(err, sql.ErrNoRows):
			return nil, repository.ErrNotFound
		default:
			return nil, err
		}
	}
	return &project, nil
}

// Update updates a project record.
func (r *Repository) Update(ctx context.Context, project *model.Project) error {
	query := `
		UPDATE project
		SET name = $1, description = $2, start_date = $3, target_end_date = $4, actual_end_date = $5, modified_on = CURRENT_TIMESTAMP(0), modified_by = $6, version = version + 1
		WHERE id = $7 AND version = $8
		RETURNING version`
	args := []interface{}{project.Name, project.Description, project.StartDate, project.TargetEndDate, project.ActualEndDate, project.ModifiedBy, project.ID, project.Version}
	err := r.db.QueryRowContext(ctx, query, args...).Scan(&project.Version)
	if err != nil {
		switch {
		case err.Error() == "pq: canceling statement due to user request":
			return fmt.Errorf("%v: %w", err, ctx.Err())
		case errors.Is(err, sql.ErrNoRows):
			return repository.ErrEditConflict
		default:
			return err
		}
	}
	return nil
}

// Delete removes a project record by its id.
func (r *Repository) Delete(ctx context.Context, id int64) error {
	if id < 1 {
		return repository.ErrNotFound
	}
	query := `
		DELETE FROM project
		WHERE id = $1`
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		switch {
		case err.Error() == "pq: canceling statement due to user request":
			return fmt.Errorf("%v: %w", err, ctx.Err())
		default:
			return err
		}
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return repository.ErrNotFound
	}
	return nil
}
