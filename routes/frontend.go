package routes

import (
	"net/http"
	"pos/controllers"
)

func RegisterFrontendRoutes() {
	http.HandleFunc("/", controllers.HomeController)
}
