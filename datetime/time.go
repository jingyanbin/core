package datetime

import (
	"github.com/jingyanbin/core/internal"
	"strconv"
	"strings"
	_ "unsafe"
)

const FormatterYmdHMS = "%Y-%m-%d %H:%M:%S" //自定义例子:: "%Y/%m/%d %H:%M:%S", "%Y-%m-%d %H:%M:%S", "%Y%m%d%H%M%S"
const MinSec = 60
const HourSec = 3600
const DaySec = 3600 * 24 //每天的秒数
const WeekSec = 3600 * 24 * 7

const FirstYears = 365
const SecondYears = 365 + 365
const ThirdYears = 365 + 365 + 366
const FourYears = 365 + 365 + 366 + 365 //每个四年的总天数

var norMonth = [12]int{31, 28, 31, 30, 31, 30, 31, 31, 30, 31, 30, 31}  //平年
var leapMonth = [12]int{31, 29, 31, 30, 31, 30, 31, 31, 30, 31, 30, 31} //闰年

//go:linkname timeNow time.now
func timeNow() (sec int64, nsec int32)

// Unix 返回秒级时间戳
//
//	@Description: 秒级时间戳
//	@return int64
//

//go:noinline
func Unix() int64 {
	sec, _ := timeNow()
	return sec
}

// UnixMs 返回毫秒级时间戳
//
//	@Description: 毫秒级时间戳
//	@return int64
//

//go:noinline
func UnixMs() int64 {
	sec, nsec := timeNow()
	return sec*1000 + int64(nsec/1000000)
}

// UnixNano 返回纳秒级时间戳
//
//	@Description: 纳秒级时间戳
//	@return int64
//
//go:noinline
func UnixNano() int64 {
	sec, nsec := timeNow()
	return sec*1e9 + int64(nsec)
}

// IsLeapYear
//
//	@Description: 是否是闰年
//	@param year 年
//	@return bool
//

func IsLeapYear(year int) bool {
	return (year%4 == 0 && year%100 != 0) || year%400 == 0
}

// IsDoubleMonth
//
//	@Description: 是否双月
//	@param month
//	@return bool
//

func IsDoubleMonth(month int) bool {
	return month%2 == 0
}

// MonthDayNumber
//
//	@Description: 根据年月获得当月的天数
//	@param year年
//	@param month 月
//	@return int 天数
//

func MonthDayNumber(year, month int) int {
	if month < 1 || month > 12 || year < 1 || year > 9999 {
		return 0
	}
	if IsLeapYear(year) {
		return leapMonth[month-1]
	} else {
		return norMonth[month-1]
	}
}

func offsetDay(offset int64) int64 {
	if offset > DaySec {
		offset = DaySec
	} else if offset < 0 {
		offset = 0
	}
	return offset
}

func checkClock(hour, min, sec int) error {
	if hour > 23 || hour < 0 {
		return internal.NewError("check date clock error: out of range hour=%v", hour)
	}
	if min < 0 || min > 59 {
		return internal.NewError("check date clock error: out of range min=%v", min)
	}
	if sec < 0 || sec > 59 {
		return internal.NewError("check date clock error: out of range sec=%v", sec)
	}
	return nil
}

func checkDateClock(year, month, day, hour, min, sec int) error {
	if year > 9999 || year < 1 {
		return internal.NewError("check date clock error: out of range year=%v", year)
	}
	if month > 12 || month < 1 {
		return internal.NewError("check date clock error: out of range month=%v", month)
	}
	if day < 1 || day > 31 {
		return internal.NewError("check date clock error: out of range day=%v", day)
	}

	m := month % 7
	if m == 0 {
		m = 1
	}
	if month == 2 {
		if IsLeapYear(year) {
			if day > 29 {
				return internal.NewError("check date clock error: out of range day=%v", day)
			}
		} else {
			if day > 28 {
				return internal.NewError("check date clock error: out of range day=%v", day)
			}
		}
	} else {
		if m%2 == 1 {
			if day > 31 {
				return internal.NewError("check date clock error: out of range day=%v", day)
			}
		} else {
			if day > 30 {
				return internal.NewError("check date clock error: out of range day=%v", day)
			}
		}
	}
	return checkClock(hour, min, sec)
}

// DateClockToFormat
//
//	@Description: 年,月,日,时,分,秒 转换为日期时间字符串
//	@param year 年
//	@param month 月
//	@param day 日
//	@param hour 时
//	@param min 分
//	@param sec 秒
//	@param formatter 格式化模板 internal.FormatterYmdHMS
//	@return string 日期时间字符串
//

func DateClockToFormat(year, month, day, hour, min, sec int, formatter string) string {
	theTime := make([]byte, 0, 19)
	length := len(formatter)
	for i := 0; i < length; {
		c := formatter[i]
		if c == '%' {
			if i+1 == length {
				break
			}
			c2 := formatter[i+1]
			switch c2 {
			case 'Y': //四位数的年份表示（0000-9999）
				internal.ItoAW(&theTime, year, 4)
			case 'y': //两位数的年份表示（00-99）
				internal.ItoAW(&theTime, year, 2)
			case 'm': //月份（01-12）
				internal.ItoAW(&theTime, month, 2)
			case 'd': //月内中的一天（0-31）
				internal.ItoAW(&theTime, day, 2)
			case 'H': //24小时制小时数（0-23）
				internal.ItoAW(&theTime, hour, 2)
			case 'M': //分钟数（00=59）
				internal.ItoAW(&theTime, min, 2)
			case 'S': //秒（00-59）
				internal.ItoAW(&theTime, sec, 2)
			default:
				theTime = append(theTime, c2)
			}
			i += 2
		} else {
			theTime = append(theTime, c)
			i += 1
		}
	}
	return string(theTime)
}

