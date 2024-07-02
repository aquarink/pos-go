package routes

import (
	"net/http"
	"pos/controllers"
	"pos/middleware"
	"pos/services"

	"github.com/gorilla/mux"
)

func RegisterBackendRoutes(router *mux.Router, client *services.AppwriteClient) {
	router.Handle("/app", middleware.CheckSession(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		controllers.SigninController(w, r, client)
	}))).Methods("GET", "POST")

	router.Handle("/app/signin", middleware.CheckSession(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		controllers.SigninController(w, r, client)
	}))).Methods("GET", "POST")

	router.Handle("/app/signup", middleware.CheckSession(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		controllers.SignupController(w, r, client)
	}))).Methods("GET", "POST")

	router.Handle("/app/dashboard", middleware.CheckSignin(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		controllers.DashboardController(w, r, client)
	}))).Methods("GET")

	//
	router.Handle("/app/category", middleware.CheckSignin(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		controllers.CategoryList(w, r, client)
	}))).Methods("GET")

	router.Handle("/app/category/list", middleware.CheckSignin(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		controllers.CategoryList(w, r, client)
	}))).Methods("GET")

	// add form dan submit
	router.Handle("/app/category/add", middleware.CheckSignin(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		controllers.CategoryAdd(w, r, client)
	}))).Methods("GET", "POST")

	// edit form dan submit
	router.Handle("/app/category/edit", middleware.CheckSignin(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		controllers.CategoryEdit(w, r, client)
	}))).Methods("GET", "POST")
	//

	router.Handle("/app/signout", middleware.CheckSignin(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		controllers.SignoutController(w, r, client)
	}))).Methods("GET")
}
