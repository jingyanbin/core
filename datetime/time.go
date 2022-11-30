package datetime

import (
	tz "github.com/jingyanbin/core/timezone"
	_ "unsafe"
)

// Unix = internal.Unix 秒级时间戳
func Unix() int64

// UnixMs = internal.UnixMs 毫秒时间戳
func UnixMs() int64

// UnixNano = internal.UnixNano 纳秒时间戳
func UnixNano() int64

// IsLeapYear = internal.IsLeapYear 是否闰年
func IsLeapYear(year int) bool

// IsDoubleMonth = internal.IsDoubleMonth 是否双月
func IsDoubleMonth(month int) bool

// MonthDayNumber = internal.MonthDayNumber 这个月多少天
func MonthDayNumber(year, month int) int

// DateClockToFormat = internal.DateClockToFormat 格式化日期时间字符串
func DateClockToFormat(year, month, day, hour, min, sec int, formatter string) string

// DateClockToYmdHMS = internal.DateClockToYmdHMS 格式化为标准日期时间字符串
func DateClockToYmdHMS(year, month, day, hour, min, sec int) string

// DateClockToUnix = internal.DateClockToUnix 日期时间转秒级时间戳
func DateClockToUnix(year, month, day, hour, min, sec int, zone tz.TimeZone) (unix int64, yDay int, daySecond int, err error)

// FormatToUnix = internal.FormatToUnix 格式化字符串转秒级时间戳
func FormatToUnix(s, formatter string, zone tz.TimeZone, extend bool) (unix int64, err error)

// YmdHMSToUnix = internal.YmdHMSToUnix 标准日期时间字符串 转秒级时间戳
func YmdHMSToUnix(s string, zone tz.TimeZone, extend bool) (unix int64, err error)

// UnixToDateClock = internal.UnixToDateClock 秒级时间戳转日期时间
func UnixToDateClock(unix int64, zone tz.TimeZone) (year, month, day, hour, min, sec, yDay, daySecond int)

// FormatToDateClock = internal.FormatToDateClock 格式化日期时间字符串转日期时间
func FormatToDateClock(s, formatter string, extend bool) (year, month, day, hour, min, sec int, err error)

// UnixWeekdayA = internal.UnixWeekdayA 获得周几1-7
func UnixWeekdayA(unix int64, zone tz.TimeZone) (week int)

// UnixWeekdayB = internal.UnixWeekdayB 获得周几0-6
func UnixWeekdayB(unix int64, zone tz.TimeZone) (week int)

// UnixYearWeekNumAByISO = internal.UnixYearWeekNumAByISO 获得年周 0-53
func UnixYearWeekNumAByISO(unix int64, zone tz.TimeZone) (year int, wkn int)

// UnixDayNumber = internal.UnixDayNumber 1970年以来过来多少天
func UnixDayNumber(unix int64, zone tz.TimeZone) int64

// UnixPreClock = internal.UnixPreClock 上一个时间点
func UnixPreClock(unix int64, hour, min, sec int, zone tz.TimeZone) (int64, error)

// UnixNextClock = internal.UnixNextClock 下一个时间点
func UnixNextClock(unix int64, hour, min, sec int, zone tz.TimeZone) (int64, error)

// UnixNextWeekDayA = internal.UnixNextWeekDayA 下一周的周几1-7
func UnixNextWeekDayA(unix int64, week int, hour, min, sec int, zone tz.TimeZone) (int64, error)

// UnixNextWeekDayB = internal.UnixNextWeekDayB 下一周的周几0-6
func UnixNextWeekDayB(unix int64, week, hour, min, sec int, zone tz.TimeZone) (int64, error)

// UnixFutureWeekDayA = internal.UnixFutureWeekDayA 下一个最近的周几1-7
func UnixFutureWeekDayA(unix int64, week, hour, min, sec int, zone tz.TimeZone) (int64, error)

// UnixFutureWeekDayB = internal.UnixFutureWeekDayB 下一个最近的周几0-6
func UnixFutureWeekDayB(unix int64, week, hour, min, sec int, zone tz.TimeZone) (int64, error)

// UnixWeekStartTimeA = internal.UnixWeekStartTimeA 一周的开始时间 周1开始
func UnixWeekStartTimeA(unix int64, zone tz.TimeZone) int64

// UnixWeekEndTimeA = internal.UnixWeekEndTimeA 一周的结束时间 周日结束
func UnixWeekEndTimeA(unix int64, zone tz.TimeZone) int64

// UnixWeekDayStartTimeA = internal.UnixWeekDayStartTimeA  周几的开始时间1-7
func UnixWeekDayStartTimeA(unix int64, week int, zone tz.TimeZone) (int64, error)

// UnixWeekAnyTimeA = internal.UnixWeekAnyTimeA 周几的任何时间1-7
func UnixWeekAnyTimeA(unix int64, week, hour, min, sec int, zone tz.TimeZone) (int64, error)

