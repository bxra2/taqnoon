package db

import (
	"database/sql"
	_ "embed"
	"fmt"
	"log"

	"github.com/bxra2/taqnun/internal/metadata"
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
	if err := migrateDescriptionNorm(); err != nil {
		return err
	}
	return migrateGlossaryMeta()
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

// migrateGlossaryMeta adds the image_url/book_url/wikidata_url columns to the
// glossary table (for pre-existing databases) and backfills them from the
// glossary URL. It is idempotent and network-free.
func migrateGlossaryMeta() error {
	for _, col := range []string{"image_url", "book_url", "wikidata_url"} {
		has, err := columnExists("glossary", col)
		if err != nil {
			return err
		}
		if !has {
			if _, err := Conn.Exec(`ALTER TABLE glossary ADD COLUMN ` + col + ` TEXT`); err != nil {
				return fmt.Errorf("add glossary %s: %w", col, err)
			}
		}
	}

	rows, err := Conn.Query(`SELECT id, COALESCE(glossaryUrl, '') FROM glossary`)
	if err != nil {
		return fmt.Errorf("select glossary for backfill: %w", err)
	}
	defer rows.Close()

	type row struct {
		id  int64
		url string
	}
	var pending []row
	for rows.Next() {
		var r row
		if err := rows.Scan(&r.id, &r.url); err != nil {
			return fmt.Errorf("scan glossary for backfill: %w", err)
		}
		pending = append(pending, r)
	}
	if err := rows.Err(); err != nil {
		return fmt.Errorf("iterate glossary for backfill: %w", err)
	}

	for _, r := range pending {
		m := metadata.Resolve(r.url)
		if _, err := Conn.Exec(`
			UPDATE glossary SET
				image_url = COALESCE(NULLIF(image_url, ''), ?),
				book_url = COALESCE(NULLIF(book_url, ''), ?),
				wikidata_url = COALESCE(NULLIF(wikidata_url, ''), ?)
			WHERE id = ?`,
			m.ImageURL, m.BookURL, m.WikidataURL, r.id); err != nil {
			return fmt.Errorf("backfill glossary %d: %w", r.id, err)
		}
	}
	return nil
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
