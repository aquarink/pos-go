package routes

import (
	"net/http"
	"pos/controllers"
	"pos/middleware"
	"pos/services"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
)

func RegisterBackendRoutes(router *mux.Router, client *services.AppwriteClient, store *sessions.CookieStore) {
	router.Handle("/app", middleware.CheckSession(store)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		controllers.SigninController(w, r, client, store)
	}))).Methods("GET", "POST")

	router.Handle("/app/signin", middleware.CheckSession(store)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		controllers.SigninController(w, r, client, store)
	}))).Methods("GET", "POST")

	router.Handle("/app/signup", middleware.CheckSession(store)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		controllers.SignupController(w, r, client, store)
	}))).Methods("GET", "POST")

	router.Handle("/app/dashboard", middleware.CheckSignin(store)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		controllers.DashboardController(w, r, client, store)
	}))).Methods("GET")

	router.Handle("/app/category", middleware.CheckSignin(store)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		controllers.CategoryList(w, r, client, store)
	}))).Methods("GET")

	router.Handle("/app/category/list", middleware.CheckSignin(store)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		controllers.CategoryList(w, r, client, store)
	}))).Methods("GET")

	// add form dan submit
	router.Handle("/app/category/add", middleware.CheckSignin(store)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		controllers.CategoryAdd(w, r, client, store)
	}))).Methods("GET", "POST")

	// edit form
	router.Handle("/app/category/edit/{id}", middleware.CheckSignin(store)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		controllers.CategoryEdit(w, r, client, store)
	}))).Methods("GET")

	// edit form submit
	router.Handle("/app/category/edit", middleware.CheckSignin(store)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		controllers.CategoryUpdate(w, r, client, store)
	}))).Methods("POST")

	router.Handle("/app/signout", middleware.CheckSignin(store)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		controllers.SignoutController(w, r, client, store)
	}))).Methods("GET")
}
