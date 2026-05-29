package main

import (
	"bufio"
	"log"
	"os"
	"strings"

	"github.com/bxra2/taqnun/db"
	"github.com/bxra2/taqnun/handlers"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/static"
)

const staticRoot = "./frontend/dist"

func main() {
	loadDotenv(".env")

	uri := os.Getenv("MONGODB_URI")
	if uri == "" {
		log.Fatal("MONGODB_URI is not set")
	}
	db.Connect(uri)

	app := fiber.New()

	api := app.Group("/api")
	api.Get("/publishers", handlers.GetPublishers)
	api.Get("/terms", handlers.GetTerms)
	api.Get("/localizations", handlers.GetLocalizations)

	// Static files + SPA fallback to index.html for client-side routes.
	app.Get("/*", static.New(staticRoot, static.Config{
		IndexNames: []string{"index.html"},
		NotFoundHandler: func(c fiber.Ctx) error {
			c.Status(fiber.StatusOK)
			return c.SendFile(staticRoot + "/index.html")
		},
	}))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Fatal(app.Listen(":" + port))
}

// loadDotenv reads simple KEY=VALUE pairs from a .env file into the process env.
// Existing environment values take precedence. Missing file is not an error.
func loadDotenv(path string) {
	f, err := os.Open(path)
	if err != nil {
		return
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		key, value, ok := strings.Cut(line, "=")
		if !ok {
			continue
		}
		key = strings.TrimSpace(key)
		value = strings.TrimSpace(value)
		if len(value) >= 2 && (value[0] == '"' && value[len(value)-1] == '"' || value[0] == '\'' && value[len(value)-1] == '\'') {
			value = value[1 : len(value)-1]
		}
		if _, exists := os.LookupEnv(key); !exists {
			os.Setenv(key, value)
		}
	}
}