// DateClockToYmdHMS
//
//	@Description: 年,月,日,时,分,秒 转换为标准日期时间字符串
//	@param year 年
//	@param month 月
//	@param day 日
//	@param hour 时
//	@param min 分
//	@param sec 秒
//	@return string 日期时间字符串
//

func DateClockToYmdHMS(year, month, day, hour, min, sec int) string {
	return DateClockToFormat(year, month, day, hour, min, sec, FormatterYmdHMS)
}

// DateClockToUnix
//
//	@Description: 年,月,日,时,分,秒 转换为秒级时间戳
//	@param year 年
//	@param month 月
//	@param day 日
//	@param hour 时
//	@param min 分
//	@param sec 秒
//	@param zone 时区
//	@return unix 秒级时间戳
//	@return yDay 一年中第几天
//	@return daySecond 一年中第几秒
//	@return err 错误
//

func DateClockToUnix(year, month, day, hour, min, sec int, zone TimeZone) (unix int64, yDay int, daySecond int, err error) {
	err = checkDateClock(year, month, day, hour, min, sec)
	if err != nil {
		err = internal.NewError("date clock to unix error check failed: %v, %v", DateClockToYmdHMS(year, month, day, hour, min, sec), err)
		return
	}

	nCha := year - 1970
	var neg bool
	if nCha < 0 {
		nCha = -nCha
		neg = true
	}

	nYear4 := nCha >> 2
	nYearMod := nCha % 4
	nDays := nYear4 * FourYears
	pMonth := &norMonth
	if nYearMod == 1 {
		nDays += FirstYears
	} else if nYearMod == 2 {
		nDays += SecondYears
		if !neg {
			pMonth = &leapMonth
		}
	} else if nYearMod == 3 {
		nDays += ThirdYears
	}

	if neg {
		nDays = -nDays
	}

	//var yDay int
	for i := 0; i < month-1; i++ {
		nDays += pMonth[i]
		yDay += pMonth[i]
	}

	nDays += day - 1
	daySecond = hour*HourSec + min*MinSec + sec
	unix = int64(nDays*DaySec+daySecond) - zone.Offset()
	//return unix, yDay, daySecond, nil
	return
}

// FormatToUnix
//
//	@Description: 日期时间字符串 转换 秒级时间戳
//	@param s 日期时间字符串
//	@param formatter 格式化模板 internal.FormatterYmdHMS
//	@param zone 时区
//	@param extend 开启扩展模可识别 2020/1/1 0:1:1, 不开启只能识别 2020/01/01 00:01:01
//	@return unix 秒级时间戳
//	@return err 错误
//

func FormatToUnix(s, formatter string, zone TimeZone, extend bool) (unix int64, err error) {
	year, month, day, hour, min, sec, err := FormatToDateClock(s, formatter, extend)
	if err != nil {
		return 0, err
	}
	unix, _, _, err = DateClockToUnix(year, month, day, hour, min, sec, zone)
	if err != nil {
		return 0, err
	}
	return unix, nil
}

// YmdHMSToUnix
//
//	@Description: 标准日期时间字符串 转换为秒级时间戳
//	@param s 标准日期时间字符串
//	@param zone 时区
//	@param extend 开启扩展模可识别 2020/1/1 0:1:1, 不开启只能识别 2020/01/01 00:01:01
//	@return unix 秒级时间戳
//	@return err 错误
//

func YmdHMSToUnix(s string, zone TimeZone, extend bool) (unix int64, err error) {
	return FormatToUnix(s, FormatterYmdHMS, zone, extend)
}

// UnixToDateClock
//
//	@Description: 秒级时间戳转换为 年,月,日,时,分,秒
//	@param unix 秒级时间戳
//	@param zone 时区
//	@return year 年
//	@return month 月
//	@return day 日
//	@return hour 时
//	@return min 分
//	@return sec 秒
//	@return yDay 一年中第几天
//	@return daySecond 一年中第几秒
//

func UnixToDateClock(unix int64, zone TimeZone) (year, month, day, hour, min, sec, yDay, daySecond int) {
	unixLocal := unix + zone.Offset()
	var nRemain int
	if unixLocal < 0 {
		nUnixSec := -unixLocal
		nDays := int(nUnixSec / DaySec)
		daySecond = (DaySec - int(nUnixSec-int64(nDays*DaySec))) % DaySec
		nYear4 := nDays/FourYears + 1
		nRemain = nYear4*FourYears - nDays
		if daySecond == 0 {
			nRemain += 1
		}
		year = 1970 - nYear4<<2
	} else {
		nDays := int(unixLocal / DaySec)
		daySecond = int(unixLocal - int64(nDays*DaySec))
		nYear4 := nDays / FourYears
		nRemain = nDays - nYear4*FourYears + 1
		year = 1970 + nYear4<<2
	}
	pMonth := &norMonth
	if nRemain <= FirstYears {

	} else if nRemain <= SecondYears {
		year += 1
		nRemain -= FirstYears
	} else if nRemain <= ThirdYears {
		year += 2
		nRemain -= SecondYears
		pMonth = &leapMonth
	} else if nRemain <= FourYears {
		year += 3
		nRemain -= ThirdYears
	} else {
		year += 4
		nRemain -= FourYears
	}
	yDay = nRemain
	var nTemp int
	for i := 0; i < 12; i++ {
		nTemp = nRemain - pMonth[i]
		if nTemp < 1 {
			month = i + 1
			if nTemp == 0 {
				day = pMonth[i]
			} else {
				day = nRemain
			}
			break
		}
		nRemain = nTemp
	}
	hour = daySecond / HourSec
	inHourSec := daySecond - hour*HourSec
	min = inHourSec / MinSec
	sec = inHourSec - min*MinSec
	return
}

