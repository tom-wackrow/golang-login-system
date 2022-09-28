package main

import (
	"net/http"

	auth "github.com/tommytank20/login/auth"
	routes "github.com/tommytank20/login/routes"
)


func main() {
	http.HandleFunc("/login", routes.Login)
	http.HandleFunc("/dashboard", auth.RequireAuth(routes.Dashboard))
	http.ListenAndServe(":80", nil)
}