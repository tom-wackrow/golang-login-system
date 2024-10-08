package routes

import (
	"html/template"
	"log"
	"net/http"

	auth "github.com/tommytank20/login/auth"
	sockets "github.com/tommytank20/login/sockets"
)

func Dashboard(w http.ResponseWriter, r *http.Request) {
	tmpl, _ := template.ParseFiles("templates/base.html", "templates/dashboard.html")
	cookie, _ := r.Cookie("session")
	sessionToken := cookie.Value
	tmpl.Execute(w, sessionToken)
}

func Login(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		if auth.Auth(w, r) {
			http.Redirect(w, r, "/dashboard", 303)
		}
		http.Redirect(w, r, "/login", 303)
		return
	}

	tmpl, _ := template.ParseFiles("templates/base.html", "templates/login.html")
	tmpl.Execute(w, nil)
}

func Register(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		auth.Register(w, r)
		if auth.Auth(w, r) {
			http.Redirect(w, r, "/dashboard", 303)
		}
	}

	tmpl, _ := template.ParseFiles("templates/base.html", "templates/register.html")
	tmpl.Execute(w, nil)
}

func Refresh(w http.ResponseWriter, r *http.Request) {
	auth.RefreshAuth(w, r)
	http.Redirect(w, r, "/dashboard", 303)
}

func Logout(w http.ResponseWriter, r *http.Request) {
	auth.Logout(w, r)
	http.Redirect(w, r, "/login", 303)
}

func logWrapper(f http.HandlerFunc) http.HandlerFunc {
	return func (w http.ResponseWriter, r *http.Request) {
		f(w, r)

		uri := r.URL.String()
		method := r.Method

		log.Printf("%v %v", method, uri)
	}
}

func Run() {
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/", fs)
	http.HandleFunc("/socket", sockets.Handler)
	http.HandleFunc("/register", logWrapper(Register))
	http.HandleFunc("/login", logWrapper(Login))
	http.HandleFunc("/refresh", logWrapper(auth.RequireAuth(Refresh)))
	http.HandleFunc("/logout", logWrapper(auth.RequireAuth(Logout)))
	http.HandleFunc("/dashboard", logWrapper(auth.RequireAuth(Dashboard)))
	// http.ListenAndServeTLS(":80", "localhost.crt", "localhost.key", nil)
	http.ListenAndServe(":80", nil)
}