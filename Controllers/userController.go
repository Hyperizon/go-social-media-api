package Controllers

import (
	db "go-social-media-api/Config"
	models "go-social-media-api/Models"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
)

func Register(c *fiber.Ctx) error {
	var user models.Users
	if err := c.BodyParser(&user); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Invalid user data",
		})
	}

	var dbUser = models.Users{}
	db.DB.Where("username = ?", user.Username).First(&dbUser)

	if dbUser.Username == user.Username {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Username already exists",
		})
	}

	if user.Username == "" {
		return c.Status(http.StatusBadRequest).JSON(
			fiber.Map{
				"success": false,
				"message": "Username is required",
			})
	}

	if user.Password == "" {
		return c.Status(http.StatusBadRequest).JSON(
			fiber.Map{
				"success": false,
				"message": "Password is required",
			})
	}

	userData := models.Users{
		Username:  user.Username,
		Password:  user.Password,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	db.DB.Create(&userData)
	return c.Status(http.StatusOK).JSON(fiber.Map{"success": true, "message": "User created successfully", "data": userData})
}

func Login(c *fiber.Ctx) error {
	var user models.Users
	if err := c.BodyParser(&user); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Invalid user data",
		})
	}

	if user.Username == "" {
		return c.Status(http.StatusBadRequest).JSON(
			fiber.Map{
				"success": false,
				"message": "Username is required",
			})
	}

	if user.Password == "" {
		return c.Status(http.StatusBadRequest).JSON(
			fiber.Map{
				"success": false,
				"message": "Password is required",
				"error":   map[string]interface{}{},
			})
	}

	var dbUser = models.Users{}
	db.DB.Where("username = ?", user.Username).First(&dbUser)

	if user.Id == 0 {
		return c.Status(http.StatusNotFound).JSON(
			fiber.Map{
				"success": false,
				"message": "User not found",
			})
	}

	if dbUser.Password != user.Password {
		return c.Status(http.StatusBadRequest).JSON(
			fiber.Map{
				"success": false,
				"message": "Incorrect password",
			})
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId":    strconv.Itoa(int(user.Id)),
		"username":  user.Username,
		"ExpiresAt": time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return c.Status(http.StatusUnauthorized).JSON(
			fiber.Map{
				"success": false,
				"message": "Token Expired or  Invalid",
			})
	}

	cookie := fiber.Cookie{
		Name:     "token",
		Value:    tokenString,
		Expires:  time.Now().Add(time.Hour * 24),
		HTTPOnly: true,
		Secure:   true,
		SameSite: "Lax",
	}
	c.Cookie(&cookie)

	userData := make(map[string]interface{})
	userData["token"] = tokenString
	return c.Status(http.StatusOK).JSON(fiber.Map{"success": true, "message": "Logged in successfully", "data": userData})
}

func Logout(c *fiber.Ctx) error {
	cookie := fiber.Cookie{
		Name:     "token",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HTTPOnly: true,
		Secure:   true,
		SameSite: "Lax",
	}

	c.Cookie(&cookie)

	return c.Status(http.StatusOK).JSON(fiber.Map{"success": true, "message": "Logged out successfully"})
}

func GetUserPosts(c *fiber.Ctx) error {
	userId := c.Params("id")

	if userId == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"success": false, "message": "User id is required"})
	}

	posts := []models.Posts{}
	db.DB.Where("user_id = ?", userId).Find(&posts)

	return c.Status(http.StatusOK).JSON(fiber.Map{"success": true, "message": "User posts", "data": posts})
}
