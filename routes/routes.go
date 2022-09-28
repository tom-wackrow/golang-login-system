package routes

import (
	"html/template"
	"net/http"

	auth "github.com/tommytank20/login/auth"
)


func Login(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		if auth.Auth(w, r){
			http.Redirect(w, r, "/dashboard", 303)
		}
		return
	}

	tmpl, _ := template.ParseFiles("templates/login.html")
	tmpl.Execute(w, nil)
}

func Dashboard(w http.ResponseWriter, r *http.Request) {
	tmpl, _ := template.ParseFiles("templates/login.html")
	tmpl.Execute(w, nil)
}