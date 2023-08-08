package Controllers

import (
	"errors"
	"fmt"
	db "go-social-media-api/Config"
	models "go-social-media-api/Models"
	"go-social-media-api/Utils"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
	"gorm.io/gorm"
)

func CreatePost(c *fiber.Ctx) error {
	var post models.Posts
	if err := c.BodyParser(&post); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Invalid post data",
		})
	}

	validate := validator.New()
	if err := validate.Struct(post); err != nil {
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

	file, _ := c.FormFile("Image")

	if file != nil {
		imageBase64, err := Utils.UploadImage(file)
		if err != nil {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{
				"success": false,
				"message": "Error uploading image",
			})
		}

		post.Image = imageBase64
	}

	tokenString := c.Cookies("token")

	token, _ := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	claims, _ := token.Claims.(jwt.MapClaims)
	userId, err := strconv.Atoi(claims["userId"].(string))
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Error getting user",
		})
	}
	post.UserId = uint(userId)
	post.CreatedAt = time.Now()
	post.UpdatedAt = time.Now()

	db.DB.Create(&post)
	return c.Status(http.StatusCreated).JSON(fiber.Map{
		"success": true,
		"message": "Post created successfully",
		"data":    post,
	})
}

func GetPosts(c *fiber.Ctx) error {
	var posts []models.Posts
	db.DB.Preload("Comments").Preload("Likes").Order("created_at desc").Find(&posts)

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Posts listed successfully",
		"data":    posts,
	})
}

func GetPostById(c *fiber.Ctx) error {
	postId := c.Params("id")
	if postId == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Post id is required",
		})
	}

	var posts models.Posts
	if err := db.DB.Preload("Comments").Preload("Likes").Where("id = ?", postId).First(&posts).Error; err != nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": "Post not found",
		})
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Post listed successfully",
		"data":    posts,
	})
}

func UpdatePost(c *fiber.Ctx) error {
	postId := c.Params("id")
	if postId == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Post id is required",
		})
	}

	var existingPost models.Posts
	if err := db.DB.Where("id = ?", postId).First(&existingPost).Error; err != nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": "Post not found",
		})
	}

	var updateData map[string]string
	if err := c.BodyParser(&updateData); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Invalid update data",
		})
	}

	if description, ok := updateData["description"]; ok {
		existingPost.Description = description
	}
	existingPost.UpdatedAt = time.Now()

	db.DB.Save(&existingPost)

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Post updated successfully",
		"data":    existingPost,
	})
}

func DeletePost(c *fiber.Ctx) error {
	postId := c.Params("id")
	if postId == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Post id is required",
		})
	}

	tokenString := c.Cookies("token")

	token, _ := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	claims, _ := token.Claims.(jwt.MapClaims)
	userId, err := strconv.Atoi(claims["userId"].(string))
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Error getting user",
		})
	}
	var post models.Posts
	if err := db.DB.Where("id = ?", postId).Where("user_id = ?", userId).First(&post).Error; err != nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": "Post not found",
		})
	}

	db.DB.Delete(&post)
	return c.Status(http.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Post deleted successfully",
	})
}

func LikePost(c *fiber.Ctx) error {
	postId := c.Params("id")
	if postId == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Post id is required",
		})
	}

	var post models.Posts
	if err := db.DB.Where("id = ?", postId).First(&post).Error; err != nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": "Post not found",
		})
	}

	tokenString := c.Cookies("token")

	token, _ := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	claims, _ := token.Claims.(jwt.MapClaims)
	userId, err := strconv.Atoi(claims["userId"].(string))
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Error getting user",
		})
	}

	var postLike models.PostLikes
	err = db.DB.Where("post_id = ?", post.Id).Where("user_id = ?", userId).First(&postLike).Error
	if err == nil {
		post.LikesCount--
		db.DB.Delete(&postLike)
	} else if errors.Is(err, gorm.ErrRecordNotFound) {
		postLike = models.PostLikes{
			PostId: post.Id,
			UserId: uint(userId),
		}
		db.DB.Create(&postLike)
		post.LikesCount++
	} else {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "An error occurred",
		})
	}

	db.DB.Save(&post)

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Post liked successfully",
		"data":    postLike,
	})
}

func CommentToPost(c *fiber.Ctx) error {
	fmt.Println("CommentToPost")
	postId := c.Params("id")
	if postId == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Post id is required",
		})
	}

	var post models.Posts
	if err := db.DB.Where("id = ?", postId).First(&post).Error; err != nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": "Post not found",
		})
	}

	tokenString := c.Cookies("token")

	token, _ := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	claims, _ := token.Claims.(jwt.MapClaims)
	userId, err := strconv.Atoi(claims["userId"].(string))
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Error getting user",
		})
	}

	var postComment models.PostCommets
	if err := c.BodyParser(&postComment); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Invalid comment data",
		})
	}

	if postComment.Comment == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Comment is required",
		})
	}

	postComment.PostId = post.Id
	postComment.UserId = uint(userId)
	postComment.CreatedAt = time.Now()
	postComment.UpdatedAt = time.Now()

	db.DB.Create(&postComment)
	return c.Status(http.StatusCreated).JSON(fiber.Map{
		"success": true,
		"message": "Comment created successfully",
		"data":    postComment,
	})
}
