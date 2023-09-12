package main

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"

	"github.com/ahui2016/temp-files/util"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/samber/lo"
)

// TextMsg 用于向前端返回一个简单的文本消息。
type TextMsg struct {
	Text string `json:"text"`
}

// TextMsg 用于向前端返回一个简单的整数。
type Int64Msg struct {
	Data int64 `json:"data"`
}

var validate = validator.New()

// noCache is a middleware.
// Cache-Control: no-store will refrain from caching.
// You will always get the up-to-date response.
func noCache(c *fiber.Ctx) error {
	c.Set("Cache-Control", "no-store")
	return c.Next()
}

func checkLoginHandler(c *fiber.Ctx) error {
	if isLoggedIn(c) {
		return nil
	}
	return fmt.Errorf("Require Login")
}

func loginHandler(c *fiber.Ctx) error {
	if isLoggedIn(c) {
		return nil
	}
	if err := checkPassword(c); err != nil {
		return err
	}
	return sessionSet(c, cookieName, true)
}

func logoutHandler(c *fiber.Ctx) error {
	if isLoggedOut(c) {
		return nil
	}
	return sessionDelete(c, cookieName)
}

func getAllFiles(c *fiber.Ctx) error {
	files, err := allFiles("")
	if err != nil {
		return err
	}
	return c.JSON(files)
}

func getOldTextFiles(c *fiber.Ctx) error {
	files, err := allFiles("old")
	if err != nil {
		return err
	}
	slices.Reverse(files)
	return c.JSON(files)
}

func getRecentFiles(c *fiber.Ctx) error {
	files, err := allFiles("recent")
	if err != nil {
		return err
	}
	slices.Reverse(files)
	return c.JSON(files)
}

func getTotalSize(c *fiber.Ctx) error {
	paths, err := filepath.Glob(files_folder + Separator + "*")
	if err != nil {
		return err
	}
	old_paths, err := filepath.Glob(old_text_files_folder + Separator + "*")
	if err != nil {
		return err
	}
	paths = append(paths, old_paths...)

	var size int64
	for _, f := range paths {
		info, err := os.Stat(f)
		if err != nil {
			return err
		}
		size += info.Size()
	}
	return c.JSON(Int64Msg{size})
}

func deleteFile(c *fiber.Ctx) error {
	form, err := checkParseFilename(c)
	if err != nil {
		return err
	}
	folder := lo.Ternary(form.Old, old_text_files_folder, files_folder)
	filePath := filepath.Join(folder, form.Filename)
	return os.Remove(filePath)
}

func downloadFile(c *fiber.Ctx) error {
	form, err := checkParseFilename(c)
	if err != nil {
		return err
	}
	filePath := filepath.Join(files_folder, form.Filename)
	return c.SendFile(filePath)
}

func loadFileHandler(c *fiber.Ctx) error {
	form := new(FilePrefixForm)
	if err := parseValidate(form, c); err != nil {
		return err
	}
	filePath, f, err := getFileByPrefix(form.Prefix, form.Old)
	if err != nil {
		return err
	}
	file := NewFileWithContent(f.Name)
	if file.IsText {
		content, err := os.ReadFile(filePath)
		if err != nil {
			return err
		}
		file.Content = string(content)
	}
	return c.JSON(file)
}

func getFileByPrefix(prefix string, old bool) (filePath string, file *File, err error) {
	folder := lo.Ternary(old, old_text_files_folder, files_folder)
	pattern := filepath.Join(folder, prefix)
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

func checkParseFilename(c *fiber.Ctx) (*FilenameForm, error) {
	form := new(FilenameForm)
	if err := parseValidate(form, c); err != nil {
		return nil, err
	}
	return form, nil
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

func saveTextFile(c *fiber.Ctx) error {
	form := new(FileWithContent)
	if err := parseValidate(form, c); err != nil {
		return err
	}
	if int64(len(form.Content)) > app_config.UploadLimit*MB {
		return fmt.Errorf("the file is too large (> %d MB)", app_config.UploadLimit)
	}
	file := NewFileWithName(form.Name)
	if !strings.Contains(file.Name, ".") {
		// 如果忘记写后缀名，就自动添加。
		file.Name += ".txt"
	}
	filePath := filepath.Join(files_folder, file.TimeName())
	if err := os.WriteFile(filePath, []byte(form.Content), 0666); err != nil {
		return err
	}
	if err := moveOldTextFile(form.CTime); err != nil {
		return err
	}
	return c.JSON(TextMsg{strconv.FormatInt(file.CTime, 10)})
}

func zipTextFiles(c *fiber.Ctx) error {
	files, err := getTextFiles()
	if err != nil {
		return err
	}
	file := NewFileWithName("all-text-files.zip")
	zipFilePath := filepath.Join(files_folder, file.TimeName())
	return util.ZipPaths(zipFilePath, files)
}

func moveOldTextFile(ctime string) error {
	if ctime == "" {
		return nil
	}
	oldpath, file, err := getFileByPrefix(ctime+"-*", false)
	if err != nil {
		return err
	}
	// 更新时间，保证文件移动后位于列表顶端。
	newfile := NewFileWithName(file.Name)
	newpath := filepath.Join(old_text_files_folder, newfile.TimeName())
	if err := os.Rename(oldpath, newpath); err != nil {
		return err
	}
	return util.RemainNewFiles(old_text_files_folder, (app_config.OldTextFilesLimit))
}

func allFiles(filter string) ([]*File, error) {
	folder := lo.Ternary(filter == "old", old_text_files_folder, files_folder)
	paths, err := filepath.Glob(folder + Separator + "*")
	if err != nil {
		return nil, err
	}
	util.SortStrings(paths)
	if filter == "recent" {
		if int64(len(paths)) > app_config.RecentFilesLimit {
			paths = paths[:app_config.RecentFilesLimit]
		}
	}
	if len(paths) == 0 {
		return nil, nil
	}
	return pathsToFiles(paths)
}

func getTextFiles() (textFiles []string, err error) {
	paths, err := filepath.Glob(files_folder + Separator + "*")
	if err != nil {
		return
	}
	if len(paths) == 0 {
		return nil, nil
	}
	for _, path := range paths {
		file := NewFileWithName(filepath.Base(path))
		if file.IsText {
			textFiles = append(textFiles, path)
		}
	}
	return
}

func pathsToFiles(paths []string) (files []*File, err error) {
	for _, filePath := range paths {
		f, err := NewFileFromServer(filePath)
		if err != nil {
			return nil, err
		}
		files = append(files, f)
	}
	return
}

func parseValidate(form any, c *fiber.Ctx) error {
	if err := c.BodyParser(form); err != nil {
		return err
	}
	return validate.Struct(form)
}
