package handlers

import (
	"log"
	"strconv"
	"strings"

	"github.com/bxra2/taqnun/db"
	"github.com/bxra2/taqnun/internal/normalize"
	"github.com/bxra2/taqnun/models"
	"github.com/gofiber/fiber/v3"
)

// termSelect mirrors the column list used by GetTerms so that the scan order
// and COALESCE behavior stay consistent across endpoints.
const termSelect = `
	t.id, COALESCE(t.english,''), COALESCE(t.arabic,''),
	COALESCE(t.french,''), COALESCE(t.german,''),
	COALESCE(t.term_url,''), COALESCE(t.description,''),
	COALESCE(g.glossaryEn,''), COALESCE(g.glossaryAr,''), COALESCE(g.glossaryUrl,''),
	COALESCE(g.publisher_id, 0),
	COALESCE(p.publisherAr,''), COALESCE(p.publisherEn,''), COALESCE(p.publisherUrl,'')`

const termFrom = `
	FROM term t
	JOIN glossary g ON g.id = t.glossary_id
	LEFT JOIN publisher p ON p.id = g.publisher_id`

// escapeLike escapes the SQL LIKE wildcards and escape character so a
// user-supplied query is treated as a literal substring.
func escapeLike(s string) string {
	s = strings.ReplaceAll(s, `\`, `\\`)
	s = strings.ReplaceAll(s, `%`, `\%`)
	s = strings.ReplaceAll(s, `_`, `\_`)
	return s
}

// SearchTerms performs a server-side, paginated search over the term table.
func SearchTerms(c fiber.Ctx) error {
	c.Set("Cache-Control", "no-store")

	q := strings.TrimSpace(c.Query("q"))
	exact := c.Query("exact") == "true" || c.Query("exact") == "1"
	desc := c.Query("desc") == "true" || c.Query("desc") == "1"

	limit := 10
	if raw := c.Query("limit"); raw != "" {
		if n, err := strconv.Atoi(raw); err == nil {
			limit = n
		}
	}
	if limit < 1 {
		limit = 1
	}
	if limit > 100 {
		limit = 100
	}

	offset := 0
	if raw := c.Query("offset"); raw != "" {
		if n, err := strconv.Atoi(raw); err == nil {
			offset = n
		}
	}
	if offset < 0 {
		offset = 0
	}

	if q == "" {
		return c.JSON(models.SearchResult{
			Items:  []models.Term{},
			Total:  0,
			Limit:  limit,
			Offset: offset,
		})
	}

	latin := normalize.Latin(q)
	arabic := normalize.Arabic(q)
	text := normalize.Text(q)

	containsEn := "%" + escapeLike(latin) + "%"
	containsAr := "%" + escapeLike(arabic) + "%"
	containsDesc := "%" + escapeLike(text) + "%"
	prefixEn := escapeLike(latin) + "%"
	prefixAr := escapeLike(arabic) + "%"

	var where string
	var args []interface{}
	if exact {
		where = `WHERE t.english_norm = ? OR t.arabic_norm = ?`
		args = []interface{}{latin, arabic}
	} else {
		where = `WHERE t.english_norm LIKE ? ESCAPE '\' OR t.arabic_norm LIKE ? ESCAPE '\'`
		args = []interface{}{containsEn, containsAr}
		if desc {
			where += ` OR t.description_norm LIKE ? ESCAPE '\'`
			args = append(args, containsDesc)
		}
	}

	var order string
	var orderArgs []interface{}
	if exact {
		order = `ORDER BY t.id`
	} else {
		order = `ORDER BY
			CASE WHEN t.english_norm = ? OR t.arabic_norm = ? THEN 0 ELSE 1 END,
			CASE WHEN t.english_norm LIKE ? ESCAPE '\' OR t.arabic_norm LIKE ? ESCAPE '\' THEN 0 ELSE 1 END,
			t.id`
		orderArgs = []interface{}{latin, arabic, prefixEn, prefixAr}
	}

	var total int
	countQuery := `SELECT COUNT(*) ` + termFrom + ` ` + where
	if err := db.Conn.QueryRow(countQuery, args...).Scan(&total); err != nil {
		log.Printf("SearchTerms count: %v", err)
		return fiber.NewError(fiber.StatusInternalServerError, "DB error")
	}

	pageArgs := make([]interface{}, 0, len(args)+len(orderArgs)+2)
	pageArgs = append(pageArgs, args...)
	pageArgs = append(pageArgs, orderArgs...)
	pageArgs = append(pageArgs, limit, offset)

	pageQuery := `SELECT ` + termSelect + ` ` + termFrom + ` ` + where + ` ` + order + ` LIMIT ? OFFSET ?`
	rows, err := db.Conn.Query(pageQuery, pageArgs...)
	if err != nil {
		log.Printf("SearchTerms query: %v", err)
		return fiber.NewError(fiber.StatusInternalServerError, "DB error")
	}
	defer rows.Close()

	terms := make([]models.Term, 0)
	for rows.Next() {
		var t models.Term
		if err := rows.Scan(&t.ID, &t.English, &t.Arabic, &t.French, &t.German,
			&t.TURL, &t.Description,
			&t.GlossaryEn, &t.GlossaryAr, &t.GlossaryURL, &t.PublisherID, &t.PublisherAr, &t.PublisherEn, &t.PublisherURL); err != nil {
			log.Printf("SearchTerms scan: %v", err)
			return fiber.NewError(fiber.StatusInternalServerError, "Decode error")
		}
		terms = append(terms, t)
	}
	if err := rows.Err(); err != nil {
		log.Printf("SearchTerms rows: %v", err)
		return fiber.NewError(fiber.StatusInternalServerError, "DB error")
	}

	return c.JSON(models.SearchResult{
		Items:  terms,
		Total:  total,
		Limit:  limit,
		Offset: offset,
	})
}
