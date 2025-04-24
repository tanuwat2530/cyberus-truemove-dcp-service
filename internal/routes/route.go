package routes

import (
	"cyberus/truemove-dcp-service/internal/controllers"
	"net/http"
)

// SetupRoutes registers all application routes
func SetupRoutes() {
	// Register routes using http.HandleFunc
	http.HandleFunc("/tmvh/wap-redirect", controllers.WapRedirect)
	http.HandleFunc("/tmvh/subscription-callback", controllers.SubscriptionCallback)
	http.HandleFunc("/tmvh/transaction-callback", controllers.TransactionCallback)
	http.HandleFunc("/tmvh/", HomeHandler)
}

// HomeHandler for root endpoint
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Welcome to backend API power by GoLang ^_^"))
}
