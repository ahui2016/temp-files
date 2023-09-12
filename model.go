package main

import (
	"mime/multipart"
	"os"
	"slices"
	"strconv"
	"strings"
	"time"
)

var MediaSuffixList = []string{
	"html",
	"htm",
	"css",
	"xml",
	"gif",
	"jpeg",
	"jpg",
	"atom",
	"rss",
	"txt",
	"png",
	"svg",
	"webp",
	"ico",
	"bmp",
	"js",
	"json",
	"pdf",
	"mp3",
	"m4a",
	"mp4",
	"mpeg",
	"mpg",
	"flv",
	"avi",
}

type AppConfig struct {
	Host               string
	UploadLimit        int64 // MB
	Password           string
	RecentFilesLimit   int64
	OldTextFilesLimit  int64
	RepeatRequestLimit int
}

type File struct {
	CTime   int64  // 服务器保存该文件的时间
	Name    string // 原文件名
	Size    int64  // length in bytes for regular files
	IsText  bool   // true if File.Name ends with ".txt" or ".md"
	IsMedia bool   // true if File.Name ends with ".jpg" or ".pdf" etc.
}

// NewFileWithName 主要是为了更方便使用 File.TimeName(), 因此 Size 不重要。
func NewFileWithName(name string) *File {
	f := &File{
		CTime: time.Now().Unix(),
		Name:  name,
		// Size:  不重要
	}
	suffix := lastElement(strings.Split(f.Name, "."))
	if strings.HasSuffix(suffix, "txt") || strings.HasSuffix(suffix, "md") {
		f.IsText = true
	}
	if slices.Contains(MediaSuffixList, suffix) {
		f.IsMedia = true
	}
	return f
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
	if err != nil {
		return nil, err
	}
	f.Size = info.Size()
	suffix := lastElement(strings.Split(f.Name, "."))
	if strings.HasSuffix(suffix, "txt") || strings.HasSuffix(suffix, "md") {
		f.IsText = true
	}
	if slices.Contains(MediaSuffixList, suffix) {
		f.IsMedia = true
	}
	return f, nil
}

func lastElement(s []string) string {
	return s[len(s)-1]
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
	Old      bool   `json:"old" form:"old"`
}

type FilePrefixForm struct {
	Prefix string `json:"prefix" form:"prefix" validate:"required"`
	Old    bool   `json:"old" form:"old"`
}

type FileWithContent struct {
	CTime   string `json:"ctime" form:"ctime"`
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
