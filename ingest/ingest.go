package ingest

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

// FileData mirrors the "fileData" object in each glossary JSON file.
// Some files use a typo'd key "glossarysAr" instead of "glossaryAr".
type FileData struct {
	GlossaryEn   string `json:"glossaryEn"`
	GlossaryAr   string `json:"glossaryAr"`
	GlossaryURL  string `json:"glossaryUrl"`
	PublisherEn  string `json:"publisherEn"`
	PublisherAr  string `json:"publisherAr"`
	PublisherURL string `json:"publisherUrl"`
}

// Entry mirrors one element of the "entries" array. Fields are optional
// across sources, so missing keys simply unmarshal to "".
type Entry struct {
	English string `json:"english"`
	Arabic  string `json:"arabic"`
	French  string `json:"french"`
	German  string `json:"german"`
	TURL    string `json:"tURL"`
	Desc    string `json:"desc"`
	DescAlt string `json:"description"`
}

func (e Entry) description() string {
	if e.Desc != "" {
		return e.Desc
	}
	return e.DescAlt
}

type glossaryFile struct {
	FileData FileData `json:"fileData"`
	Entries  []Entry  `json:"entries"`
}

// Run ingests every file listed in <dataDir>/index.json. It is idempotent:
// files whose SHA-256 hash already matches source_file are skipped. It logs a
// summary and returns nil if the manifest is missing (so startup never fails
// when no glossary data is present).
func Run(db *sql.DB, dataDir string) error {
	manifestPath := filepath.Join(dataDir, "index.json")
	manifest, err := loadManifest(manifestPath)
	if err != nil {
		if os.IsNotExist(err) {
			log.Printf("[ingest] no manifest at %s, skipping ingestion", manifestPath)
			return nil
		}
		return err
	}

	log.Printf("[ingest] found %d files in manifest", len(manifest))
	var found, ingested, skipped, failed int
	for _, rel := range manifest {
		found++
		abs := filepath.Join(dataDir, rel)
		n, err := ingestFile(db, abs, rel)
		if err != nil {
			failed++
			log.Printf("[ingest] FAILED %s: %v", rel, err)
			continue
		}
		if n == 0 {
			skipped++
		} else {
			ingested++
			log.Printf("[ingest] ingested %s (%d terms)", rel, n)
		}
	}
	log.Printf("[ingest] done: %d found, %d ingested, %d skipped, %d failed", found, ingested, skipped, failed)
	return nil
}

func loadManifest(path string) ([]string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var rels []string
	if err := json.Unmarshal(data, &rels); err != nil {
		return nil, fmt.Errorf("parse manifest %s: %w", path, err)
	}
	return rels, nil
}

// ingestFile ingests a single glossary file, returning the number of terms
// written (0 means unchanged and skipped).
func ingestFile(db *sql.DB, absPath, relPath string) (int, error) {
	data, err := os.ReadFile(absPath)
	if err != nil {
		return 0, fmt.Errorf("read: %w", err)
	}
	sum := sha256.Sum256(data)
	hash := hex.EncodeToString(sum[:])

	var glID int64
	var prevHash string
	err = db.QueryRow(`SELECT glossary_id, file_hash FROM source_file WHERE path = ?`, relPath).Scan(&glID, &prevHash)
	switch {
	case err == nil && prevHash == hash:
		return 0, nil // unchanged
	case err != nil && err != sql.ErrNoRows:
		return 0, fmt.Errorf("source_file lookup: %w", err)
	}

	var gf glossaryFile
	if err := json.Unmarshal(data, &gf); err != nil {
		return 0, fmt.Errorf("parse json: %w", err)
	}

	tx, err := db.Begin()
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	pubID, err := upsertPublisher(tx, gf.FileData)
	if err != nil {
		return 0, err
	}

	if glID == 0 {
		res, err := tx.Exec(`
			INSERT INTO glossary (glossaryEn, glossaryAr, glossaryUrl, publisher_id, domain_id, license, imported_at)
			VALUES (?, ?, ?, ?, NULL, NULL, datetime('now'))`,
			gf.FileData.GlossaryEn, gf.FileData.GlossaryAr, gf.FileData.GlossaryURL, pubID)
		if err != nil {
			return 0, fmt.Errorf("insert glossary: %w", err)
		}
		glID, err = res.LastInsertId()
		if err != nil {
			return 0, err
		}
	} else {
		if _, err := tx.Exec(`
			UPDATE glossary SET glossaryEn = ?, glossaryAr = ?, glossaryUrl = ?, publisher_id = ?
			WHERE id = ?`,
			gf.FileData.GlossaryEn, gf.FileData.GlossaryAr, gf.FileData.GlossaryURL, pubID, glID); err != nil {
			return 0, fmt.Errorf("update glossary: %w", err)
		}
		if _, err := tx.Exec(`DELETE FROM term WHERE glossary_id = ?`, glID); err != nil {
			return 0, fmt.Errorf("delete old terms: %w", err)
		}
	}

	stmt, err := tx.Prepare(`
		INSERT INTO term (glossary_id, english, arabic, french, german, term_url, description, description_norm,
		                  english_norm, arabic_norm, french_norm, german_norm)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	for _, e := range gf.Entries {
		if _, err := stmt.Exec(glID, e.English, e.Arabic, e.French, e.German, e.TURL, e.description(), normText(e.description()),
			normLatin(e.English), normArabic(e.Arabic), normLatin(e.French), normLatin(e.German)); err != nil {
			return 0, fmt.Errorf("insert term: %w", err)
		}
	}

	if _, err := tx.Exec(`
		INSERT INTO source_file (path, file_hash, glossary_id) VALUES (?, ?, ?)
		ON CONFLICT(path) DO UPDATE SET file_hash = excluded.file_hash, glossary_id = excluded.glossary_id, ingested_at = datetime('now')`,
		relPath, hash, glID); err != nil {
		return 0, fmt.Errorf("upsert source_file: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return 0, err
	}
	return len(gf.Entries), nil
}

func upsertPublisher(tx *sql.Tx, fd FileData) (int64, error) {
	var id int64
	err := tx.QueryRow(`SELECT id FROM publisher WHERE publisherEn = ? AND publisherAr = ?`,
		fd.PublisherEn, fd.PublisherAr).Scan(&id)
	if err == nil {
		if _, err := tx.Exec(`UPDATE publisher SET publisherUrl = ? WHERE id = ?`, fd.PublisherURL, id); err != nil {
			return 0, err
		}
		return id, nil
	}
	if err != sql.ErrNoRows {
		return 0, err
	}
	res, err := tx.Exec(`INSERT INTO publisher (publisherEn, publisherAr, publisherUrl) VALUES (?, ?, ?)`,
		fd.PublisherEn, fd.PublisherAr, fd.PublisherURL)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}
