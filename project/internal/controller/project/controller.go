package project

import (
	"context"
	"time"

	"github.com/emzola/venato/project/internal/controller"
	"github.com/emzola/venato/project/pkg/model"
	"github.com/emzola/venato/project/pkg/validator"
)

type projectRepository interface {
	Create(ctx context.Context, project *model.Project) error
	// Get(ctx context.Context, id int64) (*model.Project, error)
	// Update(ctx context.Context, project *model.Project) error
	// Delete(ctx context.Context, id int64) error
}

// Controller defines a new project service controller.
type Controller struct {
	repo projectRepository
}

// New creates a project service controller.
func New(repo projectRepository) *Controller {
	return &Controller{repo}
}

// Create creates a new project.
func (c *Controller) Create(ctx context.Context, name string, description string, startDate time.Time, targetEndDate time.Time, createdBy int64) (*model.Project, error) {
	project := &model.Project{
		Name:          name,
		Description:   description,
		StartDate:     startDate,
		TargetEndDate: targetEndDate,
		CreatedBy:     createdBy,
	}
	v := validator.New()
	if model.ValidateProject(v, project); !v.Valid() {
		return nil, controller.ErrFailedValidation
	}
	err := c.repo.Create(ctx, project)
	if err != nil {
		return nil, err
	}
	return project, nil
}