// FormatToDateClock
//
//	@Description: 日期时间字符串 转换为 年,月,日,时,分,秒
//	@param s 日期时间字符串
//	@param formatter 格式化模板 internal.FormatterYmdHMS
//	@param extend 开启扩展模可识别 2020/1/1 0:1:1, 不开启只能识别 2020/01/01 00:01:01
//	@return year 年
//	@return month 月
//	@return day 日
//	@return hour 时
//	@return min 分
//	@return sec 秒
//	@return err 错误
//

func FormatToDateClock(s, formatter string, extend bool) (year, month, day, hour, min, sec int, err error) {
	defer internal.Exception(func(stack string, e error) {
		err = internal.NewError("format to date clock error exception: %v, %v \n%v", s, formatter, stack)
	})
	numbers := internal.NewNextNumber(s)
	var found bool

	length := len(formatter)
	var jump int
	for i := 0; i < length; {
		c := formatter[i]
		if c == '%' {
			if i+1 == length {
				break
			}
			c2 := formatter[i+1]
			switch c2 {
			case 'Y': //四位数的年份表示（0000-9999）
				year, found = numbers.Next(jump, 4)
			case 'm': //月份（01-12）
				month, found = numbers.Next(jump, 2)
			case 'd': //月内中的一天（01-31）
				day, found = numbers.Next(jump, 2)
			case 'H': //24小时制小时数（00-23）
				hour, found = numbers.Next(jump, 2)
			case 'M': //分钟数（00=59）
				min, found = numbers.Next(jump, 2)
			case 'S': //秒（00-59）
				sec, found = numbers.Next(jump, 2)
			default:
				err = internal.NewError("format to date clock error char: %v, %v, %v", c, c2, formatter)
				return
			}

			if !found {
				err = internal.NewError("format to date clock error time: %v, %v", s, formatter)
				return
			}
			jump = 0
			i += 2
		} else {
			if !extend {
				byt, fnd := numbers.ByteByOffset(jump)
				if (fnd && byt != c) || !fnd {
					err = internal.NewError("format to date clock error separator: %v, %v", s, formatter)
					return
				}
			}
			jump += 1
			i += 1
		}
	}

	err = checkDateClock(year, month, day, hour, min, sec)
	if err != nil {
		err = internal.NewError("format to date clock error check failed: %v, %v", s, err)
	}
	return
}

// UnixWeekdayA
//
//	@Description: 返回时间戳所在时间是周几, 星期一,为一周的开始
//	@param unix 秒级时间戳
//	@param zone 时区
//	@return week 星期 1-7
//

func UnixWeekdayA(unix int64, zone TimeZone) (week int) {
	unixLocal := unix + zone.Offset()
	if unixLocal < 0 {
		nSecond := int(unixLocal%WeekSec+WeekSec) % WeekSec
		week = nSecond/DaySec + 4
	} else {
		week = int(unixLocal%WeekSec/DaySec + 4)
	}
	if week > 7 {
		week = week - 7
	}
	return
}

// UnixWeekdayB
//
//	@Description: 返回时间戳所在的时间是周几, 星期天,为一周的开始
//	@param unix 秒级时间戳
//	@param zone 时区
//	@return week 星期 0-6
//

func UnixWeekdayB(unix int64, zone TimeZone) (week int) {
	week = UnixWeekdayA(unix, zone)
	if week == 7 {
		week = 0
	}
	return
}

// UnixYearWeekNumAByISO
//
//	@Description: 返回时间戳是在一年的第几周, 星期1为一周的开始
//	@param unix 秒级时间差
//	@param zone 时区
//	@return year 年
//	@return wkn 第几周 1-53
//

func UnixYearWeekNumAByISO(unix int64, zone TimeZone) (year int, wkn int) {
	unixWeek4, _ := UnixWeekDayStartTimeA(unix, 4, zone)
	yearStart := UnixYearStartTime(unix, zone)
	yearEnd := UnixYearEndTime(unix, zone)
	dt := UnixToDateTime(unixWeek4, zone)
	//本周四在下一年
	if unixWeek4 > yearEnd {
		return dt.Year(), 1
	} else if unixWeek4 < yearStart { //本周4在上一年
		unixWeek1 := unixWeek4 - 3*DaySec
		nSecond := int(unix - unixWeek1)
		return dt.Year(), (nSecond / WeekSec) + 1
	} else {
		unixYearStartWeek1, _ := UnixWeekDayStartTimeA(yearStart, 1, zone)
		nSecond := int(unix - unixYearStartWeek1)
		return dt.Year(), (nSecond / WeekSec) + 1
	}
}

// UnixDayNumber
//
//	@Description: 返回时间戳1970年1月1日以来的天数
//	@param unix 秒级时间戳
//	@param zone 时区
//	@return int64 天数
//

func UnixDayNumber(unix int64, zone TimeZone) int64 {
	return (unix + zone.Offset()) / int64(DaySec)
}

// UnixPreClock
//
//	@Description: 获得上一个时间点的时间戳
//	@param unix 秒级时间戳
//	@param hour 时
//	@param min 分
//	@param sec 秒
//	@param zone 时区
//	@return int64 秒级时间戳
//	@return error 错误
//

func UnixPreClock(unix int64, hour, min, sec int, zone TimeZone) (int64, error) {
	nextUnix, err := UnixDayStartTimeNext(unix, 0, hour, min, sec, zone)
	if err != nil {
		return 0, err
	}
	if unix < nextUnix {
		return nextUnix - DaySec, nil
	} else {
		return nextUnix, nil
	}
}

