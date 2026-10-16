package handlers

import (
	"log"

	"github.com/bxra2/taqnun/db"
	"github.com/bxra2/taqnun/models"
	"github.com/gofiber/fiber/v3"
)

func GetPublishers(c fiber.Ctx) error {
	pubRows, err := db.Conn.Query(`
		SELECT id, COALESCE(publisherEn,''), COALESCE(publisherAr,''), COALESCE(publisherUrl,'')
		FROM publisher`)
	if err != nil {
		log.Printf("GetPublishers query: %v", err)
		return fiber.NewError(fiber.StatusInternalServerError, "DB error")
	}
	defer pubRows.Close()

	publishers := make([]models.Publisher, 0)
	idx := make(map[int64]int) // publisher id -> position in slice

	for pubRows.Next() {
		var p models.Publisher
		if err := pubRows.Scan(&p.ID, &p.PublisherEn, &p.PublisherAr, &p.PublisherURL); err != nil {
			log.Printf("GetPublishers scan: %v", err)
			return fiber.NewError(fiber.StatusInternalServerError, "Decode error")
		}
		p.Glossaries = make([]models.Glossary, 0) // [] not null in JSON
		idx[p.ID] = len(publishers)
		publishers = append(publishers, p)
	}
	if err := pubRows.Err(); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "DB error")
	}

	glosRows, err := db.Conn.Query(`
		SELECT publisher_id, COALESCE(glossaryEn,''), COALESCE(glossaryAr,''),
		       COALESCE(glossaryUrl,''), COALESCE(license,''),
		       COALESCE(image_url,''), COALESCE(book_url,''), COALESCE(wikidata_url,'')
		FROM glossary`)
	if err != nil {
		log.Printf("GetPublishers glossaries: %v", err)
		return fiber.NewError(fiber.StatusInternalServerError, "DB error")
	}
	defer glosRows.Close()

	// Some glossaries are split across multiple source files but share the same
	// name and URL (e.g. one dictionary ingested in chunks); deduplicate them
	// so each logical glossary appears once per publisher.
	seen := make(map[int64]map[string]bool)
	for glosRows.Next() {
		var pubID int64
		var g models.Glossary
		if err := glosRows.Scan(&pubID, &g.GlossaryEn, &g.GlossaryAr, &g.GlossaryURL, &g.License,
			&g.ImageURL, &g.BookURL, &g.WikidataURL); err != nil {
			log.Printf("GetPublishers glossary scan: %v", err)
			return fiber.NewError(fiber.StatusInternalServerError, "Decode error")
		}
		i, ok := idx[pubID]
		if !ok {
			continue
		}
		key := g.GlossaryEn + "\x00" + g.GlossaryAr + "\x00" + g.GlossaryURL
		if seen[pubID] == nil {
			seen[pubID] = make(map[string]bool)
		}
		if seen[pubID][key] {
			continue
		}
		seen[pubID][key] = true
		publishers[i].Glossaries = append(publishers[i].Glossaries, g)
	}
	if err := glosRows.Err(); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "DB error")
	}

	return c.JSON(publishers)
}
