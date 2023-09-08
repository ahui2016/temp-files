package util

import (
	"archive/zip"
	"io"
	"os"
	"path/filepath"
)

// ZipFile 方便自定义每个文件的文件名。
type ZipFile struct {
	Name string
	Path string
}

// ZipPaths 创建一个压缩包，把 paths 的文件打包到 zipFilePath.
func ZipPaths(zipFilePath string, paths []string) error {
	files := pathsToZipFiles(paths)
	return ZipFiles(zipFilePath, files)
}

// ZipFiles 创建一个压缩包，把 files 打包到 zipFilePath.
func ZipFiles(zipFilePath string, files []ZipFile) error {
	zipFile, err := os.Create(zipFilePath)
	if err != nil {
		return err
	}
	defer zipFile.Close()

	zipWriter := zip.NewWriter(zipFile)
	for _, file := range files {
		addZipFile(file, zipWriter)
	}

	// 如果发生错误，要删除文件。
	if err := zipWriter.Close(); err != nil {
		zipFile.Close()
		if e2 := os.Remove(zipFilePath); e2 != nil {
			err = WrapErrors(err, e2)
		}
		return err
	}
	return nil
}

func pathsToZipFiles(paths []string) (zipFiles []ZipFile) {
	for _, path := range paths {
		f := ZipFile{
			Name: filepath.Base(path),
			Path: path,
		}
		zipFiles = append(zipFiles, f)
	}
	return
}

func addZipFile(file ZipFile, zipWriter *zip.Writer) error {
	fileIndeed, err := os.Open(file.Path)
	if err != nil {
		return err
	}
	defer fileIndeed.Close()

	body, err := io.ReadAll(fileIndeed)
	if err != nil {
		return err
	}

	fileInZip, err := zipWriter.Create(file.Name)
	if err != nil {
		return err
	}

	if _, err = fileInZip.Write(body); err != nil {
		return err
	}
	return nil
}
