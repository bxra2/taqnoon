package handlers

import (
	"log"
	"sync"
	"time"

	"github.com/bxra2/taqnun/db"
	"github.com/bxra2/taqnun/models"
	"github.com/gofiber/fiber/v3"
)

const suggestTTL = 5 * time.Minute

var (
	suggestMu     sync.Mutex
	suggestCached []models.Suggestion
	suggestExpiry time.Time
)

// GetSuggestions returns a small set of single-word terms chosen at random.
// Results are cached in memory for 5 minutes (populated lazily on first use).
func GetSuggestions(c fiber.Ctx) error {
	c.Set("Cache-Control", "private, max-age=300")

	suggestMu.Lock()
	defer suggestMu.Unlock()

	if suggestCached != nil && time.Now().Before(suggestExpiry) {
		return c.JSON(suggestCached)
	}

	rows, err := db.Conn.Query(`
		SELECT id, COALESCE(english,''), COALESCE(arabic,'')
		FROM term
		WHERE (arabic_norm <> '' AND arabic_norm NOT LIKE '% %')
		   OR (english_norm <> '' AND english_norm NOT LIKE '% %')
		ORDER BY RANDOM()
		LIMIT 10`)
	if err != nil {
		log.Printf("GetSuggestions query: %v", err)
		return fiber.NewError(fiber.StatusInternalServerError, "DB error")
	}
	defer rows.Close()

	suggestions := make([]models.Suggestion, 0, 10)
	for rows.Next() {
		var s models.Suggestion
		if err := rows.Scan(&s.ID, &s.English, &s.Arabic); err != nil {
			log.Printf("GetSuggestions scan: %v", err)
			return fiber.NewError(fiber.StatusInternalServerError, "Decode error")
		}
		suggestions = append(suggestions, s)
	}
	if err := rows.Err(); err != nil {
		log.Printf("GetSuggestions rows: %v", err)
		return fiber.NewError(fiber.StatusInternalServerError, "DB error")
	}

	suggestCached = suggestions
	suggestExpiry = time.Now().Add(suggestTTL)

	return c.JSON(suggestions)
}
