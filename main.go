package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New(fiber.Config{
		Immutable: true, // 以后试试删除该设定
	})

	app.Static("/", public_folder)

	api := app.Group("/api")

	api.Get("/files", getFileList)

	log.Fatal(app.Listen(app_config.Host))
}