// UnixNextClock
//
//	@Description: 获得下一个时间点的时间戳
//	@param unix 秒级时间戳
//	@param hour 时
//	@param min 分
//	@param sec 秒
//	@param zone 时区
//	@return int64 秒级时间戳
//	@return error 错误
//

func UnixNextClock(unix int64, hour, min, sec int, zone TimeZone) (int64, error) {
	nextUnix, err := UnixDayStartTimeNext(unix, 0, hour, min, sec, zone)
	if err != nil {
		return 0, err
	}
	if unix < nextUnix {
		return nextUnix, nil
	} else {
		return nextUnix + DaySec, nil
	}
}

// UnixNextWeekDayA
//
//	@Description: 返回时间戳下一周的星期几的秒级时间戳(星期1为周的开始)
//	@param unix 秒级时间戳
//	@param week 周几 1-7
//	@param hour 时
//	@param min 分
//	@param sec 秒
//	@param zone 时区
//	@return int64 秒级时间差
//	@return error 错误
//

func UnixNextWeekDayA(unix int64, week int, hour, min, sec int, zone TimeZone) (int64, error) {
	if week < 1 || week > 7 {
		return 0, internal.NewError("week out of range(1,7): %v", week)
	}
	w := UnixWeekdayA(unix, zone)
	days := week - w
	return UnixDayStartTime(unix, zone) + int64(WeekSec+days*DaySec+hour*HourSec+min*MinSec+sec), nil
}

// UnixNextWeekDayB
//
//	@Description: 返回时间戳下一周的星期几的秒级时间戳(星期天为周的开始)
//	@param unix 秒级时间戳
//	@param week 周几 0-6
//	@param hour 时
//	@param min 分
//	@param sec 秒
//	@param zone 时区
//	@return int64 秒级时间戳
//	@return error 错误
//

func UnixNextWeekDayB(unix int64, week, hour, min, sec int, zone TimeZone) (int64, error) {
	if week < 0 || week > 6 {
		return 0, internal.NewError("week out of range(0, 6): %v", week)
	}
	w := UnixWeekdayB(unix, zone)
	days := week - w
	return UnixDayStartTime(unix, zone) + int64(WeekSec+days*DaySec+hour*HourSec+min*MinSec+sec), nil
}

// UnixFutureWeekDayA
//
//	@Description: 返回时间戳下一个最近的星期几的秒级时间戳(星期1为周的开始)
//	@param unix 秒级时间戳
//	@param week 周几 1-7
//	@param hour 时
//	@param min 分
//	@param sec 秒
//	@param zone 时区
//	@return int64 秒级时间戳
//	@return error 错误
//

func UnixFutureWeekDayA(unix int64, week, hour, min, sec int, zone TimeZone) (int64, error) {
	if week < 1 || week > 7 {
		return 0, internal.NewError("week out of range(1,7): %v", week)
	}
	if err := checkClock(hour, min, sec); err != nil {
		return 0, err
	}
	w := UnixWeekdayA(unix, zone)
	days := week - w
	unixDaySecond := hour*HourSec + min*MinSec + sec
	if days > 0 {
		return UnixDayStartTime(unix, zone) + int64(days*DaySec+unixDaySecond), nil
	} else if days == 0 {
		_, _, _, _, _, _, _, daySecond := UnixToDateClock(unix, zone)
		if unixDaySecond > daySecond {
			return UnixDayStartTime(unix, zone) + int64(unixDaySecond), nil
		} else {
			return UnixDayStartTime(unix, zone) + int64(WeekSec+days*DaySec+unixDaySecond), nil
		}
	} else {
		return UnixDayStartTime(unix, zone) + int64(WeekSec+days*DaySec+unixDaySecond), nil
	}
}

// UnixFutureWeekDayB
//
//	@Description: 返回时间戳下一个最近的星期几的秒级时间戳(星期天为周的开始)
//	@param unix 秒级时间戳
//	@param week 周几 0-6
//	@param hour 时
//	@param min 分
//	@param sec 秒
//	@param zone 时区
//	@return int64 秒级时间戳
//	@return error 错误
//

func UnixFutureWeekDayB(unix int64, week, hour, min, sec int, zone TimeZone) (int64, error) {
	if week < 0 || week > 6 {
		return 0, internal.NewError("week out of range(0, 6): %v", week)
	}
	if err := checkClock(hour, min, sec); err != nil {
		return 0, err
	}
	w := UnixWeekdayB(unix, zone)
	days := week - w
	unixDaySecond := hour*HourSec + min*MinSec + sec
	if days > 0 {
		return UnixDayStartTime(unix, zone) + int64(days*DaySec+unixDaySecond), nil
	} else if days == 0 {
		_, _, _, _, _, _, _, daySecond := UnixToDateClock(unix, zone)
		if unixDaySecond > daySecond {
			return UnixDayStartTime(unix, zone) + int64(unixDaySecond), nil
		} else {
			return UnixDayStartTime(unix, zone) + int64(WeekSec+days*DaySec+unixDaySecond), nil
		}
	} else {
		return UnixDayStartTime(unix, zone) + int64(WeekSec+days*DaySec+unixDaySecond), nil
	}
}

// UnixWeekStartTimeA
//
//	@Description: 返回时间戳所在周的开始时间（周一的0点）
//	@param unix 秒级时间戳
//	@param zone 时区
//	@return int64 秒级时间戳
//

func UnixWeekStartTimeA(unix int64, zone TimeZone) int64 {
	return UnixDayStartTime(unix, zone) + int64(1-UnixWeekdayA(unix, zone))*DaySec
}

