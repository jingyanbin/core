package basal

import (
	"os"
	"path/filepath"
)

const (
	FILE_FLAG_WRITER = os.O_WRONLY | os.O_APPEND | os.O_CREATE
	FILE_PERM_ALL    = 0777
)

var ErrNotFolder = NewError("not folder")
var ErrNotFile = NewError("not file")

// 打开文件 自动创建目录
func OpenFile(folderPath string, fileName string, flag int, perm os.FileMode) (file *os.File, err error) {
	fp := filepath.Join(folderPath, fileName)
	file, err = os.OpenFile(fp, flag, perm)
	if err != nil {
		var has bool
		has, err = IsExistFolder(folderPath)
		if err != nil {
			return
		}
		if has == false {
			err = os.MkdirAll(folderPath, os.ModeDir)
			if err != nil {
				return
			}
			file, err = os.OpenFile(fp, flag, perm)
		}
	}
	if !IsExistBy(file, err) {
		err = os.ErrNotExist
	}
	return
}

// 打开文件 自动创建目录
func OpenFileB(filePath string, flag int, perm os.FileMode) (file *os.File, err error) {
	folderPath, fileName := filepath.Split(filePath)
	return OpenFile(folderPath, fileName, flag, perm)
}

// 文件是否存在
func IsExistBy(f *os.File, err error) bool {
	if f == nil || os.IsNotExist(err) {
		return false
	}
	return true
}

// 文件或文件夹是否存在
func IsExist(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

// 文件夹是否存在
func IsExistFolder(path string) (bool, error) {
	handle, err := os.Stat(path)
	if err == nil {
		if handle.IsDir() {
			return true, nil
		} else {
			return false, ErrNotFolder
		}
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

// 文件是否存在
func IsExistFile(path string) (bool, error) {
	handle, err := os.Stat(path)
	if err == nil {
		if handle.IsDir() {
			return false, ErrNotFile
		} else {
			return true, nil
		}
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
