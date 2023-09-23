package http

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (h *Handler) Routes() http.Handler {
	router := httprouter.New()
	router.NotFound = http.HandlerFunc(h.notFoundResponse)
	router.MethodNotAllowed = http.HandlerFunc(h.methodNotAllowedResponse)
	router.HandlerFunc(http.MethodGet, "/projects", h.getAllProjects)
	router.HandlerFunc(http.MethodPost, "/projects", h.createProject)
	router.HandlerFunc(http.MethodGet, "/projects/:id", h.getProject)
	router.HandlerFunc(http.MethodPatch, "/projects/:id", h.updateProject)
	router.HandlerFunc(http.MethodDelete, "/projects/:id", h.deleteProject)
	return router
}
