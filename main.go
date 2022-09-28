package main

import (
	routes "github.com/tommytank20/login/routes"
)


func main() {
	// http.HandleFunc("/login", routes.Login)
	// http.HandleFunc("/refresh", routes.Refresh)
	// http.HandleFunc("/logout", routes.Logout)
	// http.HandleFunc("/dashboard", auth.RequireAuth(routes.Dashboard))
	// http.ListenAndServe(":80", nil)
	routes.Run()
}