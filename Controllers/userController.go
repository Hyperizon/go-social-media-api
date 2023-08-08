package Controllers

import (
	"fmt"
	db "go-social-media-api/Config"
	Middleware "go-social-media-api/Middlewares"
	models "go-social-media-api/Models"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
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

	validate := validator.New()
	validate.RegisterValidation("usernameValid", Middleware.UsernameValid)
	if err := validate.Struct(user); err != nil {
		var errMsgs []string
		for _, err := range err.(validator.ValidationErrors) {
			fieldName := strings.ToLower(err.Field())
			errMsg := fmt.Sprintf("Field validation for '%s' failed on the '%s' tag", fieldName, err.Tag())
			errMsgs = append(errMsgs, errMsg)
		}

		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Validation error",
			"errors":  errMsgs,
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

	alreadyLogin := c.Cookies("token")
	if alreadyLogin != "" {
		Logout(c)
	}

	var user models.Users
	if err := c.BodyParser(&user); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Invalid user data",
		})
	}

	validate := validator.New()
	validate.RegisterValidation("usernameValid", Middleware.UsernameValid)
	if err := validate.Struct(user); err != nil {
		var errMsgs []string
		for _, err := range err.(validator.ValidationErrors) {
			fieldName := strings.ToLower(err.Field())
			errMsg := fmt.Sprintf("Field validation for '%s' failed on the '%s' tag", fieldName, err.Tag())
			errMsgs = append(errMsgs, errMsg)
		}

		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Validation error",
			"errors":  errMsgs,
		})
	}

	var dbUser = models.Users{}
	db.DB.Where("username = ?", user.Username).First(&dbUser)

	if dbUser.Id == 0 {
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
		"userId":    strconv.Itoa(int(dbUser.Id)),
		"username":  dbUser.Username,
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
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "User id is required",
		})
	}

	var user models.Users
	if err := db.DB.Where("id = ?", userId).First(&user).Error; err != nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": "User not found",
		})
	}

	var posts []models.Posts
	db.DB.Preload("Comments").Preload("Likes").Order("created_at desc").Where("user_id = ?", userId).Find(&posts)

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "User Posts listed successfully",
		"data":    posts,
	})
}

func GetUserComments(c *fiber.Ctx) error {
	userId := c.Params("id")

	if userId == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "User id is required",
		})
	}

	var user models.Users
	if err := db.DB.Where("id = ?", userId).First(&user).Error; err != nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": "User not found",
		})
	}

	var comments []models.PostCommets
	db.DB.Order("created_at desc").Where("user_id = ?", userId).Find(&comments)

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "User Comments listed successfully",
		"data":    comments,
	})
}

func GetUserLikes(c *fiber.Ctx) error {
	userId := c.Params("id")

	if userId == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "User id is required",
		})
	}

	var user models.Users
	if err := db.DB.Where("id = ?", userId).First(&user).Error; err != nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": "User not found",
		})
	}

	var likes []models.PostLikes
	db.DB.Where("user_id = ?", userId).Find(&likes)

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "User Likes listed successfully",
		"data":    likes,
	})
}
