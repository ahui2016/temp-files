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
	api.Post("/files", getFileList)
	api.Post("/upload-file", uploadFileHandler)
	api.Post("/delete-file", deleteFile)
	api.Post("/download-file", downloadFile)

	log.Fatal(app.Listen(app_config.Host))
}
