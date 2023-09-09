package http

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (h *Handler) Routes() http.Handler {
	router := httprouter.New()
	router.NotFound = http.HandlerFunc(h.notFoundResponse)
	router.MethodNotAllowed = http.HandlerFunc(h.methodNotAllowedResponse)
	router.HandlerFunc(http.MethodGet, "/project", h.listProject)
	router.HandlerFunc(http.MethodPost, "/project", h.createProject)
	router.HandlerFunc(http.MethodGet, "/project/:id", h.getProject)
	router.HandlerFunc(http.MethodPatch, "/project/:id", h.updateProject)
	router.HandlerFunc(http.MethodDelete, "/project/:id", h.deleteProject)
	return router
}
