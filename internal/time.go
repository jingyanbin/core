package internal

import (
	_ "unsafe"
)

//go:linkname timeNow time.now
func timeNow() (sec int64, nsec int32)

// Unix 返回秒级时间戳
//
//go:linkname Unix github.com/jingyanbin/core/datetime.Unix
//go:noinline
func Unix() int64 {
	sec, _ := timeNow()
	return sec
}

// UnixMs 返回毫秒级时间戳
//
//go:linkname UnixMs github.com/jingyanbin/core/datetime.UnixMs
//go:noinline
func UnixMs() int64 {
	sec, nsec := timeNow()
	return sec*1000 + int64(nsec/1000000)
}

// UnixNano 返回纳秒级时间戳
//
//go:linkname UnixNano github.com/jingyanbin/core/datetime.UnixNano
//go:noinline
func UnixNano() int64 {
	sec, nsec := timeNow()
	return sec*1e9 + int64(nsec)
}
