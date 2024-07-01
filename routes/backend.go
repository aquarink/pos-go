package routes

import (
	"net/http"
	"pos/controllers"
	"pos/services"
)

func registerRouteWithPrefix(path string, handler http.HandlerFunc) {
	http.HandleFunc("/app"+path, handler)
}

func RegisterBackendRoutes(client *services.AppwriteClient) {
	registerRouteWithPrefix("/signup", func(w http.ResponseWriter, r *http.Request) {
		controllers.SignupController(w, r, client)
	})
	registerRouteWithPrefix("/signin", func(w http.ResponseWriter, r *http.Request) {
		controllers.SigninController(w, r, client)
	})
	registerRouteWithPrefix("/users", func(w http.ResponseWriter, r *http.Request) {
		controllers.UserController(w, r, client)
	})
}
