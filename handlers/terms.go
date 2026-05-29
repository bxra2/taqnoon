package handlers

import (
	"context"
	"time"

	"github.com/bxra2/taqnun/db"
	"github.com/bxra2/taqnun/models"
	"github.com/gofiber/fiber/v3"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func GetTerms(c fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	cursor, err := db.GetCollection("taqnoon", "terms").Find(ctx, bson.M{})
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "DB error")
	}
	defer cursor.Close(ctx)

	terms := make([]models.Term, 0)
	if err := cursor.All(ctx, &terms); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Decode error")
	}

	return c.JSON(terms)
}
