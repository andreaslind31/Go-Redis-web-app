package main

import (
	"html/template"
	"net/http"

	"github.com/go-redis/redis"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"golang.org/x/crypto/bcrypt"
)

var store = sessions.NewCookieStore(([]byte("paZZw0rd123")))
var client *redis.Client
var templates *template.Template

func main() {
	client = redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	templates = template.Must(template.ParseGlob("templates/*.html"))
	r := mux.NewRouter()
	r.HandleFunc("/", authRequired(indexGetHandler)).Methods("GET")
	r.HandleFunc("/", authRequired(indexPostHandler)).Methods("POST")
	r.HandleFunc("/confirm", authRequired(confirmGetHandler)).Methods("GET")
	r.HandleFunc("/confirm", authRequired(confirmPostHandler)).Methods("POST")
	r.HandleFunc("/create", adminRequired(createGetHandler)).Methods("GET")
	r.HandleFunc("/create", adminRequired(createPostHandler)).Methods("POST")
	r.HandleFunc("/login", loginGetHandler).Methods("GET")
	r.HandleFunc("/login", loginPostHandler).Methods("POST")
	r.HandleFunc("/register", registerGetHandler).Methods("GET")
	r.HandleFunc("/register", registerPostHandler).Methods("POST")
	fs := http.FileServer(http.Dir("./assets/"))
	r.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", fs))
	http.Handle("/", r)
	http.ListenAndServe(":8080", nil)
}
func adminRequired(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, _ := store.Get(r, "session")

		// Check if user is authenticated - if not send to login page
		if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}

		handler.ServeHTTP(w, r)
	}
}

func authRequired(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, _ := store.Get(r, "session")

		// // Check if user has username - if not send to login page
		_, ok := session.Values["username"]
		if !ok {
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}

		handler.ServeHTTP(w, r)
	}
}
func indexGetHandler(w http.ResponseWriter, r *http.Request) {
	comments, err := client.LRange("comments", 0, 10).Result()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal server error"))
		return
	}
	templates.ExecuteTemplate(w, "index.html", comments)
}
func indexPostHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	comment := r.PostForm.Get("comment")
	err := client.LPush("comments", comment).Err()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal server error"))
		return
	}
	http.Redirect(w, r, "/", http.StatusFound)
}
func loginGetHandler(w http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(w, "login.html", nil)
}
func loginPostHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	username := r.PostForm.Get("username")
	password := r.PostForm.Get("password")
	hash, err := client.Get("user:" + username).Bytes()
	if err == redis.Nil {
		templates.ExecuteTemplate(w, "login.html", "unknown user")
		return
	} else if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal server error"))
		return
	}
	err = bcrypt.CompareHashAndPassword(hash, []byte(password))
	if err != nil {
		templates.ExecuteTemplate(w, "login.html", "invalid login")
		return
	}
	session, _ := store.Get(r, "session")
	session.Values["username"] = username

	if username == "admin" {
		session.Values["authenticated"] = true
	}

	session.Save(r, w)
	http.Redirect(w, r, "/", http.StatusFound)
}

// func logoutHandler(w http.ResponseWriter, r *http.Request) {
// 	session, _ := store.Get(r, "username")

// 	// Revoke users authentication
// 	session.Values["username"] = nil
// 	session.Values["authenticated"] = false
// 	session.Save(r, w)
// }

func registerGetHandler(w http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(w, "register.html", nil)
}
func registerPostHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	username := r.PostForm.Get("username")
	password := r.PostForm.Get("password")
	cost := bcrypt.DefaultCost
	hash, err := bcrypt.GenerateFromPassword([]byte(password), cost)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal server error"))
		return
	}
	err = client.Set("user:"+username, hash, 0).Err() //the 0 is to tell Redis that the key shouldnt expire
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal server error"))
		return
	}
	http.Redirect(w, r, "/login", http.StatusFound)
}
func confirmGetHandler(w http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(w, "confirm.html", nil)
}
func confirmPostHandler(w http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(w, "confirm.html", nil)
}
func createGetHandler(w http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(w, "create.html", nil)
}
func createPostHandler(w http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(w, "create.html", nil)
}
