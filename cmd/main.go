package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"

	"github.com/mehdighachoui/url-shortener/handlers"
)

func main() {
	action := handlers.NewAction()
	action.RedisConnection()

	engine := html.New("./templates", ".html")
	app := fiber.New(fiber.Config{Views: engine})

	app.Get("/", func(f *fiber.Ctx) error {
		return f.Render("index", fiber.Map{})
	})
	app.Get("/:url", action.GetUrlHandler)
	app.Post("/shorten", action.CreateUrlHandler)

	app.Listen(":8001")
}
