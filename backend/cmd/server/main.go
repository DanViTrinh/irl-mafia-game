package main

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
    db, err := sql.Open("sqlite3", "./db/theflu.db")
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

    // Optional: test connection
    err = db.Ping()
    if err != nil {
        log.Fatal(err)
    }

    log.Println("Database connected successfully!")
}
