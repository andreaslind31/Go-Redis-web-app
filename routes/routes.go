package routes

import (
	"net/http"

	"github.com/andreaslind31/Go-Redis-web-app/middleware"
	"github.com/andreaslind31/Go-Redis-web-app/models"
	"github.com/andreaslind31/Go-Redis-web-app/sessions"
	"github.com/andreaslind31/Go-Redis-web-app/utils"

	"github.com/gorilla/mux"
)

func NewRouter() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/", middleware.AuthRequired(indexGetHandler)).Methods("GET")
	r.HandleFunc("/", middleware.AuthRequired(indexPostHandler)).Methods("POST")
	r.HandleFunc("/confirm", middleware.AuthRequired(confirmGetHandler)).Methods("GET")
	r.HandleFunc("/confirm", middleware.AuthRequired(confirmPostHandler)).Methods("POST")
	r.HandleFunc("/create", middleware.AdminRequired(createGetHandler)).Methods("GET")
	r.HandleFunc("/create", middleware.AdminRequired(createPostHandler)).Methods("POST")
	r.HandleFunc("/login", loginGetHandler).Methods("GET")
	r.HandleFunc("/login", loginPostHandler).Methods("POST")
	r.HandleFunc("/register", registerGetHandler).Methods("GET")
	r.HandleFunc("/register", registerPostHandler).Methods("POST")
	fs := http.FileServer(http.Dir("./assets/"))
	r.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", fs))
	return r
}

func indexGetHandler(w http.ResponseWriter, r *http.Request) {
	comments, err := models.GetComments()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal server error"))
		return
	}
	utils.ExecuteTemplate(w, "index.html", comments)
}
func indexPostHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	comment := r.PostForm.Get("comment")
	err := models.PostComment(comment)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal server error"))
		return
	}
	http.Redirect(w, r, "/", http.StatusFound)
}
func loginGetHandler(w http.ResponseWriter, r *http.Request) {
	utils.ExecuteTemplate(w, "login.html", nil)
}
func loginPostHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	username := r.PostForm.Get("username")
	password := r.PostForm.Get("password")
	err := models.AuthenticateUser(username, password)
	if err != nil {
		switch err {
		case models.ErrUserNotFound:
			utils.ExecuteTemplate(w, "login.html", "unknown user")
		case models.ErrInvalidLogin:
			utils.ExecuteTemplate(w, "login.html", "invalid login")
		default:
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Internal server error"))
		}
		return
	}

	session, _ := sessions.Store.Get(r, "session")
	session.Values["username"] = username

	if username == "admin" {
		session.Values["authenticated"] = true
	}

	session.Save(r, w)
	http.Redirect(w, r, "/", http.StatusFound)
}

// func logoutHandler(w http.ResponseWriter, r *http.Request) {
// 	session, _ := sessions.Store.Get(r, "username")

// 	// Revoke users authentication
// 	session.Values["username"] = nil
// 	session.Values["authenticated"] = false
// 	session.Save(r, w)
// }

func registerGetHandler(w http.ResponseWriter, r *http.Request) {
	utils.ExecuteTemplate(w, "register.html", nil)
}
func registerPostHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	username := r.PostForm.Get("username")
	password := r.PostForm.Get("password")
	err := models.RegisterUser(username, password)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal server error"))
		return
	}
	http.Redirect(w, r, "/login", http.StatusFound)
}
func confirmGetHandler(w http.ResponseWriter, r *http.Request) {
	utils.ExecuteTemplate(w, "confirm.html", nil)
}
func confirmPostHandler(w http.ResponseWriter, r *http.Request) {
	utils.ExecuteTemplate(w, "confirm.html", nil)
}
func createGetHandler(w http.ResponseWriter, r *http.Request) {
	utils.ExecuteTemplate(w, "create.html", nil)
}
func createPostHandler(w http.ResponseWriter, r *http.Request) {
	utils.ExecuteTemplate(w, "create.html", nil)
}
