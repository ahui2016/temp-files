package main

import (
	"mime/multipart"
	"os"
	"strconv"
	"strings"
	"time"
)

type AppConfig struct {
	Host        string
	UploadLimit int64 // MB
	Password    string
}

type File struct {
	CTime  int64  // 服务器保存该文件的时间
	Name   string // 原文件名
	Size   int64  // length in bytes for regular files
	IsText bool   // true if File.Name ends with ".txt" or ".md"
}

// NewFileWithName 主要是为了更方便使用 File.TimeName(),
// 因此 Size 和 IsText 不重要。
func NewFileWithName(name string) *File {
	return &File{
		CTime: time.Now().Unix(),
		Name:  name,
		// Size:  不重要
		// IsText: 此时文件类型不重要
	}
}

// NewFileFromUser 根据用户上传的文件解析出一个 File
func NewFileFromUser(file *multipart.FileHeader) *File {
	return &File{
		CTime: time.Now().Unix(),
		Name:  file.Filename,
		Size:  file.Size,
		// IsText: 此时文件类型不重要
	}
}

// NewFileFromServer 根据服务器中的文件名解析出一个 File
func NewFileFromServer(filePath string) (*File, error) {
	f := new(File)
	info, err := os.Stat(filePath)
	if err != nil {
		return nil, err
	}
	f.CTime, f.Name, err = splitTimeName(info.Name())
	f.Size = info.Size()
	if strings.HasSuffix(f.Name, ".txt") || strings.HasSuffix(f.Name, ".md") {
		f.IsText = true
	}
	return f, err
}

func splitTimeName(timeName string) (time int64, name string, err error) {
	array := strings.SplitN(timeName, "-", 2)
	if len(array) == 2 {
		time, err = strconv.ParseInt(array[0], 10, 64)
		return time, array[1], err
	}
	return 0, timeName, nil
}

// TimeName 返回 CTime-Name (用连字号连接两个字符串)
// 用来作为保存在服务器时的文件名.
func (f *File) TimeName() string {
	return strconv.FormatInt(f.CTime, 10) + "-" + f.Name
}

type FilenameForm struct {
	Filename string `json:"filename" form:"filename" validate:"required"`
}

type FileWithContent struct {
	Name    string `json:"name" form:"name" validate:"required"`
	IsText  bool   `json:"isText" form:"isText"`
	Content string `json:"content" form:"content"`
}

func NewFileWithContent(name string) *FileWithContent {
	file := &FileWithContent{Name: name}
	if strings.HasSuffix(name, ".txt") || strings.HasSuffix(name, ".md") {
		file.IsText = true
	}
	return file
}
