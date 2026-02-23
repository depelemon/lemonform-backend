package auth

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func App(app fiber.Router, db *gorm.DB) {
	r := app.Group("/auth")

	ctr := NewController(db)

	r.Post("/register", ctr.Register)
	r.Post("/login", ctr.Login)
}
