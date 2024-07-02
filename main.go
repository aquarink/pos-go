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

	"github.com/gorilla/csrf"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/joho/godotenv"
)

var sessStore *sessions.CookieStore

func main() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Some required environment variables are missing in env")
	}

	appwriteEndpoint := os.Getenv("APPWRITE_ENDPOINT")
	appwriteProjectID := os.Getenv("APPWRITE_PROJECT_ID")
	appwriteAPIKey := os.Getenv("APPWRITE_API_KEY")
	appwriteDatabaseID := os.Getenv("POS_DB")

	if appwriteEndpoint == "" || appwriteProjectID == "" || appwriteAPIKey == "" || appwriteDatabaseID == "" {
		log.Fatalf("Some required environment variables are missing in appwrite")
	}

	csrfAuthKey := os.Getenv("CSRF_AUTH_KEY")
	cookiesKey := os.Getenv("COOKIES_KEY")

	if csrfAuthKey == "" || cookiesKey == "" {
		log.Fatalf("Some required environment variables are missing in csrf")
	}

	if len(cookiesKey) < 32 {
		log.Fatal("COOKIES_KEY must be at least 32 bytes long")
	}

	sessStore = sessions.NewCookieStore([]byte(cookiesKey))
	log.Printf("Cookie store initialized with key: %s", cookiesKey)

	// Initialize Appwrite client
	client := services.NewAppwriteClient(appwriteEndpoint, appwriteProjectID, appwriteAPIKey, appwriteDatabaseID)
	router := mux.NewRouter()

	csrfMiddleware := csrf.Protect([]byte(csrfAuthKey), csrf.Secure(false))

	// middleware
	router.Use(middleware.SessionMiddleware(sessStore))

	routes.RegisterFrontendRoutes(router, sessStore)
	routes.RegisterBackendRoutes(router, client, sessStore)

	fs := http.FileServer(http.Dir("assets"))
	router.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", fs))

	router.NotFoundHandler = http.HandlerFunc(Handle404)

	if err := http.ListenAndServe(":8080", csrfMiddleware(router)); err != nil {
		log.Fatalf("Some required environment variables are missing in ListenAndServe")
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
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
