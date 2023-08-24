package main

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

var validate = validator.New()

// noCache is a middleware.
// Cache-Control: no-store will refrain from caching.
// You will always get the up-to-date response.
func noCache(c *fiber.Ctx) error {
	c.Set("Cache-Control", "no-store")
	return c.Next()
}

func getFileList(c *fiber.Ctx) error {
	if err := checkPassword(c); err != nil {
		return err
	}
	files, err := allFiles()
	if err != nil {
		return err
	}
	slices.Reverse(files)
	return c.JSON(files)
}

func deleteFile(c *fiber.Ctx) error {
	if err := checkPassword(c); err != nil {
		return err
	}
	type Form struct {
		Filename string `json:"filename" form:"filename" validate:"required"`
	}
	form := new(Form)
	if err := parseValidate(form, c); err != nil {
		return err
	}
	filePath := filepath.Join(files_folder, form.Filename)
	return os.Remove(filePath)
}

func uploadFileHandler(c *fiber.Ctx) error {
	file, err := c.FormFile("file")
	if err != nil {
		return err
	}
	if err := checkPassword(c); err != nil {
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

func checkPassword(c *fiber.Ctx) error {
	type Pass struct {
		Word string `json:"pwd" form:"pwd"`
	}
	pwd := new(Pass)
	if err := c.BodyParser(pwd); err != nil {
		return err
	}
	if pwd.Word != app_config.Password {
		return fmt.Errorf("wrong password")
	}
	return nil
}

func parseValidate(form any, c *fiber.Ctx) error {
	if err := c.BodyParser(form); err != nil {
		return err
	}
	return validate.Struct(form)
}
