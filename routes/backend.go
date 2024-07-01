package routes

import (
	"net/http"
	"pos/controllers"
	"pos/services"

	"github.com/gorilla/mux"
)

func RegisterBackendRoutes(router *mux.Router, client *services.AppwriteClient) {
	router.HandleFunc("/app/signup", func(w http.ResponseWriter, r *http.Request) {
		controllers.SignupController(w, r, client)
	}).Methods("GET", "POST")
	router.HandleFunc("/app/signin", func(w http.ResponseWriter, r *http.Request) {
		controllers.SigninController(w, r, client)
	}).Methods("GET", "POST")
	router.HandleFunc("/app/dashboard", func(w http.ResponseWriter, r *http.Request) {
		controllers.DashboardController(w, r, client)
	}).Methods("GET")
}
