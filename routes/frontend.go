package routes

import (
	"pos/controllers"

	"github.com/gorilla/mux"
)

func RegisterFrontendRoutes(router *mux.Router) {
	router.HandleFunc("/", controllers.HomeController)
}
