package controllers

import (
	"html/template"
	"net/http"
	"os"
	"pos/services"
)

func UserController(w http.ResponseWriter, r *http.Request, client *services.AppwriteClient) {
	if r.Method == http.MethodGet {
		users, err := client.GetAllUsers(os.Getenv("USERS"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		tmpl, err := template.ParseFiles(
			"views/templates/backend.html",
			"views/pages/users/list.html",
		)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		tmpl.ExecuteTemplate(w, "layout", users)
	}
}

func HomeController(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("views/templates/frontend.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tmpl.Execute(w, nil)
}

// Contoh controller untuk dashboard
func DashboardController(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("views/templates/backend.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tmpl.Execute(w, nil)
}
