//go:build linux

package basal

import (
	"os"
	"syscall"
	_ "unsafe"
)

// IsExistByFileInfo
//
//	@Description: 文件是否存在
//	@param info 文件信息
//	@return bool 是否存在
func IsExistByFileInfo(info os.FileInfo) bool {
	sys2 := info.Sys()
	t := sys2.(*syscall.Stat_t)
	return t.Nlink > 0
}
