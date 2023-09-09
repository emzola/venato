package http

import (
	"fmt"
	"net/http"

	"go.uber.org/zap"
)

func (h *Handler) logError(r *http.Request, err error) {
	logger, _ := zap.NewProduction()
	defer logger.Sync()
	logger.Info(fmt.Sprintf("%s", err),
		zap.String("request_method", r.Method),
		zap.String("request_url", r.URL.String()),
	)
}

func (h *Handler) errorResponse(w http.ResponseWriter, r *http.Request, status int, message interface{}) {
	env := envelop{"error": message}
	err := h.encodeJSON(w, status, env, nil)
	if err != nil {
		h.logError(r, err)
		w.WriteHeader(500)
	}
}

func (h *Handler) serverErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	h.logError(r, err)
	message := "the server encountered a problem and could not process your request"
	h.errorResponse(w, r, http.StatusInternalServerError, message)
}

func (h *Handler) notFoundResponse(w http.ResponseWriter, r *http.Request) {
	message := "the requested resource could not be found"
	h.errorResponse(w, r, http.StatusNotFound, message)
}

func (h *Handler) methodNotAllowedResponse(w http.ResponseWriter, r *http.Request) {
	message := fmt.Sprintf("the %s method is not supported for this resource", r.Method)
	h.errorResponse(w, r, http.StatusMethodNotAllowed, message)
}

func (h *Handler) badRequestResponse(w http.ResponseWriter, r *http.Request, err error) {
	h.errorResponse(w, r, http.StatusBadRequest, err.Error())
}

func (h *Handler) editConflictResponse(w http.ResponseWriter, r *http.Request) {
	message := "unable to update the record due to an edit conflict, please try again"
	h.errorResponse(w, r, http.StatusConflict, message)
}

func (h *Handler) failedValidationResponse(w http.ResponseWriter, r *http.Request, err error) {
	h.errorResponse(w, r, http.StatusUnprocessableEntity, err.Error())
}
