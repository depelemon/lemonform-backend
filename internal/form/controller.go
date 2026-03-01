package form

import (
	"strconv"
	"time"

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
	CreateQuestion(c *fiber.Ctx) error
	UpdateQuestion(c *fiber.Ctx) error
	DeleteQuestion(c *fiber.Ctx) error
}

type controller struct {
	db *gorm.DB
}

func NewController(db *gorm.DB) Controller {
	return &controller{db: db}
}

// ── Request types ───────────────────────────────────────────────

type createFormRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

type updateFormRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Status      string `json:"status"` // "open" or "closed"
}

type createQuestionRequest struct {
	Label    string `json:"label"`
	Type     string `json:"type"`     // short_answer, radio, checkbox, dropdown
	Options  string `json:"options"`  // JSON-encoded options (for non-short_answer types)
	Required *bool  `json:"required"` // default true
	Order    int    `json:"order"`
}

type updateQuestionRequest struct {
	Label    string `json:"label"`
	Type     string `json:"type"`
	Options  string `json:"options"`
	Required *bool  `json:"required"`
	Order    int    `json:"order"`
}

// ── Form CRUD ───────────────────────────────────────────────────

// List godoc
// @Summary      List forms
// @Description  Get all forms owned by the authenticated user. Supports search, status filter, date range, sorting, and pagination.
// @Tags         Forms
// @Produce      json
// @Security     BearerAuth
// @Param        search          query  string  false  "Search by title (case-insensitive)"
// @Param        status          query  string  false  "Filter by status: open or closed"
// @Param        created_after   query  string  false  "Filter forms created after this date (RFC3339, e.g. 2025-01-01T00:00:00Z)"
// @Param        created_before  query  string  false  "Filter forms created before this date (RFC3339)"
// @Param        sort_by         query  string  false  "Field to sort by: created_at (default) or title"
// @Param        sort            query  string  false  "Sort direction: asc or desc (default desc)"
// @Param        page            query  int     false  "Page number (default 1)"
// @Param        limit           query  int     false  "Items per page (default 20, max 100)"
// @Success      200  {object} map[string]interface{}
// @Router       /api/v1/forms [get]
func (ctr *controller) List(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)

	query := ctr.db.Where("owner_id = ?", userID)

	// Search by title
	if search := c.Query("search"); search != "" {
		query = query.Where("title ILIKE ?", "%"+search+"%")
	}

	// Filter by status
	if status := c.Query("status"); status == "open" || status == "closed" {
		query = query.Where("status = ?", status)
	}

	// Filter by date range
	if after := c.Query("created_after"); after != "" {
		if t, err := time.Parse(time.RFC3339, after); err == nil {
			query = query.Where("created_at >= ?", t)
		}
	}
	if before := c.Query("created_before"); before != "" {
		if t, err := time.Parse(time.RFC3339, before); err == nil {
			query = query.Where("created_at <= ?", t)
		}
	}

	// Sort by field + direction
	sortBy := c.Query("sort_by", "created_at")
	if sortBy != "title" {
		sortBy = "created_at" // default & safeguard
	}
	sortDir := "DESC"
	if c.Query("sort") == "asc" {
		sortDir = "ASC"
	}
	query = query.Order(sortBy + " " + sortDir)

	// Count total (for pagination metadata)
	var total int64
	countQuery := *query
	countQuery.Model(&models.Form{}).Count(&total)

	// Pagination
	page, _ := strconv.Atoi(c.Query("page", "1"))
	if page < 1 {
		page = 1
	}
	limit, _ := strconv.Atoi(c.Query("limit", "20"))
	if limit < 1 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}
	offset := (page - 1) * limit

	var forms []models.Form
	if err := query.Offset(offset).Limit(limit).Find(&forms).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"ok":    false,
			"error": "failed to fetch forms",
		})
	}

	return c.JSON(fiber.Map{
		"ok":    true,
		"forms": forms,
		"meta": fiber.Map{
			"page":  page,
			"limit": limit,
			"total": total,
		},
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
	if err := ctr.db.Preload("Owner").Preload("Questions", func(db *gorm.DB) *gorm.DB {
		return db.Order("\"order\" ASC")
	}).Where("id = ? AND owner_id = ?", formID, userID).First(&form).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"ok":    false,
			"error": "form does not exist or you do not have permission to access it",
		})
	}

	// Count responses so frontend knows if questions are locked for editing
	var responseCount int64
	ctr.db.Model(&models.Response{}).Where("form_id = ?", formID).Count(&responseCount)

	return c.JSON(fiber.Map{
		"ok":             true,
		"form":           form,
		"owner_email":    form.Owner.Email,
		"response_count": responseCount,
	})
}

// Create godoc
// @Summary      Create a form
// @Description  Create a new form (defaults to status "open").
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
		Status:      "open",
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
// @Description  Update a form's title, description, or status (open/closed).
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
			"error": "form does not exist or you do not have permission to access it",
		})
	}

	// Check if form has responses — if so, only status may be changed
	var responseCount int64
	ctr.db.Model(&models.Response{}).Where("form_id = ?", formID).Count(&responseCount)

	var req updateFormRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"ok":    false,
			"error": "invalid request body",
		})
	}

	if responseCount > 0 {
		// Only allow status changes when form has responses
		if req.Title != "" || req.Description != "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"ok":    false,
				"error": "cannot edit title or description of a form that already has responses; you may only change the status",
			})
		}
	} else {
		if req.Title != "" {
			form.Title = req.Title
		}
		if req.Description != "" {
			form.Description = req.Description
		}
	}
	if req.Status == "open" || req.Status == "closed" {
		form.Status = req.Status
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
// @Description  Delete a form together with its questions and responses.
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
			"error": "form does not exist or you do not have permission to access it",
		})
	}

	// Delete associated answers (through responses), then responses, then questions
	var responseIDs []uint
	ctr.db.Model(&models.Response{}).Where("form_id = ?", formID).Pluck("id", &responseIDs)
	if len(responseIDs) > 0 {
		ctr.db.Where("response_id IN ?", responseIDs).Delete(&models.Answer{})
	}
	ctr.db.Where("form_id = ?", formID).Delete(&models.Response{})
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

