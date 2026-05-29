package handlers

import (
	"github.com/bxra2/taqnun/localizations"
	"github.com/gofiber/fiber/v3"
)

func GetLocalizations(c fiber.Ctx) error {
	entries, err := localizations.FetchAll(c.Context())
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "failed to load localizations")
	}
	c.Set("Cache-Control", "public, max-age=21600, stale-while-revalidate=86400")
	return c.JSON(entries)
}