// UnixWeekEndTimeA
//
//	@Description: 返回时间戳所在周的结束时间（周日的23:59:59）
//	@param unix 秒级时间戳
//	@param zone 时区
//	@return int64 秒级时间戳
//

func UnixWeekEndTimeA(unix int64, zone TimeZone) int64 {
	return UnixDayStartTime(unix, zone) + int64(8-UnixWeekdayA(unix, zone))*DaySec - 1
}

// UnixWeekDayStartTimeA
//
//	@Description: 返回时间戳所在周几开始时间（周一的0点）
//	@param unix 秒级时间戳
//	@param week 周几 1-7
//	@param zone 时区
//	@return int64 秒级时间戳
//	@return error 错误
//

func UnixWeekDayStartTimeA(unix int64, week int, zone TimeZone) (int64, error) {
	if week < 1 || week > 7 {
		return 0, internal.NewError("week out of range(1, 7): %v", week)
	}
	return UnixDayStartTimeNext(unix, week-UnixWeekdayA(unix, zone), 0, 0, 0, zone)
}

// UnixWeekAnyTimeA
//
//	@Description: 返回时间戳所在周的任何时间
//	@param unix 秒级时间戳
//	@param week 周几 1-7
//	@param hour 时
//	@param min 分
//	@param sec 秒
//	@param zone 时区
//	@return int64 秒级时间戳
//	@return error 错误
//

func UnixWeekAnyTimeA(unix int64, week, hour, min, sec int, zone TimeZone) (int64, error) {
	if week < 1 || week > 7 {
		return 0, internal.NewError("week out of range(1, 7): %v", week)
	}
	return UnixDayStartTimeNext(unix, week-UnixWeekdayA(unix, zone), hour, min, sec, zone)
}

// UnixYearAnyTime
//
//	@Description: 返回时间戳所在年内某月的的任何时间
//	@param unix 秒级时间戳
//	@param month 月
//	@param day 日
//	@param hour 时
//	@param min 分
//	@param sec 秒
//	@param zone 时区
//	@return int64 秒级时间戳
//	@return error 错误
//

func UnixYearAnyTime(unix int64, month, day, hour, min, sec int, zone TimeZone) (int64, error) {
	year, _, _, _, _, _, _, _ := UnixToDateClock(unix, zone)
	if unixMonth, _, _, err := DateClockToUnix(year, month, day, hour, min, sec, zone); err != nil {
		return 0, err
	} else {
		return unixMonth, nil
	}
}

// UnixMonthAnyTime
//
//	@Description: 返回时间戳所在月的任意时间点
//	@param unix 秒级时间戳
//	@param day 日
//	@param hour 时
//	@param min 分
//	@param sec 秒
//	@param zone 时区
//	@return int64 秒级时间戳
//	@return error 错误
//

func UnixMonthAnyTime(unix int64, day, hour, min, sec int, zone TimeZone) (int64, error) {
	year, month, _, _, _, _, _, _ := UnixToDateClock(unix, zone)
	if unixMonth, _, _, err := DateClockToUnix(year, month, day, hour, min, sec, zone); err != nil {
		return 0, err
	} else {
		return unixMonth, nil
	}
}

// UnixDayAnyTime
//
//	@Description: 根据时间戳得到当天的任意时间
//	@param unix 秒级时间戳
//	@param hour 时
//	@param min 分
//	@param sec 秒
//	@param zone 时区
//	@return int64 秒级时间戳
//	@return error 错误
//

func UnixDayAnyTime(unix int64, hour, min, sec int, zone TimeZone) (int64, error) {
	if secNum, err := ClockToSec(hour, min, sec); err != nil {
		return 0, err
	} else {
		return UnixDayStartTime(unix, zone) + secNum, nil
	}
}

// UnixYearStartTime
//
//	@Description: 返回时间戳所在年的开始时间
//	@param unix 秒级时间戳
//	@param zone 时区
//	@return int64 秒级时间戳
//

func UnixYearStartTime(unix int64, zone TimeZone) int64 {
	_, _, _, _, _, _, yDay, daySecond := UnixToDateClock(unix, zone)
	unixLocal := unix + zone.Offset()
	return unixLocal - int64((yDay-1)*DaySec+daySecond) - zone.Offset()
	//yearStartUnix, _ := UnixYearAnyTime(unix, 1, 1, 0, 0, 0, zone)
	//return yearStartUnix
}

// UnixYearEndTime
//
//	@Description: 返回时间戳所在年的结束时间
//	@param unix 秒级时间戳
//	@param zone 时区
//	@return int64 秒级时间戳
//

func UnixYearEndTime(unix int64, zone TimeZone) int64 {
	year, _, _, _, _, _, _, _ := UnixToDateClock(unix, zone)
	nextYearStartUnix, _, _, _ := DateClockToUnix(year+1, 1, 1, 0, 0, 0, zone)
	return nextYearStartUnix - 1
}

// UnixMonthStartTime
//
//	@Description: 返回时间戳所在月1日0时的秒级时间戳
//	@param unix 秒级时间戳
//	@param zone 时区
//	@return int64 秒级时间戳
//

func UnixMonthStartTime(unix int64, zone TimeZone) int64 {
	year, month, _, _, _, _, _, _ := UnixToDateClock(unix, zone)
	unixMonth, _, _, _ := DateClockToUnix(year, month, 1, 0, 0, 0, zone)
	return unixMonth
}

// UnixMonthEndTime
//
//	@Description: 返回时间戳所在月的最后一天23:59:59
//	@param unix 秒级时间戳
//	@param zone 时区
//	@return int64 秒级时间戳
//

