package main

import (
	"fmt"
	"log"

	"github.com/SMSLYCLOUD/smsly-code/smsly-code-api/internal/config"
	"github.com/SMSLYCLOUD/smsly-code/smsly-code-api/internal/database"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	cfg := config.LoadConfig()

	// Try to connect to DB, but don't crash if it fails (for now)
	db, err := database.Connect(cfg)
	if err != nil {
		log.Printf("Warning: Failed to connect to database: %v", err)
	} else {
		sqlDB, _ := db.DB()
		defer sqlDB.Close()
	}

	app := fiber.New()

	app.Use(logger.New())
	app.Use(cors.New())

	app.Get("/health", func(c *fiber.Ctx) error {
		status := "ok"
		if db == nil {
			status = "degraded (no db)"
		}
		return c.JSON(fiber.Map{
			"status": status,
			"service": "smsly-code-api",
		})
	})

	addr := fmt.Sprintf(":%s", cfg.Port)
	log.Printf("Server listening on %s", addr)
	log.Fatal(app.Listen(addr))
}
