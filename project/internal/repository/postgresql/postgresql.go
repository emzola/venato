package postgresql

import (
	"context"
	"database/sql"
	"os"

	"github.com/emzola/venato/project/pkg/model"
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

// Create adds a project record.
func (r *Repository) Create(ctx context.Context, project *model.Project) error {
	query := `
			INSERT INTO project (name, description, start_date, target_end_date, created_by)
			VALUES ($1, $2, $3, $4)
		  	RETURNING id, created_on, version`
	args := []interface{}{project.Name, project.Description, project.StartDate, project.TargetEndDate, project.CreatedBy}
	return r.db.QueryRowContext(ctx, query, args...).Scan(&project.ID, &project.CreatedOn, &project.Version)
}
