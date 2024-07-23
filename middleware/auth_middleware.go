package middleware

import (
	"net/http"
	"pos/models"

	"github.com/gorilla/sessions"
)

// ngecek apakah user masih login atau masih ada session
// pada selain auth page
func CheckSignin(store *sessions.CookieStore) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			session, _ := store.Get(r, "session")
			if session.Values["user_id"] == nil {
				http.Redirect(w, r, "/app/signin", http.StatusSeeOther)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

// ngecek apakah user masih login atau masih ada session
// pada auth page
func CheckSession(store *sessions.CookieStore) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			session, _ := store.Get(r, "session")
			if session.Values["user_id"] != nil {
				http.Redirect(w, r, "/app/dashboard", http.StatusSeeOther)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

// ini supaya sessionnya bisa di akses global
func SessionMiddleware(store *sessions.CookieStore) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			session, _ := store.Get(r, "session")
			user_id, _ := session.Values["user_id"].(string)
			user_name, _ := session.Values["user_name"].(string)
			role, _ := session.Values["role"].(string)

			models.GlobalSessionData = models.SessionData{
				UserId:   user_id,
				UserName: user_name,
				Role:     role,
			}

			next.ServeHTTP(w, r)
		})
	}
}