func UnixMonthEndTime(unix int64, zone TimeZone) int64 {
	year, month, _, _, _, _, _, _ := UnixToDateClock(unix, zone)
	if month < 12 {
		month += 1
	} else {
		month = 1
		year += 1
	}
	unixMonth, _, _, _ := DateClockToUnix(year, month, 1, 0, 0, 0, zone)
	return unixMonth - 1
}

// UnixDayStartTime
//
//	@Description: 返回时间戳当天0时的秒级时间戳
//	@param unix 秒级时间戳
//	@param zone 时区
//	@return int64 秒级时间戳
//

func UnixDayStartTime(unix int64, zone TimeZone) int64 {
	unixLocal := unix + zone.Offset()
	if unixLocal < 0 {
		nSecond := (unixLocal%DaySec + DaySec) % DaySec
		return unixLocal - nSecond - zone.Offset()
	}
	return unixLocal - unixLocal%DaySec - zone.Offset()
}

// UnixDayEndTime
//
//	@Description: 返回时间戳当天23:59:59
//	@param unix 秒级时间戳
//	@param zone 时区
//	@return int64 秒级时间戳
//

func UnixDayEndTime(unix int64, zone TimeZone) int64 {
	return UnixDayStartTime(unix, zone) + DaySec - 1
}

// UnixHourStartTime
//
//	@Description: 返回时间戳本小时的0分的秒级时间戳
//	@param unix 秒级时间戳
//	@return int64 秒级时间戳
//

func UnixHourStartTime(unix int64) int64 {
	if unix < 0 {
		nSecond := (unix%HourSec + HourSec) % DaySec
		return unix - nSecond
	}
	return unix - unix%HourSec
}

// UnixMinStartTime
//
//	@Description: 当前分的0秒时间戳（秒）
//	@param unix 秒级时间戳
//	@return int64 秒级时间戳
//

func UnixMinStartTime(unix int64) int64 {
	return (unix / 60) * 60
}

// UnixDayStartTimeNext
//
//	@Description: 返回时间戳当然开始时间后N天特定时间的秒级时间戳
//	@param unix 秒级时间戳
//	@param days 天数
//	@param hour 小时
//	@param min 分钟
//	@param sec 秒数
//	@param zone 时区
//	@return int64 秒级时间戳
//	@return error 错误
//

func UnixDayStartTimeNext(unix int64, days, hour, min, sec int, zone TimeZone) (int64, error) {
	if err := checkClock(hour, min, sec); err != nil {
		return 0, err
	}
	start := UnixDayStartTime(unix, zone)
	return start + int64(days*DaySec+hour*HourSec+min*MinSec+sec), nil
}

// UnixHourStartTimeNext
//
//	@Description: 返回时间戳小时开始时间后N天特定时间的秒级时间戳
//	@param unix 秒级时间戳
//	@param days 天数
//	@param hour 小时
//	@param min 分钟
//	@param sec 秒数
//	@param zone 时区
//	@return int64 时间戳
//	@return error 错误
//

func UnixHourStartTimeNext(unix int64, days, hour, min, sec int, zone TimeZone) (int64, error) {
	if err := checkClock(hour, min, sec); err != nil {
		return 0, err
	}
	start := UnixHourStartTime(unix)
	return start + int64(days*DaySec+hour*HourSec+min*MinSec+sec), nil
}

// UnixBetweenWeekTimeA
//
//	@Description: 判断时间戳是否在一段时间范围 支持跨周 eg: 7_12:00:00 => 3_13:00:00 / 7+12:00:00 => 3+13:00:00
//	@param unix 秒级时间戳
//	@param startTime 开始时间 eg: 7_12:00:00 (周几 1-7)
//	@param endTime 结束时间 3_13:00:00
//	@param zone 时区
//	@return bool
//	@return error 错误
//

func UnixBetweenWeekTimeA(unix int64, startTime, endTime string, zone TimeZone) (bool, error) {
	startNums, endNums := internal.NewNextNumber(startTime).Numbers(), internal.NewNextNumber(endTime).Numbers()
	if len(startNums) != 4 || len(endNums) != 4 {
		return false, internal.NewError("time error: %s, %s", startTime, endTime)
	}
	startUnix, err := UnixWeekAnyTimeA(unix, startNums[0], startNums[1], startNums[2], startNums[3], zone)
	if err != nil {
		return false, err
	}
	endUnix, err := UnixWeekAnyTimeA(unix, endNums[0], endNums[1], endNums[2], endNums[3], zone)
	if err != nil {
		return false, err
	}
	chaSec := endUnix - startUnix
	if chaSec >= 0 {
		if unix >= startUnix && unix <= endUnix {
			return true, nil
		}
	} else {
		weekStartUnix := UnixWeekStartTimeA(unix, zone)
		if unix >= weekStartUnix && unix <= endUnix {
			return true, nil
		}
		weekEndUnix := UnixWeekEndTimeA(unix, zone)
		if unix >= startUnix && unix <= weekEndUnix {
			return true, nil
		}
	}
	return false, nil
}

// UnixBetweenMonthTime
//
//	@Description: 判断时间戳是否在一段时间范围内 eg: 08-15 12:00:00 => 01-15 13:00:00
//	@param unix 秒级时间戳
//	@param startTime 开始时间 eg: 08-15 12:00:00
//	@param endTime 结束时间 eg: 01-15 13:00:00
//	@param zone 时区
//	@return bool 是否时间范围内
//	@return error 错误
//

