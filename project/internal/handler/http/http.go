package http

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/emzola/venato/project/internal/controller"
	"github.com/emzola/venato/project/pkg/model"
	"github.com/emzola/venato/project/pkg/validator"
)

type projectController interface {
	Create(ctx context.Context, name, description string, startDate, targetEndDate time.Time, createdBy, modifiedBy int64) (*model.Project, error)
	Get(ctx context.Context, id int64) (*model.Project, error)
	// GetAll(ctx context.Context) ([]*model.Project, model.Metadata, error)
	Update(ctx context.Context, id int64, name, description *string, startDate, targetEndDate, actualEndDate *time.Time, modifiedBy int64) (*model.Project, error)
	Delete(ctx context.Context, id int64) error
}

// Handler defines a project HTTP handler.
type Handler struct {
	ctrl projectController
}

// New creates a new project HTTP handler.
func New(ctrl projectController) *Handler {
	return &Handler{ctrl}
}

// createProject handles POST /projects requests for creating a new project.
func (h *Handler) createProject(w http.ResponseWriter, r *http.Request) {
	var requestBody struct {
		Name          string    `json:"name"`
		Description   string    `json:"description"`
		StartDate     time.Time `json:"start_date"`
		TargetEndDate time.Time `json:"target_end_date"`
	}
	err := h.decodeJSON(w, r, &requestBody)
	if err != nil {
		h.badRequestResponse(w, r, err)
		return
	}
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()
	project, err := h.ctrl.Create(ctx, requestBody.Name, requestBody.Description, requestBody.StartDate, requestBody.TargetEndDate, 1, 1)
	if err != nil {
		switch {
		case errors.Is(err, controller.ErrFailedValidation):
			h.failedValidationResponse(w, r, err)
		default:
			h.serverErrorResponse(w, r, err)
		}
		return
	}
	header := make(http.Header)
	header.Set("Location", fmt.Sprintf("/project/%d", project.ID))
	err = h.encodeJSON(w, http.StatusCreated, envelop{"project": project}, header)
	if err != nil {
		h.serverErrorResponse(w, r, err)
	}
}

// getProject handles GET /projects requests for retrieving a project.
func (h *Handler) getProject(w http.ResponseWriter, r *http.Request) {
	id, err := h.readIDParam(r, "id")
	if err != nil {
		h.badRequestResponse(w, r, err)
		return
	}
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()
	project, err := h.ctrl.Get(ctx, id)
	if err != nil {
		switch {
		case errors.Is(err, context.Canceled):
			return
		case errors.Is(err, controller.ErrNotFound):
			h.notFoundResponse(w, r)
		default:
			h.serverErrorResponse(w, r, err)
		}
		return
	}
	err = h.encodeJSON(w, http.StatusOK, envelop{"project": project}, nil)
	if err != nil {
		h.serverErrorResponse(w, r, err)
	}
}

// getAllProjects handles GET /projects requests for retrieving a paginated list of all projects.
func (h *Handler) getAllProjects(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name string
		model.Filters
	}
	v := validator.New()
	qs := r.URL.Query()
	input.Name = h.readString(qs, "name", "")
	input.Filters.Page = h.readInt(qs, "page", 1, v)
	input.Filters.PageSize = h.readInt(qs, "page_size", 20, v)
	input.Filters.Sort = h.readString(qs, "sort", "id")
	input.Filters.SortSafelist = []string{"id", "name", "created_on", "modified_on", "-id", "-name", "-created_on", "- modified_on"}
	// if !v.Valid() {
	// 	h.failedValidationResponse(w, r, v.Errors)
	// }
}

// updateProject handles PATCH /projects requests for updating a project.
func (h *Handler) updateProject(w http.ResponseWriter, r *http.Request) {
	id, err := h.readIDParam(r, "id")
	if err != nil {
		h.notFoundResponse(w, r)
		return
	}
	var requestBody struct {
		Name          *string    `json:"name"`
		Description   *string    `json:"description"`
		StartDate     *time.Time `json:"start_date"`
		TargetEndDate *time.Time `json:"target_end_date"`
		ActualEndDate *time.Time `json:"actual_end_date,omitempty"`
	}
	err = h.decodeJSON(w, r, &requestBody)
	if err != nil {
		h.badRequestResponse(w, r, err)
		return
	}
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()
	project, err := h.ctrl.Update(ctx, id, requestBody.Name, requestBody.Description, requestBody.StartDate, requestBody.TargetEndDate, requestBody.ActualEndDate, 1)
	if err != nil {
		switch {
		case errors.Is(err, context.Canceled):
			return
		case errors.Is(err, controller.ErrFailedValidation):
			h.failedValidationResponse(w, r, err)
		case errors.Is(err, controller.ErrEditConflict):
			h.editConflictResponse(w, r)
		default:
			h.serverErrorResponse(w, r, err)
		}
		return
	}
	err = h.encodeJSON(w, http.StatusOK, envelop{"project": project}, nil)
	if err != nil {
		h.serverErrorResponse(w, r, err)
	}
}

// deleteProject handles DELETE /projects requests for deleting a project.
func (h *Handler) deleteProject(w http.ResponseWriter, r *http.Request) {
	id, err := h.readIDParam(r, "id")
	if err != nil {
		h.notFoundResponse(w, r)
		return
	}
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()
	err = h.ctrl.Delete(ctx, id)
	if err != nil {
		switch {
		case errors.Is(err, context.Canceled):
			return
		case errors.Is(err, controller.ErrNotFound):
			h.notFoundResponse(w, r)
		default:
			h.serverErrorResponse(w, r, err)
		}
		return
	}
	err = h.encodeJSON(w, http.StatusOK, envelop{"message": "project successfully deleted"}, nil)
	if err != nil {
		h.serverErrorResponse(w, r, err)
	}
}
