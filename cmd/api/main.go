package main

import (
	"github.com/gofiber/fiber/v3"
	"log"
)

func HandleRoute(c fiber.Ctx) error {
	return c.SendString("Hello, World 👋!")
}

func main() {
	app := fiber.New()

	app.Get("/", HandleRoute)

	log.Fatal(app.Listen(":3000"))
}
