package main

import (
	"log"
	"net/http"
	"os"
	"pos/routes"
	"pos/services"

	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	appwriteEndpoint := os.Getenv("APPWRITE_ENDPOINT")
	appwriteProjectID := os.Getenv("APPWRITE_PROJECT_ID")
	appwriteAPIKey := os.Getenv("APPWRITE_API_KEY")
	appwriteDatabaseID := os.Getenv("POS_DB")

	if appwriteEndpoint == "" || appwriteProjectID == "" || appwriteAPIKey == "" || appwriteDatabaseID == "" {
		log.Fatalf("Some required environment variables are missing")
	}

	log.Printf("Appwrite endpoint: %s", appwriteEndpoint)
	log.Printf("Appwrite project ID: %s", appwriteProjectID)

	// Initialize Appwrite client
	client := services.NewAppwriteClient(appwriteEndpoint, appwriteProjectID, appwriteAPIKey, appwriteDatabaseID)

	routes.RegisterFrontendRoutes()
	routes.RegisterBackendRoutes(client)

	fs := http.FileServer(http.Dir("assets"))
	http.Handle("/assets/", http.StripPrefix("/assets/", fs))

	log.Println("Server started at :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
