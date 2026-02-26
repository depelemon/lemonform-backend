package response

import (
	"github.com/crlnravel/go-fiber-template/internal/auth"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func App(app fiber.Router, db *gorm.DB) {
	ctr := NewController(db)

	// Public — anyone can submit a response to an open form
	app.Post("/forms/:id/responses", ctr.Submit)

	// Protected — only the form owner can view responses
	app.Get("/forms/:id/responses", auth.JWTMiddleware(), ctr.List)
}
