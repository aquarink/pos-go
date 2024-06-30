package main

import (
	"log"
	"net/http"
	"pos/routes"
	"pos/services"
)

func main() {
	// Inisialisasi Appwrite client
	services.InitAppwriteClient("http://localhost/v1", "YOUR_PROJECT_ID", "YOUR_API_KEY")

	routes.RegisterFrontendRoutes()
	routes.RegisterBackendRoutes()

	log.Println("Server started at :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
