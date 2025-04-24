package main

import (
	"cyberus/truemove-dcp-service/internal/routes"
	"fmt"
	"net/http"
)

func main() {
	// Setup all routes
	routes.SetupRoutes()

	// Start the server on port 8080
	fmt.Println("Starting cyberus-truemove-dcp-service server on port 8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println("Error starting server:", err)
	}

}
