package db

import (
	"database/sql"
	_ "embed"
	"fmt"
	"log"

	"github.com/bxra2/taqnun/internal/normalize"
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
	if _, err := Conn.Exec(schemaSQL); err != nil {
		return err
	}
	return migrateDescriptionNorm()
}

// migrateDescriptionNorm is an additive migration that ensures the term table
// has a description_norm column and backfills it for any existing rows. It is
// idempotent: on subsequent runs no column is added and no NULL rows remain.
func migrateDescriptionNorm() error {
	has, err := columnExists("term", "description_norm")
	if err != nil {
		return err
	}
	if !has {
		if _, err := Conn.Exec(`ALTER TABLE term ADD COLUMN description_norm TEXT`); err != nil {
			return fmt.Errorf("add description_norm column: %w", err)
		}
	}

	tx, err := Conn.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	rows, err := tx.Query(`SELECT id, COALESCE(description, '') FROM term WHERE description_norm IS NULL`)
	if err != nil {
		return fmt.Errorf("select terms to backfill: %w", err)
	}

	type termRow struct {
		id   int64
		desc string
	}
	var pending []termRow
	for rows.Next() {
		var r termRow
		if err := rows.Scan(&r.id, &r.desc); err != nil {
			rows.Close()
			return fmt.Errorf("scan term to backfill: %w", err)
		}
		pending = append(pending, r)
	}
	if err := rows.Err(); err != nil {
		rows.Close()
		return fmt.Errorf("iterate terms to backfill: %w", err)
	}
	rows.Close()

	stmt, err := tx.Prepare(`UPDATE term SET description_norm = ? WHERE id = ?`)
	if err != nil {
		return fmt.Errorf("prepare backfill: %w", err)
	}
	defer stmt.Close()

	for _, r := range pending {
		if _, err := stmt.Exec(normalize.Text(r.desc), r.id); err != nil {
			return fmt.Errorf("backfill description_norm for term %d: %w", r.id, err)
		}
	}
	if err := stmt.Close(); err != nil {
		return err
	}
	return tx.Commit()
}

// columnExists reports whether the named column exists on the given table.
func columnExists(table, column string) (bool, error) {
	// table must be a trusted, hardcoded identifier (PRAGMA has no parameters).
	rows, err := Conn.Query(`PRAGMA table_info(` + table + `)`)
	if err != nil {
		return false, err
	}
	defer rows.Close()

	for rows.Next() {
		// PRAGMA table_info columns: cid, name, type, notnull, dflt_value, pk
		var cid int
		var name, typ string
		var notnull int
		var dflt, pk interface{}
		if err := rows.Scan(&cid, &name, &typ, &notnull, &dflt, &pk); err != nil {
			return false, err
		}
		if name == column {
			return true, nil
		}
	}
	return false, rows.Err()
}

func Close() {
	if Conn != nil {
		Conn.Close()
	}
}
