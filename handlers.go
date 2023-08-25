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
	filename, err := checkParseFilename(c)
	if err != nil {
		return err
	}
	filePath := filepath.Join(files_folder, filename)
	return os.Remove(filePath)
}

func downloadFile(c *fiber.Ctx) error {
	filename, err := checkParseFilename(c)
	if err != nil {
		return err
	}
	filePath := filepath.Join(files_folder, filename)
	return c.SendFile(filePath)
}

func loadFileHandler(c *fiber.Ctx) error {
	prefix, err := checkParseFilename(c)
	if err != nil {
		return err
	}
	filePath, file, err := getFileByPrefix(prefix)
	if err != nil {
		return err
	}
	content, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}
	return c.JSON(FileWithContent{
		Name:    file.Name,
		Content: string(content),
	})
}

func getFileByPrefix(prefix string) (filePath string, file *File, err error) {
	pattern := filepath.Join(files_folder, prefix)
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return
	}
	if len(matches) < 1 {
		err = fmt.Errorf("file not found: %s", prefix)
		return
	}
	filePath = matches[0]
	file, err = NewFileFromServer(filePath)
	return
}

func checkParseFilename(c *fiber.Ctx) (filename string, err error) {
	if err = checkPassword(c); err != nil {
		return
	}
	form := new(FilenameForm)
	if err = parseValidate(form, c); err != nil {
		return
	}
	return form.Filename, nil
}

func uploadFileHandler(c *fiber.Ctx) error {
	if err := checkPassword(c); err != nil {
		return err
	}
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

func saveTextFile(c *fiber.Ctx) error {
	if err := checkPassword(c); err != nil {
		return err
	}
	form := new(FileWithContent)
	if err := parseValidate(form, c); err != nil {
		return err
	}
	if int64(len(form.Content)) > app_config.UploadLimit*MB {
		return fmt.Errorf("the file is too large (> %d MB)", app_config.UploadLimit)
	}
	fileName := NewFileWithName(form.Name).TimeName()
	filePath := filepath.Join(files_folder, fileName)
	return os.WriteFile(filePath, []byte(form.Content), 0666)
}

func renameFile(c *fiber.Ctx) error {}

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
