package internal

import (
	tz "github.com/jingyanbin/core/timezone"
	_ "unsafe"
)

type DateTime struct {
	unix      int64 //秒级时间戳
	year      int
	month     int
	day       int
	hour      int
	min       int
	sec       int
	yDay      int //年中第几天
	zone      tz.TimeZone
	daySecond int //天中第几秒
}

func (m *DateTime) Unix() int64 {
	return m.unix
}

func (m *DateTime) Year() int {
	return m.year
}

func (m *DateTime) Month() int {
	return m.month
}

func (m *DateTime) Day() int {
	return m.day
}

func (m *DateTime) Hour() int {
	return m.hour
}

func (m *DateTime) Min() int {
	return m.min
}

func (m *DateTime) Sec() int {
	return m.sec
}

func (m *DateTime) YDay() int {
	return m.yDay
}

func (m *DateTime) Zone() tz.TimeZone {
	return m.zone
}

func (m *DateTime) DaySecond() int {
	return m.daySecond
}

func (m *DateTime) SetZone(zone tz.TimeZone) {
	if m.zone.Offset() != zone.Offset() {
		m.zone = zone
		if m.unix == 0 {
			return
		}
		m.flush()
	}
}

// 标准返回年周
func (m *DateTime) UnixYearWeekNumAByISO() (int, int) {
	return UnixYearWeekNumAByISO(m.unix, m.zone)
}

// 周几1-7
func (m *DateTime) WeekdayA() int {
	return UnixWeekdayA(m.unix, m.zone)
}

// 周几0-6
func (m *DateTime) WeekdayB() (week int) {
	return UnixWeekdayB(m.unix, m.zone)
}

// 天数
func (m *DateTime) UnixDayNumber() int64 {
	return UnixDayNumber(m.unix, m.zone)
}

func (m *DateTime) flush() {
	m.year, m.month, m.day, m.hour, m.min, m.sec, m.yDay, m.daySecond = UnixToDateClock(m.unix, m.zone)
}

// 刷新到当前时间
func (m *DateTime) Flush() {
	unix := Unix()
	if unix == m.unix {
		return
	}
	m.unix = unix
	m.flush()
}

// 刷新时间到指定秒级时间戳
func (m *DateTime) FlushToUnix(unix int64) {
	if unix == m.unix {
		return
	}
	m.unix = unix
	m.flush()
}

// 刷新时间到指定日期时间
func (m *DateTime) FlushToDateClock(year, month, day, hour, min, sec int) error {
	unix, yDay, daySecond, err := DateClockToUnix(year, month, day, hour, min, sec, m.zone)
	if err != nil {
		return err
	}
	if m.unix == unix {
		return nil
	}
	m.unix = unix
	m.year = year
	m.month = month
	m.day = day
	m.hour = hour
	m.min = min
	m.sec = sec
	m.yDay = yDay
	m.daySecond = daySecond
	return nil
}

// 刷新时间到 日期时间字符串
func (m *DateTime) FlushToFormat(s, formatter string, extend bool) error {
	year, month, day, hour, min, sec, err := FormatToDateClock(s, formatter, extend)
	if err != nil {
		return err
	}
	return m.FlushToDateClock(year, month, day, hour, min, sec)
}

// 刷新到 标准日期时间字符串
func (m *DateTime) FlushToYmdHMS(s string, extend bool) error {
	return m.FlushToFormat(s, FormatterYmdHMS, extend)
}

// 返回时间戳所在年的开始时间
func (m *DateTime) UnixYearStartTime() int64 {
	return UnixYearStartTime(m.unix, m.zone)
}

// 返回时间戳所在月1日0时的秒级时间戳
func (m *DateTime) UnixMonthStartTime() int64 {
	unixMon, _, _, _ := DateClockToUnix(m.year, m.month, 1, 0, 0, 0, m.zone)
	return unixMon
}

// 返回当日0时的秒级时间戳
func (m *DateTime) UnixDayStartTime() int64 {
	return UnixDayStartTime(m.unix, m.zone)
}

// 返回时间戳本小时的0分的秒级时间戳
func (m *DateTime) UnixHourStartTime() int64 {
	return UnixHourStartTime(m.unix)
}

// 返回本小时开始时间的下面的花的时间
func (m *DateTime) UnixHourStartTimeNext(days, hour, min, sec int) (int64, error) {
	return UnixHourStartTimeNext(m.unix, days, hour, min, sec, m.zone)
}

// 下周几1-7
func (m *DateTime) UnixNextWeekDayA(week int, hour, min, sec int) (int64, error) {
	return UnixNextWeekDayA(m.unix, week, hour, min, sec, m.zone)
}

// 下周几0-6
func (m *DateTime) UnixNextWeekDayB(week int, hour, min, sec int) (int64, error) {
	return UnixNextWeekDayB(m.unix, week, hour, min, sec, m.zone)
}

// 下一个周几1-7
func (m *DateTime) UnixFutureWeekDayA(week, hour, min, sec int) (int64, error) {
	return UnixFutureWeekDayA(m.unix, week, hour, min, sec, m.zone)
}

