package main

import (
	"html/template"
	"log"
	"net/http"
	"os"
	"pos/middleware"
	"pos/models"
	"pos/routes"
	"pos/services"

	"github.com/gorilla/mux"
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
	router := mux.NewRouter()

	// middleware
	router.Use(middleware.SessionMiddleware)

	routes.RegisterFrontendRoutes(router)
	routes.RegisterBackendRoutes(router, client)

	fs := http.FileServer(http.Dir("assets"))
	router.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", fs))

	router.NotFoundHandler = http.HandlerFunc(Handle404)

	log.Println("Server started at :8080")
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

func Handle404(w http.ResponseWriter, r *http.Request) {
	data := models.PublicData{
		Title: "Page Not Found",
		Data:  map[string]interface{}{},
		Error: r.URL.Query().Get("error"),
		Msg:   r.URL.Query().Get("msg"),
	}

	tmpl, err := template.ParseFiles("views/pages/auth/404.html")
	if err != nil {
		log.Println("Error parsing 404 template:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Println("Rendering 404 template")
	err = tmpl.Execute(w, data)
	if err != nil {
		log.Println("Error executing 404 template:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
