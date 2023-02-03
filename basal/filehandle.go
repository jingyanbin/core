package basal

import (
	"os"
	"path/filepath"
	_ "unsafe"
)

type FileHandle struct {
	folderPath string
	fileName   string
	flag       int
	perm       os.FileMode
	handle     *os.File
}

func (m *FileHandle) PathName() string {
	return filepath.Join(m.folderPath, m.fileName)
}

func (m *FileHandle) SetPathName(folderPath, fileName string) bool {
	if m.folderPath == folderPath && m.fileName == fileName {
		return false
	}
	m.Close()
	m.folderPath = folderPath
	m.fileName = fileName
	return true
}

func (m *FileHandle) WriteString(s string) (n int, err error) {
	n, err = m.Write([]byte(s))
	return
}

func (m *FileHandle) checkReopen() (err error) {
	if m.handle == nil {
		m.handle, err = OpenFile(m.folderPath, m.fileName, m.flag, m.perm)
		return err
	} else {
		var fi os.FileInfo
		fi, err = m.handle.Stat()
		if err != nil {
			m.handle.Close()
			m.handle, err = OpenFile(m.folderPath, m.fileName, m.flag, m.perm)
			return err
		}
		if !IsExistByFileInfo(fi) {
			m.handle.Close()
			m.handle, err = OpenFile(m.folderPath, m.fileName, m.flag, m.perm)
			return err
		}
	}
	return nil
}

func (m *FileHandle) Write(b []byte) (n int, err error) {
	if err = m.checkReopen(); err != nil {
		return 0, err
	}
	n, err = m.handle.Write(b)
	return
}

func (m *FileHandle) Close() {
	if m.handle == nil {
		return
	}
	m.handle.Close()
	m.handle = nil
}

func OpenFileHandle(folderPath string, fileName string, flag int, perm os.FileMode) (*FileHandle, error) {
	var err error
	if folderPath == "" {
		err = NewError("OpenHandleFile Error: folderPath is nil")
		return nil, err
	}
	if fileName == "" {
		err = NewError("OpenHandleFile Error: fileName is nil")
		return nil, err
	}
	hf := &FileHandle{folderPath: folderPath, fileName: fileName, flag: flag, perm: perm}
	hf.handle, err = OpenFile(folderPath, fileName, flag, perm)
	if err != nil {
		return nil, err
	}
	return hf, nil
}
