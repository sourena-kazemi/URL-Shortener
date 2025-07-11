package main

import (
	"URL-Shortener/internal/handlers"
	"database/sql"
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	db, err := sql.Open("sqlite3", "./internal/db/app.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	http.HandleFunc("GET /favicon.ico", handlers.FavIconHandler)
	http.HandleFunc("GET /static/", handlers.StaticFilesHandler)
	http.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		handlers.HomePageHandler(w, r)
	})
	http.HandleFunc("GET /dashboard", handlers.DashboardPageHandler)
	http.HandleFunc("POST /shorten/{url}", handlers.ShortenUrl)
	http.HandleFunc("GET /auth/github/callback", handlers.OAuthCallback)

	http.ListenAndServe(":8080", nil)
}
