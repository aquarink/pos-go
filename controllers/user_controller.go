package controllers

import (
	"html/template"
	"net/http"
	"os"
	"pos/models"
	"pos/services"
	"pos/utils"

	"github.com/gorilla/sessions"
	"golang.org/x/crypto/bcrypt"
)

func Password(w http.ResponseWriter, r *http.Request, client *services.AppwriteClient, store *sessions.CookieStore) {
	if r.Method == http.MethodGet {
		data := models.PublicData{
			Title:   "Change Password",
			Data:    map[string]interface{}{},
			Error:   r.URL.Query().Get("error"),
			Msg:     r.URL.Query().Get("msg"),
			Session: models.GlobalSessionData,
		}

		utils.RenderTemplateWithSidebar(w, r, "views/templates/backend.html", "views/pages/users/password.html", data)
		return
	}

	if r.Method == http.MethodPost {
		oldPassword := r.FormValue("old")
		newPassword := r.FormValue("new")
		rePassword := r.FormValue("re")

		if oldPassword == "" || newPassword == "" || rePassword == "" {
			http.Redirect(w, r, "/app/password?error=form tidak lengkap", http.StatusSeeOther)
			return
		}

		if len(newPassword) < 8 {
			http.Redirect(w, r, "/app/password?error=password less than 8 character", http.StatusSeeOther)
			return
		}

		if newPassword != rePassword {
			http.Redirect(w, r, "/app/password?error=password not match", http.StatusSeeOther)
			return
		}

		existingUser, err := client.GetUserByID(os.Getenv("USERS"), models.GlobalSessionData.UserId)
		if err != nil {
			http.Redirect(w, r, "/app/password?error=user not found", http.StatusSeeOther)
			return
		}

		err = bcrypt.CompareHashAndPassword([]byte(existingUser.Password), []byte(oldPassword))
		if err != nil {
			http.Redirect(w, r, "/app/password?error=password not valid", http.StatusSeeOther)
			return
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
		if err != nil {
			http.Redirect(w, r, "/app/password?error=internal server error", http.StatusSeeOther)
			return
		}

		userChange := models.User{
			Name:          existingUser.Name,
			Password:      string(hashedPassword),
			EmailVerified: existingUser.EmailVerified,
			Role:          existingUser.Role,
			Email:         existingUser.Email,
		}

		err = client.UpdateUser(os.Getenv("USERS"), existingUser.ID, userChange)
		if err != nil {
			http.Redirect(w, r, "/app/password?error=internal server error "+err.Error(), http.StatusSeeOther)
			return
		}

		http.Redirect(w, r, "/app/signout?message=password was changed", http.StatusSeeOther)
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
