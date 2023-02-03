package basal

import (
	"os"
	"path"
	"path/filepath"
	"runtime"
)

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
