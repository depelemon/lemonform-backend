package response

import (
	"strconv"

	"github.com/crlnravel/go-fiber-template/internal/models"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func itoa(id uint) string {
	return strconv.FormatUint(uint64(id), 10)
}

type Controller interface {
	Submit(c *fiber.Ctx) error
	List(c *fiber.Ctx) error
	GetPublic(c *fiber.Ctx) error
}

type controller struct {
	db *gorm.DB
}

func NewController(db *gorm.DB) Controller {
	return &controller{db: db}
}

// ── Request types ───────────────────────────────────────────────

type answerPayload struct {
	QuestionID uint   `json:"question_id"`
	Value      string `json:"value"`
}

type submitRequest struct {
	Answers []answerPayload `json:"answers"`
}

// ── Handlers ────────────────────────────────────────────────────

// Submit godoc
// @Summary      Submit a response
// @Description  Submit answers to a form. The form must be open (status = "open"). No authentication required.
// @Tags         Responses
// @Accept       json
// @Produce      json
// @Param        id    path  int            true  "Form ID"
// @Param        body  body  submitRequest  true  "Response payload"
// @Success      201  {object} map[string]interface{}
// @Failure      400  {object} map[string]interface{}
// @Failure      403  {object} map[string]interface{}
// @Failure      404  {object} map[string]interface{}
// @Router       /api/v1/forms/{id}/responses [post]
func (ctr *controller) Submit(c *fiber.Ctx) error {
	formID, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"ok":    false,
			"error": "invalid form id",
		})
	}

	// Find the form
	var form models.Form
	if err := ctr.db.First(&form, formID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"ok":    false,
			"error": "form does not exist",
		})
	}

	// Check that the form is open
	if form.Status != "open" {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"ok":    false,
			"error": "this form is closed and no longer accepts responses",
		})
	}

	var req submitRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"ok":    false,
			"error": "invalid request body",
		})
	}

	if len(req.Answers) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"ok":    false,
			"error": "answers are required",
		})
	}

	// Load this form's questions so we can validate
	var questions []models.Question
	ctr.db.Where("form_id = ?", formID).Find(&questions)

	// Build a set of valid question IDs for this form
	validQIDs := make(map[uint]bool, len(questions))
	for _, q := range questions {
		validQIDs[q.ID] = true
	}

	// Validate every submitted answer references a question in this form
	for _, a := range req.Answers {
		if !validQIDs[a.QuestionID] {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"ok":    false,
				"error": "question_id " + itoa(a.QuestionID) + " does not belong to this form",
			})
		}
	}

	// Check that every required question has an answer
	submittedQIDs := make(map[uint]bool, len(req.Answers))
	for _, a := range req.Answers {
		submittedQIDs[a.QuestionID] = true
	}
	for _, q := range questions {
		if q.Required && !submittedQIDs[q.ID] {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"ok":    false,
				"error": "missing answer for required question: " + q.Label,
			})
		}
	}

	// Build the response with answers inside a transaction
	resp := models.Response{
		FormID: uint(formID),
	}

	err = ctr.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&resp).Error; err != nil {
			return err
		}

		for _, a := range req.Answers {
			answer := models.Answer{
				ResponseID: resp.ID,
				QuestionID: a.QuestionID,
				Value:      a.Value,
			}
			if err := tx.Create(&answer).Error; err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"ok":    false,
			"error": "failed to submit response",
		})
	}

	// Reload with answers
	ctr.db.Preload("Answers").First(&resp, resp.ID)

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"ok":       true,
		"response": resp,
	})
}

// List godoc
// @Summary      List responses
// @Description  Get all responses for a form. Only the form owner can access this.
// @Tags         Responses
// @Produce      json
// @Security     BearerAuth
// @Param        id  path  int  true  "Form ID"
// @Success      200  {object} map[string]interface{}
// @Failure      404  {object} map[string]interface{}
// @Router       /api/v1/forms/{id}/responses [get]
func (ctr *controller) List(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)
	formID, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"ok":    false,
			"error": "invalid form id",
		})
	}

	// Verify form ownership
	var form models.Form
	if err := ctr.db.Where("id = ? AND owner_id = ?", formID, userID).First(&form).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"ok":    false,
			"error": "form does not exist or you do not have permission to access it",
		})
	}

	var responses []models.Response
	if err := ctr.db.Preload("Answers").Where("form_id = ?", formID).Order("created_at DESC").Find(&responses).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"ok":    false,
			"error": "failed to fetch responses",
		})
	}

	return c.JSON(fiber.Map{
		"ok":        true,
		"responses": responses,
	})
}

// GetPublic godoc
// @Summary      Get a form publicly
// @Description  Get a form with its questions (no auth required). Only returns open forms.
// @Tags         Responses
// @Produce      json
// @Param        id  path  int  true  "Form ID"
// @Success      200  {object} map[string]interface{}
// @Failure      404  {object} map[string]interface{}
// @Router       /api/v1/forms/{id}/public [get]
func (ctr *controller) GetPublic(c *fiber.Ctx) error {
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
	}).First(&form, formID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"ok":    false,
			"error": "form does not exist",
		})
	}

	if form.Status != "open" {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"ok":    false,
			"error": "this form is closed and no longer accepts responses",
		})
	}

	return c.JSON(fiber.Map{
		"ok":   true,
		"form": form,
	})
}
