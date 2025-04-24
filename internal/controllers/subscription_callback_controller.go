package controllers

import (
	"CyberusGolangShareLibrary/utilities"
	services "cyberus/truemove-dcp-service/internal/services"

	"net/http"
)

func SubscriptionCallback(w http.ResponseWriter, r *http.Request) {
	// Check if the method is POST
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	response := services.SubscriptionCallbackProcessRequest(r)

	utilities.ResponseWithJSON(w, http.StatusOK, response)
}
