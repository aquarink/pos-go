package controllers

import (
	"net/http"
	"pos/models"
	"pos/services"
	"text/template"

	"github.com/gorilla/sessions"
	"golang.org/x/crypto/bcrypt"
)

var store = sessions.NewCookieStore([]byte("something-very-secret"))

func SignupController(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		tmpl, err := template.ParseFiles(
			"views/templates/auth.html",
			"views/templates/auth.html",
		)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		tmpl.ExecuteTemplate(w, "layout", nil)
		return
	}

	if r.Method == http.MethodPost {
		name := r.FormValue("name")
		email := r.FormValue("email")
		password := r.FormValue("password")

		if name == "" || email == "" || password == "" {
			http.Redirect(w, r, "/app/signup?error=Form tidak lengkap", http.StatusSeeOther)
			return
		}

		if len(password) < 8 {
			http.Redirect(w, r, "/app/signup?error=Password kurang dari 8 karakter", http.StatusSeeOther)
			return
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		user := models.User{
			Name:     name,
			Email:    email,
			Password: string(hashedPassword),
		}

		err = services.CreateUser(user)
		if err != nil {
			if err.Error() == "Email already exists" {
				http.Redirect(w, r, "/app/signup?error=Email sudah digunakan", http.StatusSeeOther)
				return
			}
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/app/signin?message=Silahkan cek email anda untuk verifikasi", http.StatusSeeOther)
	}
}

func SigninController(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		tmpl, err := template.ParseFiles(
			"views/templates/auth.html",
			"views/pages/auth/signin.html",
		)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		tmpl.ExecuteTemplate(w, "layout", nil)
		return
	}

	if r.Method == http.MethodPost {
		email := r.FormValue("email")
		password := r.FormValue("password")

		if email == "" || password == "" {
			http.Redirect(w, r, "/app/signin?error=Form tidak lengkap", http.StatusSeeOther)
			return
		}

		user, err := services.GetUserByEmail(email)
		if err != nil {
			http.Redirect(w, r, "/app/signin?error=Email atau password salah", http.StatusSeeOther)
			return
		}

		err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
		if err != nil {
			http.Redirect(w, r, "/app/signin?error=Email atau password salah", http.StatusSeeOther)
			return
		}

		// Set session
		session, _ := store.Get(r, "session")
		session.Values["user_id"] = user.ID
		session.Save(r, w)

		http.Redirect(w, r, "/app/dashboard", http.StatusSeeOther)
	}
}

func SignoutController(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session")
	delete(session.Values, "user_id")
	session.Save(r, w)
	http.Redirect(w, r, "/app/signin", http.StatusSeeOther)
}
