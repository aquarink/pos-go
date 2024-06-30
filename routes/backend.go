package routes

import (
	"net/http"
	"pos/controllers"
)

// registerRouteWithPrefix menambahkan prefiks "/app" pada setiap rute
func registerRouteWithPrefix(path string, handler http.HandlerFunc) {
	http.HandleFunc("/app"+path, handler)
}

func RegisterBackendRoutes() {
	registerRouteWithPrefix("/users", controllers.UserController)
}
