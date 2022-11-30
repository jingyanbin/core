package datetime

import (
	"github.com/jingyanbin/core/internal"
	tz "github.com/jingyanbin/core/timezone"
)

type DateTime = internal.DateTime // 日期时间对象

// Now = internal.Now 当前时间对象
func Now() (dt *DateTime)

// UnixToDateTime = internal.UnixToDateTime 秒级时间戳转日期时间对象
func UnixToDateTime(unix int64, zone tz.TimeZone) (dt *DateTime)

// FormatToDateTime = internal.FormatToDateTime 日期时间字符串转日期时间对象
func FormatToDateTime(s, formatter string, zone tz.TimeZone, extend bool) (dt *DateTime, err error)

// YmdHMSToDateTime = internal.YmdHMSToDateTime 标志日期时间字符串转日期时间对象
func YmdHMSToDateTime(s string, zone tz.TimeZone, extend bool) (dt *DateTime, err error)

// DateClockToDateTime = internal.DateClockToDateTime 日期时间转日期时间对象
func DateClockToDateTime(year, month, day, hour, min, sec int, zone tz.TimeZone) (dt *DateTime, err error)

// YearMonthByAddMonthNum = internal.YearMonthByAddMonthNum 年约加上月数得到新的年月
func YearMonthByAddMonthNum(year, month, addMonthNum int) (y int, m int)

// DateTimeDiffMonth = internal.DateTimeDiffMonth 相差月份
func DateTimeDiffMonth(start, end *DateTime) int
