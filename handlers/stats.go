package handlers

import (
	"log"

	"github.com/bxra2/taqnun/db"
	"github.com/bxra2/taqnun/models"
	"github.com/gofiber/fiber/v3"
)

func GetStats(c fiber.Ctx) error {
	c.Set("Cache-Control", "private, max-age=300")

	var stats models.Stats
	err := db.Conn.QueryRow(`
		SELECT COUNT(*),
		       COUNT(CASE WHEN english_norm <> '' THEN 1 END),
		       COUNT(CASE WHEN arabic_norm <> '' THEN 1 END),
		       COUNT(CASE WHEN french_norm <> '' THEN 1 END),
		       COUNT(CASE WHEN german_norm <> '' THEN 1 END)
		FROM term`).Scan(&stats.TotalTerms, &stats.English, &stats.Arabic, &stats.French, &stats.German)
	if err != nil {
		log.Printf("GetStats term query: %v", err)
		return fiber.NewError(fiber.StatusInternalServerError, "DB error")
	}

	if err := db.Conn.QueryRow(`SELECT COUNT(*) FROM glossary`).Scan(&stats.Glossaries); err != nil {
		log.Printf("GetStats glossary query: %v", err)
		return fiber.NewError(fiber.StatusInternalServerError, "DB error")
	}

	return c.JSON(stats)
}