func UnixBetweenMonthTime(unix int64, startTime string, endTime string, zone TimeZone) (bool, error) {
	startNums, endNums := internal.NewNextNumber(startTime).Numbers(), internal.NewNextNumber(endTime).Numbers()
	if len(startNums) != 5 || len(endNums) != 5 {
		return false, internal.NewError("time error: %s, %s", startTime, endTime)
	}
	startUnix, err := UnixYearAnyTime(unix, startNums[0], startNums[1], startNums[2], startNums[3], startNums[4], zone)
	if err != nil {
		return false, err
	}
	endUnix, err := UnixYearAnyTime(unix, endNums[0], endNums[1], endNums[2], endNums[3], endNums[4], zone)
	if err != nil {
		return false, err
	}

	chaSec := endUnix - startUnix
	if chaSec >= 0 {
		if unix >= startUnix && unix <= endUnix {
			return true, nil
		}
	} else {
		yearStartTime := UnixYearStartTime(unix, zone)
		if unix >= yearStartTime && unix <= endUnix {
			return true, nil
		}
		yearEndTime := UnixYearEndTime(unix, zone)
		if unix >= startUnix && unix <= yearEndTime {
			return true, nil
		}
	}
	return false, nil
}

// UnixBetweenHourTime
//
//	@Description: 判断时间戳是否在时间范围内 小时级别 eg: 12:00:00 => 13:00:00
//	@param unix 秒级时间戳
//	@param startTime 开始时间 eg: 12:00:00
//	@param endTime 结束时间 eg: 13:00:00
//	@param zone 时区
//	@return bool 是否在时间范围内
//	@return error 错误
//

func UnixBetweenHourTime(unix int64, startTime string, endTime string, zone TimeZone) (bool, error) {
	startNums, endNums := internal.NewNextNumber(startTime).Numbers(), internal.NewNextNumber(endTime).Numbers()
	if len(startNums) != 3 || len(endNums) != 3 {
		return false, internal.NewError("time error: %s, %s", startTime, endTime)
	}
	startUnix, err := UnixDayAnyTime(unix, startNums[0], startNums[1], startNums[2], zone)
	if err != nil {
		return false, err
	}
	endUnix, err := UnixDayAnyTime(unix, endNums[0], endNums[1], endNums[2], zone)
	if err != nil {
		return false, err
	}

	chaSec := endUnix - startUnix
	if chaSec >= 0 {
		if unix >= startUnix && unix <= endUnix {
			return true, nil
		}
	} else {
		dayStartTime := UnixDayStartTime(unix, zone)
		if unix >= dayStartTime && unix <= endUnix {
			return true, nil
		}
		dayEndTime := UnixDayEndTime(unix, zone)
		if unix >= startUnix && unix <= dayEndTime {
			return true, nil
		}
	}
	return false, nil
}

// UnixBetweenHourTimeEx
//
//	@Description: 判断时间戳是否在时间范围内 小时级别 eg: 12:00:00 => 13:00:00
//	@param unix 秒级时间戳
//	@param startHms 开始时间 eg: 12:00:00
//	@param endHms 结束时间 eg: 13:00:00
//	@param zone 时区
//	@return bool 在时间范围内
//	@return int64 开始时间戳
//	@return int64 结束时间戳
//	@return error 错误
//

func UnixBetweenHourTimeEx(unix int64, startHms string, endHms string, zone TimeZone) (bool, int64, int64, error) {
	startNums, endNums := internal.NewNextNumber(startHms).Numbers(), internal.NewNextNumber(endHms).Numbers()
	if len(startNums) != 3 || len(endNums) != 3 {
		return false, 0, 0, internal.NewError("time error: %s, %s", startHms, endHms)
	}
	startUnix, err := UnixDayAnyTime(unix, startNums[0], startNums[1], startNums[2], zone)
	if err != nil {
		return false, 0, 0, err
	}
	endUnix, err := UnixDayAnyTime(unix, endNums[0], endNums[1], endNums[2], zone)
	if err != nil {
		return false, 0, 0, err
	}
	chaSec := endUnix - startUnix
	if chaSec >= 0 {
		return unix >= startUnix && unix <= endUnix, startUnix, endUnix, nil
	} else {
		dayStartTime := UnixDayStartTime(unix, zone)
		if unix >= dayStartTime && unix <= endUnix {
			start, err := UnixDayStartTimeNext(unix, -1, startNums[0], startNums[1], startNums[2], zone)
			if err != nil {
				return false, 0, 0, err
			}
			return true, start, endUnix, nil
		}
		dayEndTime := UnixDayEndTime(unix, zone)
		if unix >= startUnix && unix <= dayEndTime {
			end, err := UnixDayStartTimeNext(unix, 1, endNums[0], endNums[1], endNums[2], zone)
			if err != nil {
				return false, 0, 0, err
			}
			return true, startUnix, end, nil
		}
		nextStart, err := UnixDayStartTimeNext(unix, 0, startNums[0], startNums[1], startNums[2], zone)
		if err != nil {
			return false, 0, 0, err
		}
		nextEnd, err := UnixDayStartTimeNext(unix, 0, endNums[0], endNums[1], endNums[2], zone)
		if err != nil {
			return false, 0, 0, err
		}
		return false, nextStart, nextEnd, nil
	}
}

var hmsLenErr = internal.NewError("hms string len error")

// HmsToClock
//
//	@Description: 时间字符串 转换为时分秒
//	@param hms 时间字符串
//	@return hour 时
//	@return min 分
//	@return sec 秒
//	@return err 错误
//

func HmsToClock(hms string) (hour, min, sec int, err error) {
	clocks := strings.Split(hms, ":")
	if len(clocks) < 3 {
		return 0, 0, 0, hmsLenErr
	}
	if hour, err = strconv.Atoi(clocks[0]); err != nil {
		return 0, 0, 0, err
	}
	if min, err = strconv.Atoi(clocks[1]); err != nil {
		return 0, 0, 0, err
	}
	if sec, err = strconv.Atoi(clocks[2]); err != nil {
		return 0, 0, 0, err
	}
	return
}

