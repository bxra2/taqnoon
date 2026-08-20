package handlers

import (
	"log"

	"github.com/bxra2/taqnun/db"
	"github.com/bxra2/taqnun/models"
	"github.com/gofiber/fiber/v3"
)

func GetTerms(c fiber.Ctx) error {
	rows, err := db.Conn.Query(`
		SELECT t.id, COALESCE(t.english,''), COALESCE(t.arabic,''),
		       COALESCE(t.french,''), COALESCE(t.german,''),
		       COALESCE(t.term_url,''), COALESCE(t.description,''),
		       COALESCE(g.glossaryEn,''), COALESCE(g.glossaryAr,''), COALESCE(g.glossaryUrl,''),
		       COALESCE(g.publisher_id, 0)
		FROM term t
		JOIN glossary g ON g.id = t.glossary_id`)
	if err != nil {
		log.Printf("GetTerms query: %v", err)
		return fiber.NewError(fiber.StatusInternalServerError, "DB error")
	}
	defer rows.Close()

	terms := make([]models.Term, 0)
	for rows.Next() {
		var t models.Term
		if err := rows.Scan(&t.ID, &t.English, &t.Arabic, &t.French, &t.German,
			&t.TURL, &t.Description,
			&t.GlossaryEn, &t.GlossaryAr, &t.GlossaryURL, &t.PublisherID); err != nil {
			log.Printf("GetTerms scan: %v", err)
			return fiber.NewError(fiber.StatusInternalServerError, "Decode error")
		}
		terms = append(terms, t)
	}
	if err := rows.Err(); err != nil {
		log.Printf("GetTerms rows: %v", err)
		return fiber.NewError(fiber.StatusInternalServerError, "DB error")
	}
	return c.JSON(terms)
}
