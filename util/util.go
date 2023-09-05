package util

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"slices"

	"github.com/samber/lo"
)

// WrapErrors 把多个错误合并为一个错误.
func WrapErrors(allErrors ...error) (wrapped error) {
	for _, err := range allErrors {
		if err != nil {
			if wrapped == nil {
				wrapped = err
			} else {
				wrapped = fmt.Errorf("%w | %w", wrapped, err)
			}
		}
	}
	return
}

// TODO: 改名 PathNotExists
func PathNotExists(name string) (ok bool) {
	_, err := os.Lstat(name)
	if os.IsNotExist(err) {
		ok = true
		err = nil
	}
	lo.Must0(err)
	return
}

func PathExists(name string) bool {
	return !PathNotExists(name)
}

// https://stackoverflow.com/questions/30376921/how-do-you-copy-a-file-in-go
func CopyFile(dstPath, srcPath string) error {
	src, err := os.Open(srcPath)
	if err != nil {
		return err
	}
	defer src.Close()

	dst, err := os.Create(dstPath)
	if err != nil {
		return err
	}
	defer dst.Close()

	_, err1 := io.Copy(dst, src)
	err2 := dst.Sync()
	return WrapErrors(err1, err2)
}

func SortStrings(x []string) {
	if slices.IsSorted(x) {
		return
	}
	slices.Sort(x)
}

// DeleteOldFiles 刪除 folder 裡 n 個最舊的檔案。
func DeleteOldFiles(folder string, n int) error {
	Separator := string(filepath.Separator)
	files, err := filepath.Glob(folder + Separator + "*")
	if err != nil {
		return err
	}
	SortStrings(files)
	if len(files) < n {
		n = len(files)
	}
	for _, f := range files[:n] {
		if err := os.Remove(f); err != nil {
			return err
		}
	}
	return nil
}

// RemainNewFiles 刪除 folder 裡的舊檔案, 剩下 n 個最新的檔案。
func RemainNewFiles(folder string, n int64) error {
	Separator := string(filepath.Separator)
	files, err := filepath.Glob(folder + Separator + "*")
	if err != nil {
		return err
	}
	SortStrings(files)

	length := int64(len(files))
	if length <= n {
		return nil
	}

	n_del := length - n
	for _, f := range files[:n_del] {
		if err := os.Remove(f); err != nil {
			return err
		}
	}
	return nil
}