// 下一个周几0-6
func (m *DateTime) UnixFutureWeekDayB(week, hour, min, sec int) (int64, error) {
	return UnixFutureWeekDayB(m.unix, week, hour, min, sec, m.zone)
}

// 当月开始时间
func (m *DateTime) MonthStartDateTime(addMonthNum int) *DateTime {
	year, month := m.year, m.month
	if addMonthNum != 0 {
		year, month = YearMonthByAddMonthNum(m.year, m.month, addMonthNum)
	}
	dt, err := DateClockToDateTime(year, month, 1, 0, 0, 0, m.zone)
	if err != nil {
		panic(err)
	}
	return dt
}

// 当月结束时间
func (m *DateTime) MonthEndDateTime(addMonthNum int) *DateTime {
	year, month := m.year, m.month
	if addMonthNum != 0 {
		year, month = YearMonthByAddMonthNum(m.year, m.month, addMonthNum)
	}
	days := MonthDayNumber(year, month)
	dt, err := DateClockToDateTime(year, month, days, 23, 59, 59, m.zone)
	if err != nil {
		panic(err)
	}
	return dt
}

// 月份差
func (m *DateTime) DiffMonth(dt *DateTime) (monthNum int) {
	return DateTimeDiffMonth(dt, m)
}

// 加天数
func (m *DateTime) Add(days, hour, min, sec int) *DateTime {
	return UnixToDateTime(m.unix+int64(days*DaySec+hour*HourSec+min*MinSec+sec), m.zone)
}

// 加秒数
func (m *DateTime) AddSec(sec int) *DateTime {
	return UnixToDateTime(m.unix+int64(sec), m.zone)
}

// 格式化日期时间字符串
func (m *DateTime) Format(formatter string) string {
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
				ItoAW(&theTime, m.year, 4)
			case 'y': //两位数的年份表示（00-99）
				ItoAW(&theTime, m.year, 2)
			case 'm': //月份（01-12）
				ItoAW(&theTime, m.month, 2)
			case 'd': //月内中的一天（0-31）
				ItoAW(&theTime, m.day, 2)
			case 'H': //24小时制小时数（0-23）
				ItoAW(&theTime, m.hour, 2)
			case 'M': //分钟数（00=59）
				ItoAW(&theTime, m.min, 2)
			case 'S': //秒（00-59）
				ItoAW(&theTime, m.sec, 2)
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

// 返回标准日期时间字符串
func (m *DateTime) YmdHMS() string {
	return m.Format(FormatterYmdHMS)
}

// 当前日期时间
//
//go:linkname Now github.com/jingyanbin/core/datetime.Now
func Now() (dt *DateTime) {
	dt = &DateTime{zone: tz.Local()}
	dt.Flush()
	return
}

// 时间戳转日期时间对象
//
//go:linkname UnixToDateTime github.com/jingyanbin/core/datetime.UnixToDateTime
func UnixToDateTime(unix int64, zone tz.TimeZone) (dt *DateTime) {
	dt = &DateTime{unix: unix, zone: zone}
	dt.flush()
	return
}

// 格式化日期时间字符串 转换为DateTime
//
//go:linkname FormatToDateTime github.com/jingyanbin/core/datetime.FormatToDateTime
func FormatToDateTime(s, formatter string, zone tz.TimeZone, extend bool) (dt *DateTime, err error) {
	dt = &DateTime{zone: zone}
	err = dt.FlushToFormat(s, formatter, extend)
	if err != nil {
		return nil, err
	}
	return dt, nil
}

// 标准日期时间字符串 转换为DateTime
//
//go:linkname YmdHMSToDateTime github.com/jingyanbin/core/datetime.YmdHMSToDateTime
func YmdHMSToDateTime(s string, zone tz.TimeZone, extend bool) (dt *DateTime, err error) {
	dt = &DateTime{zone: zone}
	err = dt.FlushToYmdHMS(s, extend)
	if err != nil {
		return nil, err
	}
	return dt, nil
}

// 年,月,日,时,分,秒 转换为DateTime
//
//go:linkname DateClockToDateTime github.com/jingyanbin/core/datetime.DateClockToDateTime
func DateClockToDateTime(year, month, day, hour, min, sec int, zone tz.TimeZone) (dt *DateTime, err error) {
	dt = &DateTime{zone: zone}
	err = dt.FlushToDateClock(year, month, day, hour, min, sec)
	if err != nil {
		return nil, err
	}
	return dt, nil
}

// 根据年月加减月数获得新的年月
//
//go:linkname YearMonthByAddMonthNum github.com/jingyanbin/core/datetime.YearMonthByAddMonthNum
func YearMonthByAddMonthNum(year, month, addMonthNum int) (y int, m int) {
	monthCount := month + addMonthNum
	m = monthCount % 12
	y = year + monthCount/12
	if m < 1 {
		m = 12 + m
		y--
	}
	return y, m
}

// 开始时间和结束时间相差月份
//
//go:linkname DateTimeDiffMonth github.com/jingyanbin/core/datetime.DateTimeDiffMonth
func DateTimeDiffMonth(start, end *DateTime) int {
	yearCha := end.Year() - start.Year()
	return yearCha*12 + (end.Month() - start.Month())
}
