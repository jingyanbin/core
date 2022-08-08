package datetime

import (
	"github.com/jingyanbin/core/basal"
	internal "github.com/jingyanbin/core/internal"
	tz "github.com/jingyanbin/core/timezone"
	"strconv"
	"strings"
	_ "unsafe"
)

const formatterYmdHMS = "%Y-%m-%d %H:%M:%S"
const minSec = 60
const hourSec = 3600
const daySec = 3600 * 24 //每天的秒数
const weekSec = 3600 * 24 * 7

const firstYears = 365
const secondYears = 365 + 365
const thirdYears = 365 + 365 + 366
const fourYears = 365 + 365 + 366 + 365 //每个四年的总天数

var norMonth = [12]int{31, 28, 31, 30, 31, 30, 31, 31, 30, 31, 30, 31}  //平年
var leapMonth = [12]int{31, 29, 31, 30, 31, 30, 31, 31, 30, 31, 30, 31} //闰年

//返回秒级时间戳
func Unix() int64

//返回毫秒级时间戳
func UnixMs() int64

//返回纳秒级时间戳
func UnixNano() int64

//是否是闰年
func IsLeapYear(year int) bool {
	return (year%4 == 0 && year%100 != 0) || year%400 == 0
}

//是否双月
func IsDoubleMonth(month int) bool {
	return month%2 == 0
}

