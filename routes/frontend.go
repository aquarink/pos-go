package routes

import (
	"pos/controllers"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
)

func RegisterFrontendRoutes(router *mux.Router, store *sessions.CookieStore) {
	router.HandleFunc("/", controllers.HomeController)
}
