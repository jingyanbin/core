//go:build windows

package basal

import (
	"os"
	_ "unsafe"
)

// IsExistByFileInfo
//
//	@Description: 文件是否存在
//	@param info 文件信息
//	@return bool 是否存在
func IsExistByFileInfo(info os.FileInfo) bool {
	return true
}