//根据年月获得当月的天数
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
	if offset > daySec {
		offset = daySec
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

//@description: 年,月,日,时,分,秒 -> 日期时间字符串
//@param:       year, month, day, hour, min, sec "年,月,日,时,分,秒"
//@param:       formatter string "格式化模板" 如: "%Y/%m/%d %H:%M:%S", "%Y-%m-%d %H:%M:%S", "%Y%m%d%H%M%S"
//@param:       string "日期时间字符串"
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

//@description: 年,月,日,时,分,秒 -> 换为标准日期时间字符串
//@param:       year, month, day, hour, min, sec "年,月,日,时,分,秒"
//@param:       string "日期时间字符串"
func DateClockToYmdHMS(year, month, day, hour, min, sec int) string {
	return DateClockToFormat(year, month, day, hour, min, sec, formatterYmdHMS)
}

//@description: 年,月,日,时,分,秒 -> 转换为秒级时间戳
//@param:       year, month, day, hour, min, sec "年,月,日,时,分,秒"
//@param:       zone tz.TimeZone "时区"
//@return:      unix int64 "秒级时间戳"
//@return:      yDay "一年中第几天"
//@return:      daySecond "一天中第几秒"
//@return:      error "错误信息"
func DateClockToUnix(year, month, day, hour, min, sec int, zone tz.TimeZone) (unix int64, yDay int, daySecond int, err error) {
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
	nDays := nYear4 * fourYears
	pMonth := &norMonth
	if nYearMod == 1 {
		nDays += firstYears
	} else if nYearMod == 2 {
		nDays += secondYears
		if !neg {
			pMonth = &leapMonth
		}
	} else if nYearMod == 3 {
		nDays += thirdYears
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
	daySecond = hour*hourSec + min*minSec + sec
	unix = int64(nDays*daySec+daySecond) - zone.Offset()
	//return unix, yDay, daySecond, nil
	return
}

//@description: 日期时间字符串 -> 根据格式化模板 -> 转换为秒级时间戳
//@param:       s string "日期时间字符串"
//@param:       formatter string "格式化模板" 如: "%Y/%m/%d %H:%M:%S", "%Y-%m-%d %H:%M:%S", "%Y%m%d%H%M%S"
//@param:       zone tz.TimeZone "时区"
//@param:       extend bool "是否启用扩展增强模式" 与函数 FormatToDateClock 一样
//@return:      unix int64 "秒级时间戳"
//@return:      error "错误信息"
func FormatToUnix(s, formatter string, zone tz.TimeZone, extend bool) (unix int64, err error) {
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

//@description: 标准日期时间字符串 -> 转换为秒级时间戳
//@param:       s string "日期时间字符串"
//@param:       zone tz.TimeZone "时区"
//@param:       extend bool "是否启用扩展增强模式" 与函数 FormatToDateClock 一样
//@return:      unix int64 "秒级时间戳"
//@return:      error "错误信息"
func YmdHMSToUnix(s string, zone tz.TimeZone, extend bool) (unix int64, err error) {
	return FormatToUnix(s, formatterYmdHMS, zone, extend)
}

//@description: 秒级时间戳 ->转为 年,月,日,时,分,秒
//@param:       unix int64 "秒级时间戳"
//@param:       zone tz.TimeZone "时区"
//@return:      year, month, day, hour, min, sec "年,月,日,时,分,秒"
//@return:      yDay "一年中第几天"
//@return:      daySecond "一天中第几秒"
func UnixToDateClock(unix int64, zone tz.TimeZone) (year, month, day, hour, min, sec, yDay, daySecond int) {
	unixLocal := unix + zone.Offset()
	var nRemain int
	if unixLocal < 0 {
		nUnixSec := -unixLocal
		nDays := int(nUnixSec / daySec)
		daySecond = (daySec - int(nUnixSec-int64(nDays*daySec))) % daySec
		nYear4 := nDays/fourYears + 1
		nRemain = nYear4*fourYears - nDays
		if daySecond == 0 {
			nRemain += 1
		}
		year = 1970 - nYear4<<2
	} else {
		nDays := int(unixLocal / daySec)
		daySecond = int(unixLocal - int64(nDays*daySec))
		nYear4 := nDays / fourYears
		nRemain = nDays - nYear4*fourYears + 1
		year = 1970 + nYear4<<2
	}
	pMonth := &norMonth
	if nRemain <= firstYears {

	} else if nRemain <= secondYears {
		year += 1
		nRemain -= firstYears
	} else if nRemain <= thirdYears {
		year += 2
		nRemain -= secondYears
		pMonth = &leapMonth
	} else if nRemain <= fourYears {
		year += 3
		nRemain -= thirdYears
	} else {
		year += 4
		nRemain -= fourYears
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
	hour = daySecond / hourSec
	inHourSec := daySecond - hour*hourSec
	min = inHourSec / minSec
	sec = inHourSec - min*minSec
	return
}

//@description: 日期时间字符串中获取 -> 年,月,日,时,分,秒
//@param:       s string "日期时间字符串"
//@param:       formatter string "格式化字符串"
//@param:       extend bool "是否启用扩展模式"
//              扩展模式: 可以识别非标准日期时间字符串 如: 2020/1/1 0:1:1
//              非扩展模式: 只能识别标准日期时间字符串 如: 2020/01/01 00:01:01
//@return:      year, month, day, hour, min, sec int "日期时间"
//@return:      error "错误信息"
func FormatToDateClock(s, formatter string, extend bool) (year, month, day, hour, min, sec int, err error) {
	defer internal.Exception(func(stack string, e error) {
		err = internal.NewError("format to date clock error exception: %v, %v \n%v", s, formatter, stack)
	})
	numbers := basal.NewNextNumber(s)
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

//@description: 返回时间戳所在的时间是周几, 星期1为一周的开始
//@param:       unix int64 "秒级时间戳"
//@param:       zone tz.TimeZone "时区"
//@return:      week int "星期(1-7)"
func UnixWeekdayA(unix int64, zone tz.TimeZone) (week int) {
	unixLocal := unix + zone.Offset()
	if unixLocal < 0 {
		nSecond := int(unixLocal%weekSec+weekSec) % weekSec
		week = nSecond/daySec + 4
	} else {
		week = int(unixLocal%weekSec/daySec + 4)
	}
	if week > 7 {
		week = week - 7
	}
	return
}

//@description: 返回时间戳所在的时间是周几, 星期天为一周的开始
//@param:       unix int64 "秒级时间戳"
//@param:       zone tz.TimeZone "时区"
//@return:      week int "星期(0-6)"
func UnixWeekdayB(unix int64, zone tz.TimeZone) (week int) {
	week = UnixWeekdayA(unix, zone)
	if week == 7 {
		week = 0
	}
	return
}

//@description: 返回时间戳年中的星期数, 星期1为一周的开始
//@param:       unix int64 "秒级时间戳"
//@param:       zone tz.TimeZone "时区"
//@return:      int(1-53) "第几周"
func UnixYearWeekNumAByISO(unix int64, zone tz.TimeZone) (int, int) {
	unixWeek4, _ := UnixWeekDayZeroHourA(unix, 4, zone)
	yearStart := UnixYearZeroHour(unix, zone)
	yearEnd := UnixYearEndTime(unix, zone)
	dt := UnixToDateTime(unixWeek4, zone)
	//本周四在下一年
	if unixWeek4 > yearEnd {
		return dt.Year(), 1
	} else if unixWeek4 < yearStart { //本周4在上一年
		unixWeek1 := unixWeek4 - 3*daySec
		nSecond := int(unix - unixWeek1)
		return dt.Year(), (nSecond / weekSec) + 1
	} else {
		unixYearStartWeek1, _ := UnixWeekDayZeroHourA(yearStart, 1, zone)
		nSecond := int(unix - unixYearStartWeek1)
		return dt.Year(), (nSecond / weekSec) + 1
	}
}

//@description: 返回时间戳1970年1月1日以来的天数
//@param:       unix int64 "秒级时间戳"
//@param:       zone tz.TimeZone "时区"
//@return:      int64 "天数"
func UnixDayNumber(unix int64, zone tz.TimeZone) int64 {
	return (unix + zone.Offset()) / int64(daySec)
}

//@description: 返回时间戳所在年的1月1日0时的秒级时间戳
//@param:       unix int64 "秒级时间戳"
//@param:       zone tz.TimeZone "时区"
//@return:      int64 "秒级时间戳"
func UnixYearZeroHour(unix int64, zone tz.TimeZone) int64 {
	_, _, _, _, _, _, yDay, daySecond := UnixToDateClock(unix, zone)
	unixLocal := unix + zone.Offset()
	return unixLocal - int64((yDay-1)*daySec+daySecond) - zone.Offset()
}

//@description: 返回时间戳所在月1日0时的秒级时间戳
//@param:       unix int64 "秒级时间戳"
//@param:       zone tz.TimeZone "时区"
//@return:      int64 "秒级时间戳"
func UnixMonthZeroHour(unix int64, zone tz.TimeZone) int64 {
	year, month, _, _, _, _, _, _ := UnixToDateClock(unix, zone)
	unixMon, _, _, _ := DateClockToUnix(year, month, 1, 0, 0, 0, zone)
	return unixMon
}

//@description: 返回时间戳当天0时的秒级时间戳
//@param:       unix int64 "秒级时间戳"
//@param:       zone tz.TimeZone "时区"
//@return:      int64 "秒级时间戳"
func UnixDayZeroHour(unix int64, zone tz.TimeZone) int64 {
	unixLocal := unix + zone.Offset()
	if unixLocal < 0 {
		nSecond := (unixLocal%daySec + daySec) % daySec
		return unixLocal - nSecond - zone.Offset()
	}
	return unixLocal - unixLocal%daySec - zone.Offset()
}

//@description: 返回时间戳本小时的0分的秒级时间戳
//@param:       unix int64 "秒级时间戳"
//@return:      int64 "秒级时间戳"
func UnixHourZeroMin(unix int64) int64 {
	if unix < 0 {
		nSecond := (unix%hourSec + hourSec) % daySec
		return unix - nSecond
	}
	return unix - unix%hourSec
}

//@description: 返回时间戳当天特定时间的秒级时间戳
//@param:       unix int64 "秒级时间戳"
//@param:       hour, min, sec int "时,分,秒"
//@param:       zone tz.TimeZone "时区"
//@return:      int64 "秒级时间戳"
func UnixThisDay(unix int64, hour, min, sec int, zone tz.TimeZone) (int64, error) {
	if err := checkClock(hour, min, sec); err != nil {
		return 0, err
	}
	start := UnixDayZeroHour(unix, zone)
	return start + int64(hour*hourSec+min*minSec+sec), nil
}

//@description: 返回时间戳后N天特定时间的秒级时间戳
//@param:       unix int64 "秒级时间戳"
//@param:       days, hour, min, sec int "天数,时,分,秒"
//@param:       zone tz.TimeZone "时区"
//@return:      int64 "秒级时间戳"
func UnixDayZeroHourNext(unix int64, days, hour, min, sec int, zone tz.TimeZone) (int64, error) {
	if err := checkClock(hour, min, sec); err != nil {
		return 0, err
	}
	start := UnixDayZeroHour(unix, zone)
	return start + int64(days*daySec+hour*hourSec+min*minSec+sec), nil
}

//返回当天0点时间戳加上sec的时间戳 单位: 秒
func UnixDayZeroHourNextSec(unix int64, sec int, zone tz.TimeZone) int64 {
	start := UnixDayZeroHour(unix, zone)
	return start + int64(sec)
}

//前一个时间点
func UnixPreClock(unix int64, hour, min, sec int, zone tz.TimeZone) (int64, error) {
	nextUnix, err := UnixDayZeroHourNext(unix, 0, hour, min, sec, zone)
	if err != nil {
		return 0, err
	}
	if unix < nextUnix {
		return nextUnix - daySec, nil
	} else {
		return nextUnix, nil
	}
}

//下一个时间点
func UnixNextClock(unix int64, hour, min, sec int, zone tz.TimeZone) (int64, error) {
	nextUnix, err := UnixDayZeroHourNext(unix, 0, hour, min, sec, zone)
	if err != nil {
		return 0, err
	}
	if unix < nextUnix {
		return nextUnix, nil
	} else {
		return nextUnix + daySec, nil
	}
}

//@description: 返回时间戳下一周的星期几的秒级时间戳(星期1为周的开始)
//@param:       unix int64 "秒级时间戳"
//@param:       week, hour, min, sec int "星期几(1-7),时,分,秒"
//@param:       zone tz.TimeZone "时区"
//@return:      int64 "秒级时间戳"
//@return:      error "错误信息"
func UnixNextWeekDayA(unix int64, week int, hour, min, sec int, zone tz.TimeZone) (int64, error) {
	if week < 1 || week > 7 {
		return 0, internal.NewError("week out of range(1,7): %v", week)
	}
	w := UnixWeekdayA(unix, zone)
	days := week - w
	return UnixDayZeroHour(unix, zone) + int64(weekSec+days*daySec+hour*hourSec+min*minSec+sec), nil
}

//@description: 返回时间戳下一周的星期几的秒级时间戳(星期天为周的开始)
//@param:       unix int64 "秒级时间戳"
//@param:       week, hour, min, sec int "星期几(0-6),时,分,秒"
//@param:       zone tz.TimeZone "时区"
//@return:      int64 "秒级时间戳"
//@return:      error "错误信息"
func UnixNextWeekDayB(unix int64, week, hour, min, sec int, zone tz.TimeZone) (int64, error) {
	if week < 0 || week > 6 {
		return 0, internal.NewError("week out of range(0, 6): %v", week)
	}
	w := UnixWeekdayB(unix, zone)
	days := week - w
	return UnixDayZeroHour(unix, zone) + int64(weekSec+days*daySec+hour*hourSec+min*minSec+sec), nil
}

//@description: 返回时间戳下一个最近的星期几的秒级时间戳(星期1为周的开始)
//@param:       unix int64 "秒级时间戳"
//@param:       week, hour, min, sec int "星期几(1-7),时,分,秒"
//@param:       zone tz.TimeZone "时区"
//@return:      int64 "秒级时间戳"
//@return:      error "错误信息"
func UnixFutureWeekDayA(unix int64, week, hour, min, sec int, zone tz.TimeZone) (int64, error) {
	if week < 1 || week > 7 {
		return 0, internal.NewError("week out of range(1,7): %v", week)
	}
	if err := checkClock(hour, min, sec); err != nil {
		return 0, err
	}
	w := UnixWeekdayA(unix, zone)
	days := week - w
	unixDaySecond := hour*hourSec + min*minSec + sec
	if days > 0 {
		return UnixDayZeroHour(unix, zone) + int64(days*daySec+unixDaySecond), nil
	} else if days == 0 {
		_, _, _, _, _, _, _, daySecond := UnixToDateClock(unix, zone)
		if unixDaySecond > daySecond {
			return UnixDayZeroHour(unix, zone) + int64(unixDaySecond), nil
		} else {
			return UnixDayZeroHour(unix, zone) + int64(weekSec+days*daySec+unixDaySecond), nil
		}
	} else {
		return UnixDayZeroHour(unix, zone) + int64(weekSec+days*daySec+unixDaySecond), nil
	}
}

//@description: 返回时间戳下一个最近的星期几的秒级时间戳(星期天为周的开始)
//@param:       unix int64 "秒级时间戳"
//@param:       week, hour, min, sec int "星期几(0-6),时,分,秒"
//@param:       zone tz.TimeZone "时区"
//@return:      int64 "秒级时间戳"
//@return:      error "错误信息"
func UnixFutureWeekDayB(unix int64, week, hour, min, sec int, zone tz.TimeZone) (int64, error) {
	if week < 0 || week > 6 {
		return 0, internal.NewError("week out of range(0, 6): %v", week)
	}
	if err := checkClock(hour, min, sec); err != nil {
		return 0, err
	}
	w := UnixWeekdayB(unix, zone)
	days := week - w
	unixDaySecond := hour*hourSec + min*minSec + sec
	if days > 0 {
		return UnixDayZeroHour(unix, zone) + int64(days*daySec+unixDaySecond), nil
	} else if days == 0 {
		_, _, _, _, _, _, _, daySecond := UnixToDateClock(unix, zone)
		if unixDaySecond > daySecond {
			return UnixDayZeroHour(unix, zone) + int64(unixDaySecond), nil
		} else {
			return UnixDayZeroHour(unix, zone) + int64(weekSec+days*daySec+unixDaySecond), nil
		}
	} else {
		return UnixDayZeroHour(unix, zone) + int64(weekSec+days*daySec+unixDaySecond), nil
	}
}

//@description: 返回时间戳所在周的开始时间（周一的0点）
//@param:       unix int64 "秒级时间戳"
//@return:      int64 "秒级时间戳"
func UnixWeekStartTimeA(unix int64, zone tz.TimeZone) int64 {
	return UnixDayZeroHour(unix, zone) + int64(1-UnixWeekdayA(unix, zone))*daySec
}

//@description: 返回时间戳所在周的结束时间（周日的23:59:59）
//@param:       unix int64 "秒级时间戳"
//@return:      int64 "秒级时间戳"
func UnixWeekEndTimeA(unix int64, zone tz.TimeZone) int64 {
	return UnixDayZeroHour(unix, zone) + int64(8-UnixWeekdayA(unix, zone))*daySec - 1
}

//@description: 返回时间戳所在周几开始时间（周一的0点）
//@param:       unix int64 "秒级时间戳"
//@param:       week 1-7
//@return:      int64 "秒级时间戳"
func UnixWeekDayZeroHourA(unix int64, week int, zone tz.TimeZone) (int64, error) {
	if week < 1 || week > 7 {
		return 0, internal.NewError("week out of range(1, 7): %v", week)
	}
	return UnixDayZeroHourNext(unix, week-UnixWeekdayA(unix, zone), 0, 0, 0, zone)
}

//@description: 返回时间戳所在周的任何时间
//@param:       unix int64 "秒级时间戳"
//@param:       week 1-7
//@return:      int64 "秒级时间戳"
func UnixWeekAnyTimeA(unix int64, week, hour, min, sec int, zone tz.TimeZone) (int64, error) {
	if week < 1 || week > 7 {
		return 0, internal.NewError("week out of range(1, 7): %v", week)
	}
	return UnixDayZeroHourNext(unix, week-UnixWeekdayA(unix, zone), hour, min, sec, zone)
}

//@description: 判断时间戳是否在一段时间范围 支持跨周 eg: 7_12:00:00 => 3_13:00:00 / 7+12:00:00 => 3+13:00:00
//@param:       unix int64 "秒级时间戳"
//@param:       week 1-7
func UnixInBetweenWeekTimeA(unix int64, startTime, endTime string, zone tz.TimeZone) (bool, error) {
	startNums, endNums := basal.NewNextNumber(startTime).Numbers(), basal.NewNextNumber(endTime).Numbers()
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

//@description: 返回时间戳所在月的任意时间点
//@param:       unix int64 "秒级时间戳"
//@return:      int64 "秒级时间戳"
func UnixMonthAnyTime(unix int64, day, hour, min, sec int, zone tz.TimeZone) (int64, error) {
	year, month, _, _, _, _, _, _ := UnixToDateClock(unix, zone)
	if unixMonth, _, _, err := DateClockToUnix(year, month, day, hour, min, sec, zone); err != nil {
		return 0, err
	} else {
		return unixMonth, nil
	}
}

//@description: 返回时间戳所在年内某月的的任何时间
//@param:       unix int64 "秒级时间戳"
//@return:      int64 "秒级时间戳"
func UnixYearAnyTime(unix int64, month, day, hour, min, sec int, zone tz.TimeZone) (int64, error) {
	year, _, _, _, _, _, _, _ := UnixToDateClock(unix, zone)
	if unixMonth, _, _, err := DateClockToUnix(year, month, day, hour, min, sec, zone); err != nil {
		return 0, err
	} else {
		return unixMonth, nil
	}
}

//@description: 返回时间戳所在年的开始时间
//@param:       unix int64 "秒级时间戳"
//@return:      int64 "秒级时间戳"
func UnixYearStartTime(unix int64, zone tz.TimeZone) int64 {
	yearStartUnix, _ := UnixYearAnyTime(unix, 1, 1, 0, 0, 0, zone)
	return yearStartUnix
}

//@description: 返回时间戳所在年的结束时间
//@param:       unix int64 "秒级时间戳"
//@return:      int64 "秒级时间戳"
func UnixYearEndTime(unix int64, zone tz.TimeZone) int64 {
	year, _, _, _, _, _, _, _ := UnixToDateClock(unix, zone)
	nextYearStartUnix, _, _, _ := DateClockToUnix(year+1, 1, 1, 0, 0, 0, zone)
	return nextYearStartUnix - 1
}

//@description: 判断时间戳是否在一段时间范围内 eg: 08-15 12:00:00 => 01-15 13:00:00
func UnixInBetweenMonthTime(unix int64, startTime string, endTime string, zone tz.TimeZone) (bool, error) {
	startNums, endNums := basal.NewNextNumber(startTime).Numbers(), basal.NewNextNumber(endTime).Numbers()
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

//@description: 返回时间戳所在月1日0时的秒级时间戳
//@param:       unix int64 "秒级时间戳"
//@param:       zone tz.TimeZone "时区"
//@return:      int64 "秒级时间戳"
func UnixMonthStartTime(unix int64, zone tz.TimeZone) int64 {
	year, month, _, _, _, _, _, _ := UnixToDateClock(unix, zone)
	unixMonth, _, _, _ := DateClockToUnix(year, month, 1, 0, 0, 0, zone)
	return unixMonth
}

//@description: 返回时间戳所在月的最后一天23:59:59
//@param:       unix int64 "秒级时间戳"
//@param:       zone tz.TimeZone "时区"
//@return:      int64 "秒级时间戳"
func UnixMonthEndTime(unix int64, zone tz.TimeZone) int64 {
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

//@description: 返回时间戳当天0时的秒级时间戳
//@param:       unix int64 "秒级时间戳"
//@param:       zone tz.TimeZone "时区"
//@return:      int64 "秒级时间戳"
func UnixDayStartTime(unix int64, zone tz.TimeZone) int64 {
	unixLocal := unix + zone.Offset()
	if unixLocal < 0 {
		nSecond := (unixLocal%daySec + daySec) % daySec
		return unixLocal - nSecond - zone.Offset()
	}
	return unixLocal - unixLocal%daySec - zone.Offset()
}

//@description: 返回时间戳当天23:59:59
//@param:       unix int64 "秒级时间戳"
//@param:       zone tz.TimeZone "时区"
//@return:      int64 "秒级时间戳"
func UnixDayEndTime(unix int64, zone tz.TimeZone) int64 {
	return UnixDayStartTime(unix, zone) + daySec - 1
}

var hmsLenErr = internal.NewError("hms string len error")

//时间转换为clock
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

//时间转换为clock
func HmsToSec(hms string) (int64, error) {
	hour, min, sec, err := HmsToClock(hms)
	if err != nil {
		return 0, err
	}

	return int64(hour*hourSec + min*minSec + sec), nil
}

//@description: 转换为总秒数
func ClockToSec(hour, min, sec int) (int64, error) {
	if err := checkClock(hour, min, sec); err != nil {
		return 0, err
	}
	return int64(hour*hourSec + min*minSec + sec), nil
}

//@description: 根据时间戳得到当天的任意时间
func UnixDayAnyTime(unix int64, hour, min, sec int, zone tz.TimeZone) (int64, error) {
	if secNum, err := ClockToSec(hour, min, sec); err != nil {
		return 0, err
	} else {
		return UnixDayStartTime(unix, zone) + secNum, nil
	}
}

//@description: 判断时间戳是否在时间范围内 小时级别 eg: 12:00:00 => 13:00:00
//@param:       unix int64 "秒级时间戳"
func UnixInBetweenHourTime(unix int64, startTime string, endTime string, zone tz.TimeZone) (bool, error) {
	startNums, endNums := basal.NewNextNumber(startTime).Numbers(), basal.NewNextNumber(endTime).Numbers()
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

//@description: 判断时间戳是否在时间范围内 小时级别 eg: 12:00:00 => 13:00:00
//@param:       unix int64 "秒级时间戳"
//@return: 是否在时间范围内, 开始时间戳, 结束时间戳, 错误信息
func UnixScopeHourTime(unix int64, startHms string, endHms string, zone tz.TimeZone) (bool, int64, int64, error) {
	startNums, endNums := basal.NewNextNumber(startHms).Numbers(), basal.NewNextNumber(endHms).Numbers()
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
			start, err := UnixDayZeroHourNext(unix, -1, startNums[0], startNums[1], startNums[2], zone)
			if err != nil {
				return false, 0, 0, err
			}
			return true, start, endUnix, nil
		}
		dayEndTime := UnixDayEndTime(unix, zone)
		if unix >= startUnix && unix <= dayEndTime {
			end, err := UnixDayZeroHourNext(unix, 1, endNums[0], endNums[1], endNums[2], zone)
			if err != nil {
				return false, 0, 0, err
			}
			return true, startUnix, end, nil
		}
		nextStart, err := UnixDayZeroHourNext(unix, 0, startNums[0], startNums[1], startNums[2], zone)
		if err != nil {
			return false, 0, 0, err
		}
		nextEnd, err := UnixDayZeroHourNext(unix, 0, endNums[0], endNums[1], endNums[2], zone)
		if err != nil {
			return false, 0, 0, err
		}
		return false, nextStart, nextEnd, nil
	}
}

//返回时间戳的 格式化日期时间字符串
func UnixToFormat(unix int64, zone tz.TimeZone, formatter string) string {
	year, month, day, hour, min, sec, _, _ := UnixToDateClock(unix, zone)
	return DateClockToFormat(year, month, day, hour, min, sec, formatter)
}

//返回时间戳的 标准日期时间字符串
func UnixToYmdHMS(unix int64, zone tz.TimeZone) string {
	return UnixToFormat(unix, zone, formatterYmdHMS)
}

//判断时间戳是否是同一天
func UnixSameDay(unix1, unix2, offset int64, zone tz.TimeZone) bool {
	start := UnixDayStartTime(unix1-offset, zone)
	end := UnixDayEndTime(unix1-offset, zone)
	if unix2-offset < start {
		return false
	} else if unix2-offset > end {
		return false
	}
	return true
}

//判断时间戳是否是同一周, 星期1为周的开始
func UnixSameWeekA(unix1, unix2, offset int64, zone tz.TimeZone) bool {
	start, end := UnixWeekStartTimeA(unix1-offset, zone), UnixWeekEndTimeA(unix1-offset, zone)
	if unix2-offset < start {
		return false
	} else if unix2-offset > end {
		return false
	}
	return true
}

//判断时间戳是否是同一月
func UnixSameMonth(unix1, unix2, offset int64, zone tz.TimeZone) bool {
	dt1, dt2 := UnixToDateTime(unix1-offset, zone), UnixToDateTime(unix2-offset, zone)
	return dt1.Year() == dt2.Year() && dt1.Month() == dt2.Month()
}

//当前分的0秒时间戳（秒）
func UnixThisMinZeroSec(unix int64) int64 {
	return (unix / 60) * 60
}

//开始时间和结束时间相差月份
func UnixDiffMonth(start, end int64, zone tz.TimeZone) int {
	d1 := UnixToDateTime(start, zone)
	d2 := UnixToDateTime(end, zone)
	return DateTimeDiffMonth(d1, d2)
}

//赛季标识
//startUnix 赛季开始时间戳 单位: 秒
//unix 时间戳
//monthNum 一个赛季为多少月
//offset 偏移量 秒数 最大3600*24*28
func SeasonId(startUnix int64, unix int64, monthNum, offset int, zone tz.TimeZone) int {
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

//获得赛季时间范围
//startUnix 赛季开始时间戳 单位: 秒
//season 赛季id
//monthNum 一个赛季为多少月
//firstZeroTime 第一个赛季是否从月的第一天开始, 否则从真实时间开始
//offset 偏移量 秒数 最大3600*24*28
func SeasonDateTime(startUnix int64, seasonId, monthNum, offset int, firstZeroTime bool, zone tz.TimeZone) (begin *DateTime, end *DateTime) {
	start := UnixToDateTime(startUnix, zone)
	addMonthNum := seasonId * monthNum
	if seasonId == 0 && firstZeroTime == false {
		begin = start
	} else {
		begin = start.MonthStartDateTime(addMonthNum).AddSec(offset)
	}
	end = begin.MonthEndDateTime(monthNum - 1).AddSec(offset)
	return begin, end
}
