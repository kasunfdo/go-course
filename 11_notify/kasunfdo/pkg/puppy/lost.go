package puppy

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
)

// LostAPIHandler implements REST API handler of lost service.
type LostAPIHandler struct{}

// NewLostAPIHandler creates a LostAPIHandler.
func NewLostAPIHandler() *LostAPIHandler {
	return &LostAPIHandler{}
}

// HandlePostLostPuppy handles http request to lost puppy service.
func (a *LostAPIHandler) HandlePostLostPuppy(w http.ResponseWriter, r *http.Request) {
	var payload Payload
	if err := render.DecodeJSON(r.Body, &payload); err != nil {
		sendErrorResponse(w, http.StatusBadRequest, ErrorEf(ErrInvalid, err, ""))
		return
	}

	var responseStatus int

	switch payload.ID % 2 {
	case 0:
		responseStatus = http.StatusCreated
	case 1:
		responseStatus = http.StatusInternalServerError
	}

	time.Sleep(2 * time.Second)

	w.WriteHeader(responseStatus)
	render.JSON(w, r, fmt.Sprintf("Status: %d", responseStatus))
}

// WireLostServiceRoutes route requests to corresponding REST API handler method.
func (a *LostAPIHandler) WireLostServiceRoutes(r chi.Router) {
	r.Post("/api/lostpuppy/", a.HandlePostLostPuppy)
}
