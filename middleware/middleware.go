package middleware

import (
	"net/http"

	"github.com/andreaslind31/Go-Redis-web-app/sessions"
)

func AdminRequired(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, _ := sessions.Store.Get(r, "session")

		// Check if user is authenticated - if not send to login page
		if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}

		handler.ServeHTTP(w, r)
	}
}

func AuthRequired(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, _ := sessions.Store.Get(r, "session")

		// Check if user has username - if not send to login page
		_, ok := session.Values["username"]
		if !ok {
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}

		handler.ServeHTTP(w, r)
	}
}