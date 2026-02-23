package test

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func App(app fiber.Router, db *gorm.DB) {
	r := app.Group("/test")

	ctr := NewController()

	// Route listing

	r.Get("/", ctr.GetTestMsg)
}
