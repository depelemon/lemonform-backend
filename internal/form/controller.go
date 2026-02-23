package form

import (
	"github.com/crlnravel/go-fiber-template/internal/models"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type Controller interface {
	List(c *fiber.Ctx) error
	Get(c *fiber.Ctx) error
	Create(c *fiber.Ctx) error
	Update(c *fiber.Ctx) error
	Delete(c *fiber.Ctx) error
}

type controller struct {
	db *gorm.DB
}

func NewController(db *gorm.DB) Controller {
	return &controller{db: db}
}

type createFormRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

type updateFormRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

// List godoc
// @Summary      List forms
// @Description  Get all forms owned by the authenticated user.
// @Tags         Forms
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object} map[string]interface{}
// @Router       /api/v1/forms [get]
func (ctr *controller) List(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)

	var forms []models.Form
	if err := ctr.db.Where("owner_id = ?", userID).Order("created_at DESC").Find(&forms).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"ok":    false,
			"error": "failed to fetch forms",
		})
	}

	return c.JSON(fiber.Map{
		"ok":    true,
		"forms": forms,
	})
}

// Get godoc
// @Summary      Get form detail
// @Description  Get a single form with its questions.
// @Tags         Forms
// @Produce      json
// @Security     BearerAuth
// @Param        id   path  int  true  "Form ID"
// @Success      200  {object} map[string]interface{}
// @Router       /api/v1/forms/{id} [get]
func (ctr *controller) Get(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)
	formID, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"ok":    false,
			"error": "invalid form id",
		})
	}

	var form models.Form
	if err := ctr.db.Preload("Questions", func(db *gorm.DB) *gorm.DB {
		return db.Order("\"order\" ASC")
	}).Where("id = ? AND owner_id = ?", formID, userID).First(&form).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"ok":    false,
			"error": "form not found",
		})
	}

	return c.JSON(fiber.Map{
		"ok":   true,
		"form": form,
	})
}

// Create godoc
// @Summary      Create a form
// @Description  Create a new form for the authenticated user.
// @Tags         Forms
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        body  body  createFormRequest  true  "Form payload"
// @Success      201  {object} map[string]interface{}
// @Router       /api/v1/forms [post]
func (ctr *controller) Create(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)

	var req createFormRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"ok":    false,
			"error": "invalid request body",
		})
	}

	if req.Title == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"ok":    false,
			"error": "title is required",
		})
	}

	form := models.Form{
		Title:       req.Title,
		Description: req.Description,
		OwnerID:     userID,
	}

	if err := ctr.db.Create(&form).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"ok":    false,
			"error": "failed to create form",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"ok":   true,
		"form": form,
	})
}

// Update godoc
// @Summary      Update a form
// @Description  Update an existing form's title and description.
// @Tags         Forms
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id    path  int               true  "Form ID"
// @Param        body  body  updateFormRequest  true  "Form payload"
// @Success      200  {object} map[string]interface{}
// @Router       /api/v1/forms/{id} [put]
func (ctr *controller) Update(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)
	formID, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"ok":    false,
			"error": "invalid form id",
		})
	}

	var form models.Form
	if err := ctr.db.Where("id = ? AND owner_id = ?", formID, userID).First(&form).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"ok":    false,
			"error": "form not found",
		})
	}

	var req updateFormRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"ok":    false,
			"error": "invalid request body",
		})
	}

	if req.Title != "" {
		form.Title = req.Title
	}
	if req.Description != "" {
		form.Description = req.Description
	}

	if err := ctr.db.Save(&form).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"ok":    false,
			"error": "failed to update form",
		})
	}

	return c.JSON(fiber.Map{
		"ok":   true,
		"form": form,
	})
}

// Delete godoc
// @Summary      Delete a form
// @Description  Delete a form and its questions.
// @Tags         Forms
// @Produce      json
// @Security     BearerAuth
// @Param        id  path  int  true  "Form ID"
// @Success      200  {object} map[string]interface{}
// @Router       /api/v1/forms/{id} [delete]
func (ctr *controller) Delete(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)
	formID, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"ok":    false,
			"error": "invalid form id",
		})
	}

	var form models.Form
	if err := ctr.db.Where("id = ? AND owner_id = ?", formID, userID).First(&form).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"ok":    false,
			"error": "form not found",
		})
	}

	// Delete associated questions first
	ctr.db.Where("form_id = ?", formID).Delete(&models.Question{})

	if err := ctr.db.Delete(&form).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"ok":    false,
			"error": "failed to delete form",
		})
	}

	return c.JSON(fiber.Map{
		"ok":      true,
		"message": "form deleted successfully",
	})
}
