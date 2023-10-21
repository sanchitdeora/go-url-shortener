package main

import (
	"net/http"

	"github.com/rs/zerolog/log"
	"github.com/sanchitdeora/db"
	"github.com/sanchitdeora/urlshort"
)

func main() {
	// database initialization
	db := db.NewDatabase()
	db.Init()

	// routers initialization
	routes()

    log.Info().Msg("URL Shortener is listening on :8000")
	http.ListenAndServe(":8000", nil)
}

func routes() {
	http.Handle("/", http.FileServer(http.Dir("./static")))
	http.HandleFunc("/shorten", urlshort.HandleUrlShortener)
}