package grpc

import (
	"context"
	"errors"

	"github.com/emzola/venato/gen"
	"github.com/emzola/venato/project/internal/controller"
	"github.com/emzola/venato/project/internal/controller/project"
	"github.com/emzola/venato/project/pkg/model"
)

// Handler defines a Project gRPC handler.
type Handler struct {
	gen.UnimplementedProjectServiceServer
	ctrl *project.Controller
}

// New creates a new project gRPC handler.
func New(ctrl *project.Controller) *Handler {
	return &Handler{ctrl: ctrl}
}

// CreateProject creates a new project record.
func (h *Handler) CreateProject(ctx context.Context, req *gen.CreateProjectRequest) (*gen.CreateProjectResponse, error) {
	if req == nil {
		return nil, nilRequestError
	}
	var userID int64 = 1
	project, err := h.ctrl.Create(ctx, req.Name, req.Description, req.StartDate.AsTime(), req.TargetEndDate.AsTime(), userID, userID)
	if err != nil {
		switch {
		case errors.Is(err, controller.ErrFailedValidation):
			return nil, h.failedValidationError(err)
		default:
			return nil, internalServerError
		}
	}
	return &gen.CreateProjectResponse{Project: model.ProjectToProto(project)}, nil
}

// GetProject returns the project for a given record.
func (h *Handler) GetProject(ctx context.Context, req *gen.GetProjectRequest) (*gen.GetProjectResponse, error) {
	if req == nil {
		return nil, nilRequestError
	}
	id := req.ProjectId
	if id < 1 {
		return nil, notFoundError
	}
	project, err := h.ctrl.Get(ctx, id)
	if err != nil {
		switch {
		case errors.Is(err, context.Canceled):
			return nil, nil
		case errors.Is(err, controller.ErrNotFound):
			return nil, notFoundError
		default:
			return nil, internalServerError
		}
	}
	return &gen.GetProjectResponse{Project: model.ProjectToProto(project)}, nil
}

// UpdateProject updates the project for a given record.
func (h *Handler) UpdateProject(ctx context.Context, req *gen.UpdateProjectRequest) (*gen.UpdateProjectResponse, error) {
	if req == nil {
		return nil, nilRequestError
	}
	id := req.ProjectId
	if id < 1 {
		return nil, notFoundError
	}
	var userID int64 = 1
	startDate := req.StartDate.AsTime()
	targetEndDate := req.TargetEndDate.AsTime()
	actualEndDate := req.ActualEndDate.AsTime()
	project, err := h.ctrl.Update(ctx, id, &req.Name, &req.Description, &startDate, &targetEndDate, &actualEndDate, userID)
	if err != nil {
		switch {
		case errors.Is(err, context.Canceled):
			return nil, nil
		case errors.Is(err, controller.ErrFailedValidation):
			return nil, h.failedValidationError(err)
		case errors.Is(err, controller.ErrEditConflict):
			return nil, editConflictError
		default:
			return nil, internalServerError
		}
	}
	return &gen.UpdateProjectResponse{Project: model.ProjectToProto(project)}, nil
}

// DeleteProject deletes the project for a given record.
func (h *Handler) DeleteProject(ctx context.Context, req *gen.DeleteProjectRequest) (*gen.DeleteProjectResponse, error) {
	if req == nil {
		return nil, nilRequestError
	}
	id := req.ProjectId
	if id < 1 {
		return nil, notFoundError
	}
	err := h.ctrl.Delete(ctx, id)
	if err != nil {
		switch {
		case errors.Is(err, context.Canceled):
			return nil, nil
		case errors.Is(err, controller.ErrNotFound):
			return nil, notFoundError
		default:
			return nil, internalServerError
		}
	}
	return &gen.DeleteProjectResponse{Message: "project successfully deleted"}, nil
}
