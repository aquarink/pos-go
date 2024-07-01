package middleware

import (
	"net/http"
	"pos/models"

	"github.com/gorilla/sessions"
)

var store = sessions.NewCookieStore([]byte("something-very-secret"))

// ngecek apakah user masih login atau masih ada session
// pada selain auth page
func CheckSignin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, _ := store.Get(r, "session")
		if session.Values["user_id"] == nil {
			http.Redirect(w, r, "/app/signin", http.StatusSeeOther)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// ngecek apakah user masih login atau masih ada session
// pada auth page
func CheckSession(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, _ := store.Get(r, "session")
		if session.Values["user_id"] != nil {
			http.Redirect(w, r, "/app/dashboard", http.StatusSeeOther)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// ini supaya sessionnya bisa di akses global
func SessionMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, _ := store.Get(r, "session")
		role, _ := session.Values["role"].(string)
		lastLogin, _ := session.Values["last_login"].(string)

		models.GlobalSessionData = models.SessionData{
			Role:      role,
			LastLogin: lastLogin,
		}

		next.ServeHTTP(w, r)
	})
}
