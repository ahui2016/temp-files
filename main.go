package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New(fiber.Config{
		Immutable: true, // 以后试试删除该设定
		BodyLimit: int(app_config.UploadLimit) * MB,
	})

	app.Use(noCache)

	app.Static("/", public_folder)
	app.Static("/"+FilesFolderName, files_folder)
	app.Static("/"+OldTextFilesFolderName, old_text_files_folder)

	api := app.Group("/api")

	api.Post("/check-pwd", checkPassword)
	api.Post("/all-files", getAllFiles)
	api.Post("/recent-files", getRecentFiles)
	api.Post("/old-text-files", getOldTextFiles)
	api.Post("/total-size", getTotalSize)
	api.Post("/upload-file", uploadFileHandler)
	api.Post("/delete-file", deleteFile)
	api.Post("/download-file", downloadFile)
	api.Post("/load-file-by-prefix", loadFileHandler)
	api.Post("/save-text-file", saveTextFile)
	api.Post("/zip-text-files", zipTextFiles)

	log.Fatal(app.Listen(app_config.Host))
}