// HmsToSec
//
//	@Description: 时间转换为clock
//	@param hms 时间字符串
//	@return int64 返回秒数
//	@return error 返回错误
//

func HmsToSec(hms string) (int64, error) {
	hour, min, sec, err := HmsToClock(hms)
	if err != nil {
		return 0, err
	}

	return int64(hour*HourSec + min*MinSec + sec), nil
}

// ClockToSec
//
//	@Description: 时分秒转换为总秒数
//	@param hour 时
//	@param min 分
//	@param sec 秒
//	@return int64 返回秒数
//	@return error 错误
//

func ClockToSec(hour, min, sec int) (int64, error) {
	if err := checkClock(hour, min, sec); err != nil {
		return 0, err
	}
	return int64(hour*HourSec + min*MinSec + sec), nil
}

// UnixToFormat
//
//	@Description: 返回时间戳的 格式化日期时间字符串
//	@param unix 秒级时间戳
//	@param zone 时区
//	@param formatter 格式化模板
//	@return string 返回格式化后的日期时间字符串
//

func UnixToFormat(unix int64, zone TimeZone, formatter string) string {
	year, month, day, hour, min, sec, _, _ := UnixToDateClock(unix, zone)
	return DateClockToFormat(year, month, day, hour, min, sec, formatter)
}

// UnixToYmdHMS
//
//	@Description: 返回时间戳的 标准日期时间字符串
//	@param unix 秒级时间戳
//	@param zone 时区
//	@return string 返回标准日期时间字符串
//

func UnixToYmdHMS(unix int64, zone TimeZone) string {
	return UnixToFormat(unix, zone, FormatterYmdHMS)
}

// UnixSameDay
//
//	@Description: 判断时间戳是否是同一天
//	@param unix1 时间戳1
//	@param unix2 时间戳2
//	@param offset 偏移量
//	@param zone 时区
//	@return bool 是否同一天
//

func UnixSameDay(unix1, unix2, offset int64, zone TimeZone) bool {
	start := UnixDayStartTime(unix1-offset, zone)
	end := UnixDayEndTime(unix1-offset, zone)
	if unix2-offset < start {
		return false
	} else if unix2-offset > end {
		return false
	}
	return true
}

// UnixSameWeekA
//
//	@Description: 判断时间戳是否是同一周, 星期1为周的开始
//	@param unix1 时间戳1
//	@param unix2 时间戳2
//	@param offset 偏移量
//	@param zone 时区
//	@return bool 是否同一周
//

func UnixSameWeekA(unix1, unix2, offset int64, zone TimeZone) bool {
	start, end := UnixWeekStartTimeA(unix1-offset, zone), UnixWeekEndTimeA(unix1-offset, zone)
	if unix2-offset < start {
		return false
	} else if unix2-offset > end {
		return false
	}
	return true
}

// UnixSameMonth
//
//	@Description: 判断时间戳是否是同一月
//	@param unix1 时间戳1
//	@param unix2 时间戳2
//	@param offset 偏移量
//	@param zone 时区
//	@return bool 是否同一月
//

func UnixSameMonth(unix1, unix2, offset int64, zone TimeZone) bool {
	dt1, dt2 := UnixToDateTime(unix1-offset, zone), UnixToDateTime(unix2-offset, zone)
	return dt1.Year() == dt2.Year() && dt1.Month() == dt2.Month()
}

// UnixDiffMonth
//
//	@Description: 开始时间和结束时间相差月份
//	@param start 开始秒级时间戳
//	@param end 结束秒级时间戳
//	@param zone 时区
//	@return int 天数
//

func UnixDiffMonth(start, end int64, zone TimeZone) int {
	d1 := UnixToDateTime(start, zone)
	d2 := UnixToDateTime(end, zone)
	return DateTimeDiffMonth(d1, d2)
}

// SeasonId
//
//	@Description: 返回阶段ID
//	@param startUnix 阶段开始时间戳
//	@param unix 时间戳
//	@param monthNum 一个阶段多少个月
//	@param offset 偏移量 秒数 最大3600*24*28
//	@param zone 时区
//	@return int 阶段ID
func SeasonId(startUnix int64, unix int64, monthNum, offset int, zone TimeZone) int {
	if monthNum < 1 {
		monthNum = 1
	}
	startUnix -= int64(offset)
	unix -= int64(offset)
	monthCha := UnixDiffMonth(startUnix, unix, zone)
	if monthCha > 0 {
		return monthCha / monthNum
	} else if monthCha < 0 {
		return (monthCha - monthNum + 1) / monthNum
	} else {
		return 0
	}
}

// SeasonDateTime
//
//	@Description: 获得阶段时间范围
//	@param startUnix 阶段开始时间戳 单位: 秒
//	@param seasonId 阶段ID
//	@param monthNum 一个阶段多少月
//	@param offset 偏移量 秒数 最大3600*24*28
//	@param firstStartTime 第一个阶段是否从月的第一天开始, 否则从真实时间开始
//	@param zone 时区
//	@return begin 返回阶段的开始时间
//	@return end 返回阶段的结束时间
func SeasonDateTime(startUnix int64, seasonId, monthNum, offset int, firstStartTime bool, zone TimeZone) (begin *DateTime, end *DateTime) {
	start := UnixToDateTime(startUnix, zone)
	addMonthNum := seasonId * monthNum
	if seasonId == 0 && firstStartTime == false {
		begin = start
	} else {
		begin = start.MonthStartDateTime(addMonthNum).AddSec(offset)
	}
	end = begin.MonthEndDateTime(monthNum - 1).AddSec(offset)
	return begin, end
}
