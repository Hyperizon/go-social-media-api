package Routes

import (
	Controllers "go-social-media-api/Controllers"
	Middlewares "go-social-media-api/Middlewares"

	"github.com/gofiber/fiber/v2"
)

func Setup(app *fiber.App) {
	api := app.Group("/api")

	//Users Routes
	api.Post("/users/register", Controllers.Register)
	api.Post("/users/login", Controllers.Login)
	api.Get("/users/logout", Middlewares.Auth, Controllers.Logout)
	api.Get("/users/:id/posts", Middlewares.Auth, Controllers.GetUserPosts)

	//Posts Routes
	api.Post("/posts", Middlewares.Auth, Controllers.CreatePost)
	api.Get("/posts", Middlewares.Auth, Controllers.GetPosts)
	api.Get("/posts/:id", Middlewares.Auth, Controllers.GetPostById)
	api.Put("/posts/:id", Middlewares.Auth, Controllers.UpdatePost)
	api.Delete("/posts/:id", Middlewares.Auth, Controllers.DeletePost)
}
