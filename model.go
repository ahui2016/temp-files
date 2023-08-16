package main

type AppConfig struct {
	Host string
}

type File struct {
	CTime  int64  // 服务器保存该文件的时间
	Name   string // 原文件名
	Size   int64  // length in bytes for regular files
	IsText bool   // true if File.Name ends with ".txt" or ".md"
}

// NewFile 根据服务器中的文件名解析出一个 File
func NewFile(timeName string) *File {}

// TimeName 返回 CTime-Name (用连字号连接两个字符串)
// 用来作为保存在服务器时的文件名.
func (f *File) TimeName() string {}
