package main

import (
	db "go-social-media-api/Config"
	routes "go-social-media-api/Routes"

	"github.com/gofiber/fiber/v2"
)

func main() {

	db.Connect()

	app := fiber.New()

	routes.Setup(app)

	app.Listen(":3000")
}
