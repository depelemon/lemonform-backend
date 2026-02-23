package auth

import (
	"github.com/crlnravel/go-fiber-template/internal/models"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Controller interface {
	Register(c *fiber.Ctx) error
	Login(c *fiber.Ctx) error
}

type controller struct {
	db *gorm.DB
}

func NewController(db *gorm.DB) Controller {
	return &controller{db: db}
}

type registerRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type authResponse struct {
	OK    bool   `json:"ok"`
	Token string `json:"token"`
}

// Register godoc
// @Summary      Register a new user
// @Description  Create a new account with email and password.
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        body  body  registerRequest  true  "Register payload"
// @Success      201  {object} authResponse  "Created"
// @Failure      400  {object} common.ErrorResponse  "Bad Request"
// @Failure      409  {object} common.ErrorResponse  "Conflict (email taken)"
// @Router       /api/v1/auth/register [post]
func (ctr *controller) Register(c *fiber.Ctx) error {
	var req registerRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"ok":    false,
			"error": "invalid request body",
		})
	}

	if req.Email == "" || req.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"ok":    false,
			"error": "email and password are required",
		})
	}

	// Check if email already exists
	var existing models.User
	if err := ctr.db.Where("email = ?", req.Email).First(&existing).Error; err == nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"ok":    false,
			"error": "email already registered",
		})
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"ok":    false,
			"error": "failed to hash password",
		})
	}

	user := models.User{
		Email:    req.Email,
		Password: string(hashedPassword),
	}

	if err := ctr.db.Create(&user).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"ok":    false,
			"error": "failed to create user",
		})
	}

	// Generate JWT
	token, err := GenerateToken(user.ID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"ok":    false,
			"error": "failed to generate token",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(authResponse{
		OK:    true,
		Token: token,
	})
}

// Login godoc
// @Summary      Login
// @Description  Authenticate with email and password, receive a JWT.
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        body  body  loginRequest  true  "Login payload"
// @Success      200  {object} authResponse  "OK"
// @Failure      400  {object} common.ErrorResponse  "Bad Request"
// @Failure      401  {object} common.ErrorResponse  "Unauthorized"
// @Router       /api/v1/auth/login [post]
func (ctr *controller) Login(c *fiber.Ctx) error {
	var req loginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"ok":    false,
			"error": "invalid request body",
		})
	}

	if req.Email == "" || req.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"ok":    false,
			"error": "email and password are required",
		})
	}

	// Find user
	var user models.User
	if err := ctr.db.Where("email = ?", req.Email).First(&user).Error; err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"ok":    false,
			"error": "invalid email or password",
		})
	}

	// Compare password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"ok":    false,
			"error": "invalid email or password",
		})
	}

	// Generate JWT
	token, err := GenerateToken(user.ID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"ok":    false,
			"error": "failed to generate token",
		})
	}

	return c.JSON(authResponse{
		OK:    true,
		Token: token,
	})
}
