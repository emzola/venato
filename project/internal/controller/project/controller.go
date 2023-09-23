package project

import (
	"context"
	"errors"
	"time"

	"github.com/emzola/venato/project/internal/controller"
	"github.com/emzola/venato/project/internal/repository"
	"github.com/emzola/venato/project/pkg/model"
	"github.com/emzola/venato/project/pkg/validator"
)

type projectRepository interface {
	Create(ctx context.Context, project *model.Project) error
	Get(ctx context.Context, id int64) (*model.Project, error)
	// GetAll(ctx context.Context) ([]*model.Project, model.Metadata, error)
	Update(ctx context.Context, project *model.Project) error
	Delete(ctx context.Context, id int64) error
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
func (c *Controller) Create(ctx context.Context, name, description string, startDate, targetEndDate time.Time, createdBy, modifiedBy int64) (*model.Project, error) {
	project := &model.Project{
		Name:          name,
		Description:   description,
		StartDate:     startDate,
		TargetEndDate: targetEndDate,
		CreatedBy:     createdBy,
		ModifiedBy:    modifiedBy,
	}
	v := validator.New()
	if model.ValidateProject(v, project); !v.Valid() {
		controller.ErrFailedValidation = controller.FailedValidation(v.Errors)
		return nil, controller.ErrFailedValidation
	}
	err := c.repo.Create(ctx, project)
	if err != nil {
		return nil, err
	}
	return project, nil
}

// Get retrieves a project by id.
func (c *Controller) Get(ctx context.Context, id int64) (*model.Project, error) {
	project, err := c.repo.Get(ctx, id)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrNotFound):
			return nil, controller.ErrNotFound
		default:
			return nil, err
		}
	}
	return project, nil
}

// GetAll retrieves a paginated list of all projects.
// func (c *Controller) GetAll(ctx context.Context) ([]*model.Project, model.Metadata, error) {

// }

// Update partially updates a project record.
func (c *Controller) Update(ctx context.Context, id int64, name, description *string, startDate, targetEndDate, actualEndDate *time.Time, modifiedBy int64) (*model.Project, error) {
	project, err := c.repo.Get(ctx, id)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrNotFound):
			return nil, controller.ErrNotFound
		default:
			return nil, err
		}
	}
	// Partially update the project with new data based on whether new data is supplied by the client.
	if name != nil {
		project.Name = *name
	}
	if description != nil {
		project.Description = *description
	}
	if startDate != nil {
		project.StartDate = *startDate
	}
	if targetEndDate != nil {
		project.TargetEndDate = *targetEndDate
	}
	if actualEndDate != nil {
		project.ActualEndDate = *actualEndDate
	}
	project.ModifiedBy = modifiedBy
	v := validator.New()
	if model.ValidateProject(v, project); !v.Valid() {
		controller.ErrFailedValidation = controller.FailedValidation(v.Errors)
		return nil, controller.ErrFailedValidation
	}
	err = c.repo.Update(ctx, project)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrEditConflict):
			return nil, controller.ErrEditConflict
		default:
			return nil, err
		}
	}
	return project, nil
}

// Delete removes a project by its id.
func (c *Controller) Delete(ctx context.Context, id int64) error {
	err := c.repo.Delete(ctx, id)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrNotFound):
			return controller.ErrNotFound
		default:
			return err
		}
	}
	return nil
}
