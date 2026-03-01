package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/swagger"
	"gorm.io/gorm"

	_ "github.com/crlnravel/go-fiber-template/docs"

	"github.com/crlnravel/go-fiber-template/internal/auth"
	"github.com/crlnravel/go-fiber-template/internal/config"
	"github.com/crlnravel/go-fiber-template/internal/form"
	"github.com/crlnravel/go-fiber-template/internal/response"
	"github.com/crlnravel/go-fiber-template/internal/test"
)

type appConfig struct {
	prod bool
	port int
	db   *gorm.DB
}

func NewApp(cfg *appConfig) *fiber.App {
	app := fiber.New(fiber.Config{
		Prefork: cfg.prod,
	})

	app.Use(recover.New())
	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins:     config.GetEnv("CORS_ORIGINS", "http://localhost:3000"),
		AllowHeaders:     "*",
		AllowMethods:     "*",
		AllowCredentials: true,
	}))

	app.Get("/healthcheck", func(c *fiber.Ctx) error {
		return c.Status(200).SendString("hello from lemonform! 🍋")
	})

	app.Get("/swagger/*", swagger.HandlerDefault)

	api := app.Group("/api")

	v1 := api.Group("/v1")

	// Register route modules
	auth.App(v1, cfg.db)
	response.App(v1, cfg.db) // must be before form to avoid JWT on public submit
	form.App(v1, cfg.db)
	test.App(v1, cfg.db)

	return app
}
