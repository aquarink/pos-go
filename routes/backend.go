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

	router.Handle("/app/signout", middleware.CheckSignin(store)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		controllers.SignoutController(w, r, client, store)
	}))).Methods("GET")

	router.Handle("/app/dashboard", middleware.CheckSignin(store)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		controllers.DashboardController(w, r, client, store)
	}))).Methods("GET")

	// CATEGORY

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

	router.Handle("/app/category/remove/{id}", middleware.CheckSignin(store)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		controllers.CategoryDelete(w, r, client, store)
	}))).Methods("GET")

	// PRODUCT

	router.Handle("/app/product", middleware.CheckSignin(store)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		controllers.ProductList(w, r, client, store)
	}))).Methods("GET")

	router.Handle("/app/product/", middleware.CheckSignin(store)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		controllers.ProductList(w, r, client, store)
	}))).Methods("GET")

	router.Handle("/app/product/list", middleware.CheckSignin(store)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		controllers.ProductList(w, r, client, store)
	}))).Methods("GET")

	router.Handle("/app/product/add", middleware.CheckSignin(store)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		controllers.ProductAdd(w, r, client, store)
	}))).Methods("GET", "POST")

	router.Handle("/app/product/edit/{id}", middleware.CheckSignin(store)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		controllers.ProductEdit(w, r, client, store)
	}))).Methods("GET")

	router.Handle("/app/product/update", middleware.CheckSignin(store)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		controllers.ProductUpdate(w, r, client, store)
	}))).Methods("POST")

	router.Handle("/app/product/delete/{id}", middleware.CheckSignin(store)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		controllers.ProductDelete(w, r, client, store)
	}))).Methods("GET")

	// PACKAGE

	router.Handle("/app/package", middleware.CheckSignin(store)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		controllers.PackageList(w, r, client, store)
	}))).Methods("GET")

	router.Handle("/app/package/list", middleware.CheckSignin(store)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		controllers.PackageList(w, r, client, store)
	}))).Methods("GET")

	// add form dan submit
	router.Handle("/app/package/add", middleware.CheckSignin(store)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		controllers.PackageAdd(w, r, client, store)
	}))).Methods("GET", "POST")

	// edit form
	router.Handle("/app/package/edit/{id}", middleware.CheckSignin(store)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		controllers.PackageEdit(w, r, client, store)
	}))).Methods("GET")

	// edit form submit
	router.Handle("/app/package/update", middleware.CheckSignin(store)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		controllers.PackageUpdate(w, r, client, store)
	}))).Methods("POST")

	router.Handle("/app/package/remove/{id}", middleware.CheckSignin(store)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		controllers.PackageDelete(w, r, client, store)
	}))).Methods("GET")

	// STORE

	router.Handle("/app/store", middleware.CheckSignin(store)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		controllers.StoreEdit(w, r, client, store)
	}))).Methods("GET")

	router.Handle("/app/store/update", middleware.CheckSignin(store)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		controllers.StoreUpdate(w, r, client, store)
	}))).Methods("POST")

	router.Handle("/app/merchant", middleware.CheckSignin(store)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		controllers.MerchantList(w, r, client, store)
	}))).Methods("GET")

	// ORDER

	router.Handle("/app/order", middleware.CheckSignin(store)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		controllers.Order(w, r, client, store)
	}))).Methods("GET")

	router.Handle("/app/checkout", middleware.CheckSignin(store)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		controllers.Checkout(w, r, client, store)
	}))).Methods("POST")

	// TRANSACTION

	router.Handle("/app/transaction", middleware.CheckSignin(store)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		controllers.TransactionList(w, r, client, store)
	}))).Methods("GET")

	// CASHIER

	router.Handle("/app/cashier", middleware.CheckSignin(store)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		controllers.CashierList(w, r, client, store)
	}))).Methods("GET")

	router.Handle("/app/cashier/list", middleware.CheckSignin(store)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		controllers.CashierList(w, r, client, store)
	}))).Methods("GET")

	// add form dan submit
	router.Handle("/app/cashier/add", middleware.CheckSignin(store)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		controllers.CashierAdd(w, r, client, store)
	}))).Methods("GET", "POST")

	router.Handle("/app/cashier/remove/{id}", middleware.CheckSignin(store)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		controllers.CashierDelete(w, r, client, store)
	}))).Methods("GET")

	// BILLING

	router.Handle("/app/billing", middleware.CheckSignin(store)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		controllers.Billing(w, r, client, store)
	}))).Methods("GET")

	router.Handle("/app/billing/{id}", middleware.CheckSignin(store)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		controllers.Billing(w, r, client, store)
	}))).Methods("GET")

	// TABLE

	router.Handle("/app/table", middleware.CheckSignin(store)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		controllers.TableList(w, r, client, store)
	}))).Methods("GET")

	router.Handle("/app/table/generate", middleware.CheckSignin(store)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		controllers.TableNoGenerate(w, r, client, store)
	}))).Methods("GET")

	// USER

	router.Handle("/app/password", middleware.CheckSignin(store)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		controllers.Password(w, r, client, store)
	}))).Methods("GET", "POST")

	// USER VERIFY

	router.Handle("/app/verify/{id}", middleware.CheckSession(store)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		controllers.SignupVerifyController(w, r, client, store)
	}))).Methods("GET")
}
