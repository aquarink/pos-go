package controllers

import (
	"log"
	"net/http"
	"os"
	"pos/models"
	"pos/services"
	"pos/utils"

	"github.com/gorilla/sessions"
	"golang.org/x/crypto/bcrypt"
)

var store = sessions.NewCookieStore([]byte("something-very-secret"))

func SignupController(w http.ResponseWriter, r *http.Request, client *services.AppwriteClient) {
	if r.Method == http.MethodGet {
		utils.RenderTemplate(w, "views/templates/auth.html", "views/pages/auth/signup.html", nil)
		return
	}

	if r.Method == http.MethodPost {
		name := r.FormValue("name")
		email := r.FormValue("email")
		password := r.FormValue("password")

		if name == "" || email == "" || password == "" {
			http.Redirect(w, r, "/app/signup?error=form tidak lengkap", http.StatusSeeOther)
			return
		}

		if len(password) < 8 {
			http.Redirect(w, r, "/app/signup?error=password kurang dari 8 karakter", http.StatusSeeOther)
			return
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			log.Println("Error generating password hash:", err)
			http.Redirect(w, r, "/app/signup?error=internal server error", http.StatusSeeOther)
			return
		}

		user := models.User{
			Name:     name,
			Email:    email,
			Password: string(hashedPassword),
		}

		err = client.CreateUser(os.Getenv("USERS"), user)
		if err != nil {
			log.Println("Error creating user:", err)
			http.Redirect(w, r, "/app/signup?error=internal server error", http.StatusSeeOther)
			return
		}

		http.Redirect(w, r, "/app/signin?message=silahkan cek email anda untuk verifikasi", http.StatusSeeOther)
	}
}

func SigninController(w http.ResponseWriter, r *http.Request, client *services.AppwriteClient) {
	if r.Method == http.MethodGet {
		utils.RenderTemplate(w, "views/templates/auth.html", "views/pages/auth/signin.html", nil)
		return
	}

	if r.Method == http.MethodPost {
		email := r.FormValue("email")
		password := r.FormValue("password")

		if email == "" || password == "" {
			http.Redirect(w, r, "/app/signin?error=form tidak lengkap", http.StatusSeeOther)
			return
		}

		user, err := client.GetUserByEmail(os.Getenv("USERS"), email)
		if err != nil {
			http.Redirect(w, r, "/app/signin?error=email atau password salah", http.StatusSeeOther)
			return
		}

		err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
		if err != nil {
			http.Redirect(w, r, "/app/signin?error=email atau password salah", http.StatusSeeOther)
			return
		}

		// Set session
		session, _ := store.Get(r, "session")
		session.Values["user_id"] = user.ID
		session.Save(r, w)

		http.Redirect(w, r, "/app/dashboard", http.StatusSeeOther)
	}
}
