package apis

import (
	"time"

	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v3"
	"github.com/golang-jwt/jwt/v4"
	"github.com/miceremwirigi/go-fiber-jwt/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Handler struct {
	DB *gorm.DB
}

// --------------------------Signup handler--------------------------//
type SignupRequest struct {
	Name     string
	Email    string
	Password string
}

func (h Handler) Signup(c *fiber.Ctx) error {
	req := new(SignupRequest)
	err := c.BodyParser(&req)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Failed to parse inputed data")
	}
	if req.Name == "" || req.Email == "" || req.Password == "" {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid signup credentials")
	}

	//Save this info in database
	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return fiber.NewError(fiber.StatusExpectationFailed, "Error encrypting password")
	}

	user := &models.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: string(hash),
	}

	err = h.DB.Create(user).Error
	if err != nil {
		return fiber.NewError(fiber.StatusExpectationFailed, "Storing user credentials in database failed. Error:"+err.Error())
	}

	// Create new jwt token for user
	token, exp, err := CreateJWTToken(*user)
	if err != nil {
		return c.JSON(fiber.Map{
			"success": false,
			"message": "failed to create token",
			"user":    user,
		})
	}

	return c.JSON(fiber.Map{
		"token":   token,
		"expires": exp,
		"message": "Token created. Signed up",
		"user":    user,
	})
}

// --------------------------login handler--------------------------//
type LoginRequest struct {
	Email    string
	Password string
}

func (h Handler) Login(c *fiber.Ctx) error {
	req := new(LoginRequest)

	if err := c.BodyParser(&req); err != nil {
		return err
	}

	if req.Email == "" || req.Password == "" {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid/empty login credentials")
	}

	// Verify user in database
	user := new(models.User)
	err := h.DB.First(user, "email = ?", req.Email).Error

	if err != nil {
		return c.Status(fiber.StatusExpectationFailed).JSON(fiber.Map{
			"success": false,
			"message": "error  accessing user/ Invalid login credentials",
			"data":    err,
		})
	}

	// Check if login credetials are correct
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return c.JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	// Create new jwt token for user
	token, exp, err := CreateJWTToken(*user)
	if err != nil {
		return c.JSON(fiber.Map{
			"success": false,
			"message": "failed to create token",
			"user":    user,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"token":   token,
		"expires": exp,
		"message": "Login succesful",
		"user":    user,
	})
}

// --------------------------logout handler--------------------------//

func (h Handler) Logout(c *fiber.Ctx) error {

	return nil
}

// --------------------------Private handler--------------------------//

func (h Handler) Private(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"success": true,
		"path":    "private",
	})
}

// --------------------------Public handler--------------------------//

func (h Handler) Public(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"success": true,
		"path":    "public",
	})
}

// --------------------------Setup routes--------------------------//

func (h Handler) SetupRoutes(app *fiber.App) {
	api := app.Group("/")
	api.Post("/signup", h.Signup)
	api.Post("/login", h.Login)
	api.Get("/public", h.Public)

	private := app.Group("/private")
	private.Use(jwtware.New(jwtware.Config{
		SigningKey: []byte("secret"),
	}))
	private.Get("/", h.Private)
	private.Get("/logout", h.Logout)

}

// --------------------------Generate JWT Token----------------------//

func CreateJWTToken(user models.User) (string, int64, error) {
	exp := time.Now().Add(time.Minute * 5).Unix()
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["user_id"] = user.ID
	claims["exp"] = exp
	t, err := token.SignedString([]byte("secret")) // sign token to make it complete
	if err != nil {
		return "", 0, err
	}
	return t, exp, nil
}
