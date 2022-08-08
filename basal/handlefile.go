package basal

import (
	"github.com/jingyanbin/core/internal"
	"os"
)

// OpenFile 打开文件 自动创建目录
func OpenFile(folderPath string, fileName string, flag int, perm os.FileMode) (file *os.File, err error)

// OpenFileB 打开文件 自动创建目录
func OpenFileB(filePath string, flag int, perm os.FileMode) (file *os.File, err error)

// IsExistBy 文件是否存在
func IsExistBy(f *os.File, err error) bool

// IsExist 文件或文件夹是否存在
func IsExist(path string) (bool, error)

var ErrNotFolder = internal.ErrNotFolder

// IsExistFolder 文件夹是否存在
func IsExistFolder(path string) (bool, error)

var ErrNotFile = internal.ErrNotFolder

// IsExistFile 文件是否存在
func IsExistFile(path string) (bool, error)

type HandleFile = internal.HandleFile

const HANDLE_FILE_FLAG_WRITER = internal.HANDLE_FILE_FLAG_WRITER
const HANDLE_FILE_PERM_ALL = internal.HANDLE_FILE_PERM_ALL

func NewHandleFile(flag int, perm os.FileMode) *HandleFile

func OpenHandleFile(folderPath string, fileName string, flag int, perm os.FileMode) (*HandleFile, error)

var Path = internal.Path
