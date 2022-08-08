package datetime

import (
	internal "github.com/jingyanbin/core/internal"
	tz "github.com/jingyanbin/core/timezone"
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

func (my *DateTime) Unix() int64 {
	return my.unix
}

func (my *DateTime) Year() int {
	return my.year
}

func (my *DateTime) Month() int {
	return my.month
}

func (my *DateTime) Day() int {
	return my.day
}

func (my *DateTime) Hour() int {
	return my.hour
}

func (my *DateTime) Min() int {
	return my.min
}

func (my *DateTime) Sec() int {
	return my.sec
}

func (my *DateTime) YDay() int {
	return my.yDay
}

func (my *DateTime) Zone() tz.TimeZone {
	return my.zone
}

func (my *DateTime) DaySecond() int {
	return my.daySecond
}

func (my *DateTime) SetZone(zone tz.TimeZone) {
	if my.zone.Offset() != zone.Offset() {
		my.zone = zone
		if my.unix == 0 {
			return
		}
		my.flush()
	}
}

// UnixYearWeekNumAByISO @description: 返回星期数, 星期1为一周的开始
//@return:      int(1-53) "第几周"
func (my *DateTime) UnixYearWeekNumAByISO() (int, int) {
	return UnixYearWeekNumAByISO(my.unix, my.zone)
}

// WeekdayA @description: 返回周几, 星期1为一周的开始
//@return:      week int "星期(1-7)"
func (my *DateTime) WeekdayA() int {
	return UnixWeekdayA(my.unix, my.zone)
}

//@description: 返回周几, 星期天为一周的开始
//@return:      week int "星期(0-6)"
func (my *DateTime) WeekdayB() (week int) {
	return UnixWeekdayB(my.unix, my.zone)
}

//@description: 返回时间戳1970年1月1日以来的天数
//@return:      int64 "天数"
func (my *DateTime) UnixDayNumber() int64 {
	return UnixDayNumber(my.unix, my.zone)
}

func (my *DateTime) flush() {
	my.year, my.month, my.day, my.hour, my.min, my.sec, my.yDay, my.daySecond = UnixToDateClock(my.unix, my.zone)
}

//@description: 刷新时间为最新
func (my *DateTime) Flush() {
	unix := Unix()
	if unix == my.unix {
		return
	}
	my.unix = unix
	my.flush()
}

//@description: 刷新时间到指定秒级时间戳
//@param:       unix int64 "秒级时间戳"
func (my *DateTime) FlushToUnix(unix int64) {
	if unix == my.unix {
		return
	}
	my.unix = unix
	my.flush()
}

//@description: 刷新时间到指定日期时间
func (my *DateTime) FlushToDateClock(year, month, day, hour, min, sec int) error {
	unix, yDay, daySecond, err := DateClockToUnix(year, month, day, hour, min, sec, my.zone)
	if err != nil {
		return err
	}
	if my.unix == unix {
		return nil
	}
	my.unix = unix
	my.year = year
	my.month = month
	my.day = day
	my.hour = hour
	my.min = min
	my.sec = sec
	my.yDay = yDay
	my.daySecond = daySecond
	return nil
}

//@description: 刷新时间到 日期时间字符串
//@param:       s string "日期时间字符串"
//@param:       formatter string "格式化字符串"
//@param:       extend bool "是否启用扩展模式" 与函数 FormatToDateClock 一样
//@return:      error "错误信息"
func (my *DateTime) FlushToFormat(s, formatter string, extend bool) error {
	year, month, day, hour, min, sec, err := FormatToDateClock(s, formatter, extend)
	if err != nil {
		return err
	}
	return my.FlushToDateClock(year, month, day, hour, min, sec)
}

//@description: 刷新时间到 标准日期时间字符串
//@param:       s string "标准日期时间字符串"
//@param:       extend bool "是否启用扩展模式" 与函数 FormatToDateClock 一样
//@return:      error "错误信息"
func (my *DateTime) FlushToYmdHMS(s string, extend bool) error {
	return my.FlushToFormat(s, formatterYmdHMS, extend)
}

//@description: 返回1月1日0时的秒级时间戳
//@return:      int64 "秒级时间戳"
func (my *DateTime) UnixYearZeroHour() int64 {
	return UnixYearZeroHour(my.unix, my.zone)
}

//@description: 返回当月1日0时的秒级时间戳
//@return:      int64 "秒级时间戳"
func (my *DateTime) UnixMonthZeroHour() int64 {
	unixMon, _, _, _ := DateClockToUnix(my.year, my.month, 1, 0, 0, 0, my.zone)
	return unixMon
}

//@description: 返回当日0时的秒级时间戳
//@return:      int64 "秒级时间戳"
func (my *DateTime) UnixDayZeroHour() int64 {
	return UnixDayZeroHour(my.unix, my.zone)
}

//@description: 返回本小时0分的秒级时间戳
//@return:      int64 "秒级时间戳"
func (my *DateTime) UnixHourZeroMin() int64 {
	return UnixHourZeroMin(my.unix)
}

//@description: 返回N天后特定时间的秒级时间戳
//@param:       days, hour, min, sec int "天数,时,分,秒"
//@param:       zone tz.TimeZone "时区"
//@return:      int64 "秒级时间戳"
func (my *DateTime) UnixDayZeroHourNext(days, hour, min, sec int) (int64, error) {
	return UnixDayZeroHourNext(my.unix, days, hour, min, sec, my.zone)
}

//@description: 返回下一周的星期几的秒级时间戳(星期1为周的开始)
//@param:       week, hour, min, sec int "星期几(1-7),时,分,秒"
//@return:      int64 "秒级时间戳"
//@return:      error "错误信息"
func (my *DateTime) UnixNextWeekDayA(week int, hour, min, sec int) (int64, error) {
	return UnixNextWeekDayA(my.unix, week, hour, min, sec, my.zone)
}

//@description: 返回下一周的星期几的秒级时间戳(星期天为周的开始)
//@param:       week, hour, min, sec int "星期几(0-6),时,分,秒"
//@return:      int64 "秒级时间戳"
//@return:      error "错误信息"
func (my *DateTime) UnixNextWeekDayB(week int, hour, min, sec int) (int64, error) {
	return UnixNextWeekDayB(my.unix, week, hour, min, sec, my.zone)
}

//@description: 返回下一个最近的星期几的秒级时间戳(星期1为周的开始)
//@param:       unix int64 "秒级时间戳"
//@param:       week, hour, min, sec int "星期几(1-7),时,分,秒"
//@return:      int64 "秒级时间戳"
//@return:      error "错误信息"
func (my *DateTime) UnixFutureWeekDayA(week, hour, min, sec int) (int64, error) {
	return UnixFutureWeekDayA(my.unix, week, hour, min, sec, my.zone)
}

//@description: 返回下一个最近的星期几的秒级时间戳(星期天为周的开始)
//@param:       unix int64 "秒级时间戳"
//@param:       week, hour, min, sec int "星期几(0-6),时,分,秒"
//@return:      int64 "秒级时间戳"
//@return:      error "错误信息"
func (my *DateTime) UnixFutureWeekDayB(week, hour, min, sec int) (int64, error) {
	return UnixFutureWeekDayB(my.unix, week, hour, min, sec, my.zone)
}

//一个月的开始时间
func (my *DateTime) MonthStartDateTime(addMonthNum int) *DateTime {
	year, month := my.year, my.month
	if addMonthNum != 0 {
		year, month = YearMonthByAddMonthNum(my.year, my.month, addMonthNum)
	}
	dt, err := DateClockToDateTime(year, month, 1, 0, 0, 0, my.zone)
	if err != nil {
		panic(err)
	}
	return dt
}

//一个月的结束时间
func (my *DateTime) MonthEndDateTime(addMonthNum int) *DateTime {
	year, month := my.year, my.month
	if addMonthNum != 0 {
		year, month = YearMonthByAddMonthNum(my.year, my.month, addMonthNum)
	}
	days := MonthDayNumber(year, month)
	dt, err := DateClockToDateTime(year, month, days, 23, 59, 59, my.zone)
	if err != nil {
		panic(err)
	}
	return dt
}

//月份差
func (my *DateTime) DiffMonth(dt *DateTime) (monthNum int) {
	return DateTimeDiffMonth(dt, my)
}

//加天数
func (my *DateTime) Add(days, hour, min, sec int) *DateTime {
	return UnixToDateTime(my.unix+int64(days*daySec+hour*hourSec+min*minSec+sec), my.zone)
}

//加秒数
func (my *DateTime) AddSec(sec int) *DateTime {
	return UnixToDateTime(my.unix+int64(sec), my.zone)
}

//@description: 返回格式化日期时间字符串
//@param:       formatter string "格式化字符串"
//@return:      string "日期时间字符串"
func (my *DateTime) Format(formatter string) string {
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
				internal.ItoAW(&theTime, my.year, 4)
			case 'y': //两位数的年份表示（00-99）
				internal.ItoAW(&theTime, my.year, 2)
			case 'm': //月份（01-12）
				internal.ItoAW(&theTime, my.month, 2)
			case 'd': //月内中的一天（0-31）
				internal.ItoAW(&theTime, my.day, 2)
			case 'H': //24小时制小时数（0-23）
				internal.ItoAW(&theTime, my.hour, 2)
			case 'M': //分钟数（00=59）
				internal.ItoAW(&theTime, my.min, 2)
			case 'S': //秒（00-59）
				internal.ItoAW(&theTime, my.sec, 2)
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

//@description: 返回标准日期时间字符串
//@return:      string "标准日期时间字符串"
func (my *DateTime) YmdHMS() string {
	return my.Format(formatterYmdHMS)
}

//@description: 返回当前 DateTime
func Now() (dt *DateTime) {
	dt = &DateTime{zone: tz.Local()}
	dt.Flush()
	return
}

//@description: 秒级时间戳转换为 DateTime
//@param:       unix int64 "秒级时间戳"
//@return:      DateTime
func UnixToDateTime(unix int64, zone tz.TimeZone) (dt *DateTime) {
	dt = &DateTime{unix: unix, zone: zone}
	dt.flush()
	return
}

//@description: 格式化日期时间字符串 转换为DateTime
func FormatToDateTime(s, formatter string, zone tz.TimeZone, extend bool) (dt *DateTime, err error) {
	dt = &DateTime{zone: zone}
	err = dt.FlushToFormat(s, formatter, extend)
	if err != nil {
		return nil, err
	}
	return dt, nil
}

//@description: 标准日期时间字符串 转换为DateTime
func YmdHMSToDateTime(s string, zone tz.TimeZone, extend bool) (dt *DateTime, err error) {
	dt = &DateTime{zone: zone}
	err = dt.FlushToYmdHMS(s, extend)
	if err != nil {
		return nil, err
	}
	return dt, nil
}

//@description: 年,月,日,时,分,秒 转换为DateTime
func DateClockToDateTime(year, month, day, hour, min, sec int, zone tz.TimeZone) (dt *DateTime, err error) {
	dt = &DateTime{zone: zone}
	err = dt.FlushToDateClock(year, month, day, hour, min, sec)
	if err != nil {
		return nil, err
	}
	return dt, nil
}

//根据年月加减月数获得新的年月
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

//开始时间和结束时间相差月份
func DateTimeDiffMonth(start, end *DateTime) int {
	yearCha := end.Year() - start.Year()
	return yearCha*12 + (end.Month() - start.Month())
}
