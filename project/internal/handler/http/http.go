package http

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/emzola/venato/project/internal/controller"
	"github.com/emzola/venato/project/pkg/model"
)

type projectController interface {
	Create(ctx context.Context, name string, description string, startDate time.Time, targetEndDate time.Time, createdBy int64) (*model.Project, error)
	// Get(ctx context.Context, id int64) (*model.Project, error)
	// Update(ctx context.Context, project *model.Project) error
	// Delete(ctx context.Context, id int64) error
}

// Handler defines a project HTTP handler.
type Handler struct {
	ctrl projectController
}

// New creates a new project HTTP handler.
func New(ctrl projectController) *Handler {
	return &Handler{ctrl}
}

// Create handles POST /project requests.
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
	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()
	project, err := h.ctrl.Create(ctx, requestBody.Name, requestBody.Description, requestBody.StartDate, requestBody.TargetEndDate, 1)
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

func (h *Handler) getProject(w http.ResponseWriter, r *http.Request) {

}

func (h *Handler) updateProject(w http.ResponseWriter, r *http.Request) {

}

func (h *Handler) deleteProject(w http.ResponseWriter, r *http.Request) {

}

func (h *Handler) listProject(w http.ResponseWriter, r *http.Request) {

}
