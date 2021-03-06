package accounts

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/RichardKnop/example-api/logger"
	"github.com/RichardKnop/example-api/util/response"
	"github.com/gorilla/mux"
)

// Handles requests to complete an invitation of a user by setting password
// POST /v1/invitations/{reference:[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}}
func (s *Service) confirmInvitationHandler(w http.ResponseWriter, r *http.Request) {
	// Get the authenticated client from the request context
	_, err := GetAuthenticatedClient(r)
	if err != nil {
		response.UnauthorizedError(w, err.Error())
		return
	}

	// Request body cannot be nil
	if r.Body == nil {
		response.Error(w, "Request body cannot be nil", http.StatusBadRequest)
		return
	}

	// Read the request body
	payload, err := ioutil.ReadAll(r.Body)
	if err != nil {
		response.Error(w, "Error reading request body", http.StatusBadRequest)
		return
	}

	// Unmarshal the request body into the request prototype
	confirmInvitationRequest := new(ConfirmInvitationRequest)
	if err = json.Unmarshal(payload, confirmInvitationRequest); err != nil {
		logger.ERROR.Printf("Failed to unmarshal confirm invitation request: %s", payload)
		response.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Get the reference from request URI
	vars := mux.Vars(r)
	reference := vars["reference"]

	// Fetch the invitation we want to work with (by reference from email link)
	invitation, err := s.FindInvitationByReference(reference)
	if err != nil {
		response.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// Confirm the invitation
	if err = s.ConfirmInvitation(
		invitation,
		confirmInvitationRequest.Password,
	); err != nil {
		response.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Create invitation response
	invitationResponse, err := NewInvitationResponse(invitation)
	if err != nil {
		response.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Write the response
	response.WriteJSON(w, invitationResponse, 200)
}