// ── Question CRUD ───────────────────────────────────────────────

// CreateQuestion godoc
// @Summary      Add a question to a form
// @Description  Create a new question. Type must be one of: short_answer, multiple_choice, checkbox, dropdown.
// @Tags         Questions
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id    path  int                    true  "Form ID"
// @Param        body  body  createQuestionRequest   true  "Question payload"
// @Success      201  {object} map[string]interface{}
// @Router       /api/v1/forms/{id}/questions [post]
func (ctr *controller) CreateQuestion(c *fiber.Ctx) error {
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

	// Block question creation if form already has responses
	var respCount int64
	ctr.db.Model(&models.Response{}).Where("form_id = ?", formID).Count(&respCount)
	if respCount > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"ok":    false,
			"error": "cannot add questions to a form that already has responses",
		})
	}

	var req createQuestionRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"ok":    false,
			"error": "invalid request body",
		})
	}

	if req.Label == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"ok":    false,
			"error": "label is required",
		})
	}

	if req.Type == "" {
		req.Type = "short_answer"
	}
	if !models.ValidQuestionTypes[req.Type] {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"ok":    false,
			"error": "type must be one of: short_answer, radio, checkbox, dropdown",
		})
	}

	required := true
	if req.Required != nil {
		required = *req.Required
	}

	question := models.Question{
		FormID:   uint(formID),
		Label:    req.Label,
		Type:     req.Type,
		Options:  req.Options,
		Required: required,
		Order:    req.Order,
	}

	if err := ctr.db.Create(&question).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"ok":    false,
			"error": "failed to create question",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"ok":       true,
		"question": question,
	})
}

// UpdateQuestion godoc
// @Summary      Update a question
// @Description  Update a question's label, type, options, or order.
// @Tags         Questions
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id    path  int                    true  "Form ID"
// @Param        qid   path  int                    true  "Question ID"
// @Param        body  body  updateQuestionRequest   true  "Question payload"
// @Success      200  {object} map[string]interface{}
// @Router       /api/v1/forms/{id}/questions/{qid} [put]
func (ctr *controller) UpdateQuestion(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)
	formID, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"ok":    false,
			"error": "invalid form id",
		})
	}

	qid, err := strconv.Atoi(c.Params("qid"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"ok":    false,
			"error": "invalid question id",
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

	// Block question updates if form already has responses
	var respCount int64
	ctr.db.Model(&models.Response{}).Where("form_id = ?", formID).Count(&respCount)
	if respCount > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"ok":    false,
			"error": "cannot edit questions on a form that already has responses",
		})
	}

	var question models.Question
	if err := ctr.db.Where("id = ? AND form_id = ?", qid, formID).First(&question).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"ok":    false,
			"error": "question not found",
		})
	}

	var req updateQuestionRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"ok":    false,
			"error": "invalid request body",
		})
	}

	if req.Label != "" {
		question.Label = req.Label
	}
	if req.Type != "" {
		if !models.ValidQuestionTypes[req.Type] {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"ok":    false,
				"error": "type must be one of: short_answer, radio, checkbox, dropdown",
			})
		}
		question.Type = req.Type
	}
	if req.Options != "" {
		question.Options = req.Options
	}
	if req.Required != nil {
		question.Required = *req.Required
	}
	if req.Order != 0 {
		question.Order = req.Order
	}

	if err := ctr.db.Save(&question).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"ok":    false,
			"error": "failed to update question",
		})
	}

	return c.JSON(fiber.Map{
		"ok":       true,
		"question": question,
	})
}

// DeleteQuestion godoc
// @Summary      Delete a question
// @Description  Remove a question from a form.
// @Tags         Questions
// @Produce      json
// @Security     BearerAuth
// @Param        id   path  int  true  "Form ID"
// @Param        qid  path  int  true  "Question ID"
// @Success      200  {object} map[string]interface{}
// @Router       /api/v1/forms/{id}/questions/{qid} [delete]
func (ctr *controller) DeleteQuestion(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)
	formID, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"ok":    false,
			"error": "invalid form id",
		})
	}

	qid, err := strconv.Atoi(c.Params("qid"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"ok":    false,
			"error": "invalid question id",
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

	// Block question deletion if form already has responses
	var respCount int64
	ctr.db.Model(&models.Response{}).Where("form_id = ?", formID).Count(&respCount)
	if respCount > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"ok":    false,
			"error": "cannot delete questions from a form that already has responses",
		})
	}

	var question models.Question
	if err := ctr.db.Where("id = ? AND form_id = ?", qid, formID).First(&question).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"ok":    false,
			"error": "question not found",
		})
	}

	if err := ctr.db.Delete(&question).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"ok":    false,
			"error": "failed to delete question",
		})
	}

	return c.JSON(fiber.Map{
		"ok":      true,
		"message": "question deleted successfully",
	})
}
