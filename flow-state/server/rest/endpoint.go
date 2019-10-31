package rest

import (
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/project-flogo/core/support/log"
	"github.com/project-flogo/services/flow-state/server/rest/cors"
)

type PreflightHandler struct {
	c cors.Cors
}

// Handles the cors preflight request
func (h *PreflightHandler) handleCorsPreflight(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	h.c.HandlePreflight(w, r)
}

type EndpointManager struct {
	router *httprouter.Router
	pfh    *PreflightHandler
	paths   map[string]void
}

type void struct{}
var empty void

func NewEndpointManager(router *httprouter.Router, logger log.Logger, usePreflight bool) *EndpointManager {
	em := &EndpointManager{router: router}
	if usePreflight {
		em.paths = make(map[string]void)
		em.pfh = &PreflightHandler{c: cors.New(CorsPrefix, logger)}
	}

	em.GET("/status", status)
	return em
}

func (em *EndpointManager) addPreflight(path string) {
	if em.pfh != nil {
		if _, exists := em.paths[path]; !exists {
			em.paths[path] = empty
			em.router.OPTIONS(path, em.pfh.handleCorsPreflight) // for CORS
		}
	}
}

func (em *EndpointManager) GET(path string, handle httprouter.Handle) {
	em.router.GET(path, handle)
	em.addPreflight(path)
}

func (em *EndpointManager) POST(path string, handle httprouter.Handle) {
	em.router.POST(path, handle)
	em.addPreflight(path)
}

func (em *EndpointManager) DELETE(path string, handle httprouter.Handle) {
	em.router.DELETE(path, handle)
	em.addPreflight(path)
}

func status(response http.ResponseWriter, request *http.Request, _ httprouter.Params) {
	response.Header().Set("Content-Type", "application/json")
	response.WriteHeader(http.StatusOK)
	fmt.Fprintf(response, "%s", "{\"status\":\"ok\"}")
}
