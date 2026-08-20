package db

import (
	"database/sql"
	_ "embed"
	"log"

	_ "modernc.org/sqlite"
)

var Conn *sql.DB

//go:embed schema.sql
var schemaSQL string

func Connect() {
	var err error
	Conn, err = sql.Open("sqlite", "./data.db")
	if err != nil {
		log.Fatal(err)
	}
	Conn.SetMaxOpenConns(1)

	// SQLite disables FK enforcement by default; turn it on for ON DELETE CASCADE.
	if _, err = Conn.Exec("PRAGMA foreign_keys = ON"); err != nil {
		log.Fatal(err)
	}
	if err = Conn.Ping(); err != nil {
		log.Fatal(err)
	}
	log.Println("Connected to SQLite database")
}

func Migrate() error {
	_, err := Conn.Exec(schemaSQL)
	return err
}

func Close() {
	if Conn != nil {
		Conn.Close()
	}
}
