package main

import (
	"fmt"
	"log"

	"github.com/SMSLYCLOUD/smsly-code/smsly-code-api/internal/config"
	"github.com/SMSLYCLOUD/smsly-code/smsly-code-api/internal/database"
	"github.com/SMSLYCLOUD/smsly-code/smsly-code-api/internal/handlers"
	"github.com/SMSLYCLOUD/smsly-code/smsly-code-api/internal/middleware"
	"github.com/SMSLYCLOUD/smsly-code/smsly-code-api/internal/models"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	cfg := config.LoadConfig()

	db, err := database.Connect(cfg)
	if err != nil {
		log.Printf("Warning: Failed to connect to database: %v", err)
	} else {
		sqlDB, _ := db.DB()
		defer sqlDB.Close()

		// Auto Migrate
		log.Println("Migrating database...")
		err = db.AutoMigrate(
			&models.User{},
			&models.Repository{},
			&models.Issue{},
			&models.Comment{},
		)
		if err != nil {
			log.Printf("Migration failed: %v", err)
		}
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

	if db != nil {
		authHandler := handlers.NewAuthHandler(db, &cfg)

		api := app.Group("/api")

		// Public routes
		auth := api.Group("/auth")
		auth.Post("/register", authHandler.Register)
		auth.Post("/login", authHandler.Login)

		// Protected routes
		protected := api.Group("/")
		protected.Use(middleware.Auth(db, &cfg))
		protected.Get("/me", authHandler.GetMe)

		// Repository routes
		repoHandler := handlers.NewRepoHandler(db, &cfg)
		protected.Post("/repos", repoHandler.Create)
		protected.Get("/repos", repoHandler.List)
		protected.Get("/repos/:name/commits", repoHandler.GetCommits)
		protected.Get("/repos/:name/tree/:ref/*", repoHandler.GetTree)
		// Catch-all for root of tree with ref
		protected.Get("/repos/:name/tree/:ref", repoHandler.GetTree)

		// Issue routes
		issueHandler := handlers.NewIssueHandler(db, &cfg)
		protected.Post("/repos/:name/issues", issueHandler.Create)
		protected.Get("/repos/:name/issues", issueHandler.List)
		protected.Get("/repos/:name/issues/:id", issueHandler.Get)
		protected.Patch("/repos/:name/issues/:id", issueHandler.Update)

		// Comment routes
		protected.Post("/repos/:name/issues/:id/comments", issueHandler.CreateComment)
		protected.Get("/repos/:name/issues/:id/comments", issueHandler.ListComments)
	}

	addr := fmt.Sprintf(":%s", cfg.Port)
	log.Printf("Server listening on %s", addr)
	log.Fatal(app.Listen(addr))
}
