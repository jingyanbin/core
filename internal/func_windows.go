//go:build windows

package internal

import (
	"os"
	_ "unsafe"
)

// IsExistByFileInfo
//
//	@Description: 文件是否存在
//	@param info 文件信息
//	@return bool 是否存在
//
//go:linkname IsExistByFileInfo github.com/jingyanbin/core/basal.IsExistByFileInfo
func IsExistByFileInfo(info os.FileInfo) bool {
	return true
}
