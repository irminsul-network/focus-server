package main

import (
	"database/sql"
	"log"
	_ "modernc.org/sqlite"
	"net/http"
	"time"
)

type App struct {
	db     *sql.DB
	server *http.Server
}

func buildApp() App {

	s := &http.Server{
		Addr:         ":8080",
		Handler:      nil,
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 1 * time.Second,
	}

	db, err := sql.Open("sqlite", "data/focus.db")
	if err != nil {
		log.Fatal(err)
	}

	db.SetConnMaxLifetime(-1)
	db.SetMaxIdleConns(3)
	db.SetMaxOpenConns(5)

	ensureSchema(db)

	return App{
		db:     db,
		server: s,
	}
}

func ensureSchema(db *sql.DB) {
	_, err := db.Exec("")
	if err != nil {
		log.Fatal("could not ensure schemas")
	}
}

func main() {
	app := buildApp()
	registerHandlers(&app)
	log.Fatal(app.server.ListenAndServe())

}
