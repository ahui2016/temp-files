package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
)

func main() {
	app := fiber.New(fiber.Config{
		Immutable: true, // 以后试试删除该设定
		BodyLimit: int(app_config.UploadLimit) * MB,
	})

	limitHandler := limiter.New(limiter.Config{
		Max: app_config.RepeatRequestLimit,
	})

	app.Use(noCache)

	app.Static("/", public_folder)

	app.Use("/"+FilesFolderName, checkLoginMidWare)
	app.Static("/"+FilesFolderName, files_folder)

	app.Use("/"+OldTextFilesFolderName, checkLoginMidWare)
	app.Static("/"+OldTextFilesFolderName, old_text_files_folder)

	app.Get("/check-login", checkLoginHandler)
	app.Get("/logout", logoutHandler)

	app.Post("/login", limitHandler, loginHandler)

	api := app.Group("/api", checkLoginMidWare)

	api.Get("/all-files", getAllFiles)
	api.Get("/recent-files", getRecentFiles)
	api.Get("/old-text-files", getOldTextFiles)
	api.Get("/total-size", getTotalSize)
	api.Post("/upload-file", uploadFileHandler)
	api.Post("/delete-file", deleteFile)
	api.Post("/download-file", downloadFile)
	api.Post("/load-file-by-prefix", loadFileHandler)
	api.Post("/save-text-file", saveTextFile)
	api.Get("/zip-text-files", zipTextFiles)

	log.Fatal(app.Listen(app_config.Host))
}
