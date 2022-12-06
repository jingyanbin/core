package basal

import (
	"os"
	"path"
	"path/filepath"
	"runtime"
	_ "unsafe"
)

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

var ErrNotFolder = NewError("not folder")

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

var ErrNotFile = NewError("not file")

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

type HandleFile struct {
	folderPath string
	fileName   string
	flag       int
	perm       os.FileMode
	handle     *os.File
}

func (my *HandleFile) PathName() string {
	return filepath.Join(my.folderPath, my.fileName)
}

func (my *HandleFile) SetPathName(folderPath, fileName string) bool {
	if my.folderPath == folderPath && my.fileName == fileName {
		return false
	}
	my.Close()
	my.folderPath = folderPath
	my.fileName = fileName
	return true
}

func (my *HandleFile) WriteString(s string) (n int, err error) {
	n, err = my.Write([]byte(s))
	return
}

func (my *HandleFile) Write(b []byte) (n int, err error) {
	if my.handle == nil {
		my.handle, err = OpenFile(my.folderPath, my.fileName, my.flag, my.perm)
		if err != nil {
			return
		}
	} else {
		var hasLog bool
		hasLog, err = IsExist(filepath.Join(my.folderPath, my.fileName))
		if !hasLog {
			my.handle.Close()
			my.handle, err = OpenFile(my.folderPath, my.fileName, my.flag, my.perm)
			if err != nil {
				return
			}
		}
	}
	n, err = my.handle.Write(b)
	return
}

func (my *HandleFile) Close() {
	if my.handle == nil {
		return
	}
	my.handle.Close()
	my.handle = nil
}

const HANDLE_FILE_FLAG_WRITER = os.O_WRONLY | os.O_APPEND | os.O_CREATE
const HANDLE_FILE_PERM_ALL = 0777

func NewHandleFile(flag int, perm os.FileMode) *HandleFile {
	return &HandleFile{flag: flag, perm: perm}
}

func OpenHandleFile(folderPath string, fileName string, flag int, perm os.FileMode) (*HandleFile, error) {
	var err error
	if folderPath == "" {
		err = NewError("OpenHandleFile Error: folderPath is nil")
		return nil, err
	}
	if fileName == "" {
		err = NewError("OpenHandleFile Error: fileName is nil")
		return nil, err
	}

	hf := &HandleFile{folderPath: folderPath, fileName: fileName, flag: flag, perm: perm}

	hf.handle, err = OpenFile(folderPath, fileName, flag, perm)
	if err != nil {
		return nil, err
	}
	return hf, nil
}

type fpath struct {
	programDir string
	execDir    string
}

func (m *fpath) init() {
	var err error
	m.programDir, err = filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		panic(err)
	}
	m.execDir, err = os.Getwd()
	if err != nil {
		panic(err)
	}
}

// 程序路径
func (m *fpath) ProgramDir() string {
	return m.programDir
}

// 执行路径
func (m *fpath) ExecDir() string {
	return m.execDir
}

// 文件名
func (m *fpath) Base(p string) string {
	if runtime.GOOS == "windows" {
		return path.Base(filepath.ToSlash(p))
	} else {
		return path.Base(p)
	}
}

// 程序路径组合
func (m *fpath) ProgramDirJoin(elem ...string) string {
	var elems = make([]string, len(elem)+1)
	elems[0] = m.programDir
	copy(elems[1:], elem)
	return filepath.Join(elems...)
}

func (m *fpath) Join(elem ...string) string {
	return filepath.Join(elem...)
}

// 执行路径组合
func (m *fpath) ExecDirJoin(elem ...string) string {
	var elems = make([]string, len(elem)+1)
	elems[0] = m.execDir
	copy(elems[1:], elem)
	return filepath.Join(elems...)
}

// 获得路径最后的文件、文件件所在目录
func (m *fpath) Dir(p string) string {
	return filepath.Dir(p)
}

var Path = fpath{}

func init() {
	Path.init()
}
