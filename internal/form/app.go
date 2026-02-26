package form

import (
	"github.com/crlnravel/go-fiber-template/internal/auth"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func App(app fiber.Router, db *gorm.DB) {
	r := app.Group("/forms", auth.JWTMiddleware())

	ctr := NewController(db)

	r.Get("/", ctr.List)
	r.Get("/:id", ctr.Get)
	r.Post("/", ctr.Create)
	r.Put("/:id", ctr.Update)
	r.Delete("/:id", ctr.Delete)

	// Question CRUD (nested under forms)
	r.Post("/:id/questions", ctr.CreateQuestion)
	r.Put("/:id/questions/:qid", ctr.UpdateQuestion)
	r.Delete("/:id/questions/:qid", ctr.DeleteQuestion)
}
