package main

import (
	"fmt"
	"path/filepath"
	"slices"

	"github.com/gofiber/fiber/v2"
)

// noCache is a middleware.
// Cache-Control: no-store will refrain from caching.
// You will always get the up-to-date response.
func noCache(c *fiber.Ctx) error {
	c.Set("Cache-Control", "no-store")
	return c.Next()
}

func getFileList(c *fiber.Ctx) error {
	files, err := allFiles()
	if err != nil {
		return err
	}
	slices.Reverse(files)
	return c.JSON(files)
}

func uploadFileHandler(c *fiber.Ctx) error {
	file, err := c.FormFile("file")
	if err != nil {
		return err
	}
	if file.Size > app_config.UploadLimit*MB {
		return fmt.Errorf("the file is too large (> %d MB)", app_config.UploadLimit)
	}
	f := NewFileFromUser(file)
	filePath := filepath.Join(files_folder, f.TimeName())
	return c.SaveFile(file, filePath)
}

func allFiles() (files []*File, err error) {
	paths, err := filepath.Glob(files_folder + Separator + "*")
	if err != nil {
		return nil, err
	}
	for _, filePath := range paths {
		f, err := NewFileFromServer(filePath)
		if err != nil {
			return nil, err
		}
		files = append(files, f)
	}
	return
}