// UnixYearAnyTime = internal.UnixYearAnyTime 此年中的任何时间点
func UnixYearAnyTime(unix int64, month, day, hour, min, sec int, zone tz.TimeZone) (int64, error)

// UnixMonthAnyTime = internal.UnixMonthAnyTime 此月中的任何时间点
func UnixMonthAnyTime(unix int64, day, hour, min, sec int, zone tz.TimeZone) (int64, error)

// UnixDayAnyTime = internal.UnixDayAnyTime 此日中的任何时间点
func UnixDayAnyTime(unix int64, hour, min, sec int, zone tz.TimeZone) (int64, error)

// UnixYearStartTime = internal.UnixYearStartTime 此年的开始时间
func UnixYearStartTime(unix int64, zone tz.TimeZone) int64

// UnixYearEndTime = internal.UnixYearEndTime 此年的结束时间
func UnixYearEndTime(unix int64, zone tz.TimeZone) int64

// UnixMonthStartTime = internal.UnixMonthStartTime 此月的开始时间
func UnixMonthStartTime(unix int64, zone tz.TimeZone) int64

// UnixMonthEndTime = internal.UnixMonthEndTime 此月的结束时间
func UnixMonthEndTime(unix int64, zone tz.TimeZone) int64

// UnixDayStartTime = internal.UnixDayStartTime 此日开始时间
func UnixDayStartTime(unix int64, zone tz.TimeZone) int64

// UnixDayEndTime = internal.UnixDayEndTime 此日结束时间
func UnixDayEndTime(unix int64, zone tz.TimeZone) int64

// UnixMinStartTime = internal.UnixMinStartTime 此分的开始时间
func UnixMinStartTime(unix int64) int64

// UnixDayStartTimeNext = internal.UnixDayStartTimeNext 今日开始时间算下一个时间点
func UnixDayStartTimeNext(unix int64, days, hour, min, sec int, zone tz.TimeZone) (int64, error)

// UnixHourStartTimeNext = internal.UnixHourStartTimeNext 本小时开始算下一个是几点
func UnixHourStartTimeNext(unix int64, days, hour, min, sec int, zone tz.TimeZone) (int64, error)

// UnixBetweenWeekTimeA = internal.UnixBetweenWeekTimeA 时间是否在两个星期几之间 1-7
func UnixBetweenWeekTimeA(unix int64, startTime, endTime string, zone tz.TimeZone) (bool, error)

// UnixBetweenMonthTime = internal.UnixBetweenMonthTime 时间是否在两个日期之间
func UnixBetweenMonthTime(unix int64, startTime string, endTime string, zone tz.TimeZone) (bool, error)

// UnixBetweenHourTime = internal.UnixBetweenHourTime 时间是否在两个小时之间
func UnixBetweenHourTime(unix int64, startTime string, endTime string, zone tz.TimeZone) (bool, error)

// UnixBetweenHourTimeEx = internal.UnixBetweenHourTimeEx 时间是否在两个小数之间 返回具体范围
func UnixBetweenHourTimeEx(unix int64, startHms string, endHms string, zone tz.TimeZone) (bool, int64, int64, error)

// HmsToClock = internal.HmsToClock 字符串转时分秒
func HmsToClock(hms string) (hour, min, sec int, err error)

// HmsToSec = internal.HmsToSec 字符串转秒
func HmsToSec(hms string) (int64, error)

// ClockToSec = internal.ClockToSec 时分秒 转秒
func ClockToSec(hour, min, sec int) (int64, error)

// UnixToFormat = internal.UnixToFormat 秒级时间戳格式化字符串
func UnixToFormat(unix int64, zone tz.TimeZone, formatter string) string

// UnixToYmdHMS = internal.UnixToYmdHMS 秒级时间戳格式化标准日期时间格式字符串
func UnixToYmdHMS(unix int64, zone tz.TimeZone) string

// UnixSameDay = internal.UnixSameDay 是否同一天
func UnixSameDay(unix1, unix2, offset int64, zone tz.TimeZone) bool

// UnixSameWeekA = internal.UnixSameWeekA 是否同一周 1-7
func UnixSameWeekA(unix1, unix2, offset int64, zone tz.TimeZone) bool

// UnixSameMonth = internal.UnixSameMonth 是否同意月
func UnixSameMonth(unix1, unix2, offset int64, zone tz.TimeZone) bool

// UnixDiffMonth = internal.UnixDiffMonth 俩时间相差约
func UnixDiffMonth(start, end int64, zone tz.TimeZone) int

// SeasonId = internal.SeasonId 阶段id
func SeasonId(startUnix int64, unix int64, monthNum, offset int, zone tz.TimeZone) int

// SeasonDateTime = internal.SeasonDateTime 阶段时间范围
func SeasonDateTime(startUnix int64, seasonId, monthNum, offset int, firstStartTime bool, zone tz.TimeZone) (begin *DateTime, end *DateTime)
