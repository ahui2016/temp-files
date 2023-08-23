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

	api.Post("/check-pwd", checkPassword)
	api.Get("/files", getFileList)
	api.Post("/upload-file", uploadFileHandler)

	log.Fatal(app.Listen(app_config.Host))
}
