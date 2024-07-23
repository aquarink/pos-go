package controllers

import (
	"fmt"
	"net/http"
	"os"
	"pos/models"
	"pos/services"
	"pos/utils"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"golang.org/x/crypto/bcrypt"
)

func SignupController(w http.ResponseWriter, r *http.Request, client *services.AppwriteClient, store *sessions.CookieStore) {
	if r.Method == http.MethodGet {
		data := models.PublicData{
			Title: "Sign Up",
			Data:  map[string]interface{}{},
			Error: r.URL.Query().Get("error"),
			Msg:   r.URL.Query().Get("msg"),
		}
		utils.RenderTemplate(w, r, "views/templates/auth.html", "views/pages/auth/signup.html", data)
		return
	}

	if r.Method == http.MethodPost {
		name := r.FormValue("name")
		email := r.FormValue("email")
		password := r.FormValue("password")
		turnstileToken := r.FormValue("cf-turnstile-response")

		err := utils.VerifyTurnstile(turnstileToken)
		if err != nil {
			http.Redirect(w, r, "/app/signup?error=validasi gagal", http.StatusSeeOther)
			return
		} else {
			if name == "" || email == "" || password == "" {
				http.Redirect(w, r, "/app/signup?error=form tidak lengkap", http.StatusSeeOther)
				return
			}

			if len(password) < 8 {
				http.Redirect(w, r, "/app/signup?error=password kurang dari 8 karakter", http.StatusSeeOther)
				return
			}

			existingUser, _ := client.GetUserByEmail(os.Getenv("USERS"), email)
			if existingUser != nil {
				http.Redirect(w, r, "/app/signup?error=email sudah ada", http.StatusSeeOther)
				return
			}

			hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
			if err != nil {
				http.Redirect(w, r, "/app/signup?error=internal server failed", http.StatusSeeOther)
				return
			}

			user := models.User{
				Name:     name,
				Email:    email,
				Password: string(hashedPassword),
				Role:     models.RoleMerchant,
			}

			userID, err := client.CreateUser(os.Getenv("USERS"), user)
			if err != nil {
				http.Redirect(w, r, "/app/signup?error=internal server error "+err.Error(), http.StatusSeeOther)
				return
			}

			// KIRIM EMAIL
			subject := "Email Verification"
			text := fmt.Sprintf("Hi %s,\n\nThank you for registering with us.", name)
			html := fmt.Sprintf("Hi %s,<br><br>Thank you for registering with us.<br>Click <a href='%s%s'>here</a> to verify your email.", name, os.Getenv("EMAIL_VERIFY_URL"), userID)

			err = utils.SendEmail(email, subject, text, html)
			if err != nil {
				http.Redirect(w, r, "/app/signup?error=gagal mengirim email verifikasi", http.StatusSeeOther)
				return
			}

			// MODEL MAILS
			emailDoc := models.Mails{
				UserID:  user.ID,
				Email:   email,
				Subject: subject,
				Text:    text,
				HTML:    html,
			}

			err = client.CreateEmail(os.Getenv("MAILS"), emailDoc)
			if err != nil {
				http.Redirect(w, r, "/app/signup?error=internal server fails", http.StatusSeeOther)
				return
			}

			http.Redirect(w, r, "/app/signin?message=silahkan cek email anda untuk verifikasi", http.StatusSeeOther)
		}
	}
}

func SignupVerifyController(w http.ResponseWriter, r *http.Request, client *services.AppwriteClient, store *sessions.CookieStore) {
	if r.Method == http.MethodGet {
		vars := mux.Vars(r)
		userID := vars["id"]

		// Fetch user by ID
		user, err := client.GetUserByID(os.Getenv("USERS"), userID)
		if err != nil || user == nil {
			http.Redirect(w, r, "/app/signin?error=invalid or expired link", http.StatusSeeOther)
			return
		}

		if user.EmailVerified {
			http.Redirect(w, r, "/app/signin?error=link expired", http.StatusSeeOther)
			return
		}

		// Update user's email_verified to true
		err = client.VerifyUserEmail(os.Getenv("USERS"), userID)
		if err != nil {
			http.Redirect(w, r, "/app/signin?error=failed to verify email", http.StatusSeeOther)
			return
		}

		// Check and update cashier status if the role is "cashier"
		if user.Role == models.RoleCashier {
			err = client.UpdateCashierStatus(os.Getenv("CASHIERS"), userID, "active")
			if err != nil {
				http.Redirect(w, r, "/app/signin?error=failed to update cashier status", http.StatusSeeOther)
				return
			}
		}

		http.Redirect(w, r, "/app/signin?msg=Email verification successful", http.StatusSeeOther)
		return
	}
}

func SigninController(w http.ResponseWriter, r *http.Request, client *services.AppwriteClient, store *sessions.CookieStore) {
	if r.Method == http.MethodGet {
		data := models.PublicData{
			Title: "Sign In",
			Data:  map[string]interface{}{},
			Error: r.URL.Query().Get("error"),
			Msg:   r.URL.Query().Get("msg"),
		}

		utils.RenderTemplate(w, r, "views/templates/auth.html", "views/pages/auth/signin.html", data)
		return
	}

	if r.Method == http.MethodPost {
		email := r.FormValue("email")
		password := r.FormValue("password")
		turnstileToken := r.FormValue("cf-turnstile-response")

		err := utils.VerifyTurnstile(turnstileToken)
		if err != nil {
			http.Redirect(w, r, "/app/signup?error=validasi gagal", http.StatusSeeOther)
			return
		} else {
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

			if !user.EmailVerified {
				http.Redirect(w, r, "/app/signin?error=email belum diverifikasi", http.StatusSeeOther)
				return
			}

			// SESSION KAYA $_SESSION NYA PHP
			ses, _ := store.Get(r, "session")
			ses.Values["user_id"] = user.ID
			ses.Values["user_name"] = user.Name
			ses.Values["role"] = user.Role

			err = ses.Save(r, w)
			if err != nil {
				http.Redirect(w, r, "/app/signin?error=perbaikan sistem, mohon coba lagi nant", http.StatusSeeOther)
				return
			}

			if user.Role == "cashier" || user.Role == "Cashier" {
				http.Redirect(w, r, "/app/order", http.StatusSeeOther)
			} else {
				http.Redirect(w, r, "/app/dashboard", http.StatusSeeOther)
			}
		}
	}
}

func DashboardController(w http.ResponseWriter, r *http.Request, client *services.AppwriteClient, store *sessions.CookieStore) {
	if r.Method == http.MethodGet {
		data := models.PublicData{
			Title:   "Dashboard",
			Data:    map[string]interface{}{},
			Error:   r.URL.Query().Get("error"),
			Msg:     r.URL.Query().Get("msg"),
			Session: models.GlobalSessionData,
		}

		utils.RenderTemplateWithSidebar(w, r, "views/templates/backend.html", "views/pages/dashboard/dashboard.html", data)
		return
	}
	http.Redirect(w, r, "/app/dashboard", http.StatusSeeOther)
}

func SignoutController(w http.ResponseWriter, r *http.Request, client *services.AppwriteClient, store *sessions.CookieStore) {
	session, _ := store.Get(r, "session")
	delete(session.Values, "user_id")
	session.Save(r, w)

	msg := r.URL.Query().Get("msg")
	if msg != "" {
		http.Redirect(w, r, "/app/signin?msg="+msg, http.StatusSeeOther)
	} else {
		http.Redirect(w, r, "/app/signin", http.StatusSeeOther)
	}
}
