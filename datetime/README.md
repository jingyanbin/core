# datetime
日期时间模块

部分用法示例:

	//秒级时间戳
	unix := datetime.Unix()

	//毫秒级时间戳
	datetime.UnixMs()

	//格式化日志时间字符串转 DateTime
	datetime.FormatToDateTime("2020-09-12 00:00:00", "%Y-%m-%d %H:%M:%S", datetime.Zones.LOCAL, false)

	//标准日期时间字符串 转换为 DateTime
	datetime.YmdHMSToDateTime("2020/09/12 00:00:00",datetime.Zones.LOCAL, false)

	//标准日期时间字符串 转换为 秒级时间戳
	datetime.YmdHMSToUnix("2020/09/12 00:00:00",datetime.Zones.LOCAL, false)

	//格式化日志时间字符串转 秒级时间戳
	datetime.FormatToUnix("2020-09-12 00:00:00", "%Y-%m-%d %H:%M:%S", datetime.Zones.LOCAL, false)

	//年月日时分秒转换为 DateTime
	datetime.DateClockToDateTime(2020, 9,12,1,1,1,datetime.Zones.LOCAL)


	//当前时间的DateTime
	dt := datetime.Now()

	//刷新为最新
	dt.Flush()

	//刷新到指定时间戳
	dt.FlushToUnix(unix)

	//刷新到指定标准时间字符串
	dt.FlushToYmdHMS("2020/09/11 00:00:00", false)

	//刷新到指定格式化时间字符串
	dt.FlushToFormat("2020-09-12 00:00:00", "%Y-%m-%d %H:%M:%S", false)

	//刷新到指定 年月日时分秒
	dt.FlushToDateClock(2020, 10, 23, 1, 2,3)

	//得到格式化日期时间字符串
	dt.Format("%Y-%m-%d %H:%M:%S")

	//得到标准日期时间字符串
	dt.YmdHMS()

	//星期1-7
	dt.WeekdayA()

	//秒级时间戳
	dt.Unix()

	//下一个周1零点时间戳
	dt.UnixFutureWeekDayA(1, 0, 0, 0)

	//下一周的周1零点时间戳
	dt.UnixNextWeekDayA(1, 0,0,0)

	//一年第一天0点时间戳
	dt.UnixYearZeroHour()

	//这个月1号0点时间戳
	dt.UnixMonthZeroHour()

	//今天0点时间戳
	dt.UnixDayZeroHour()

	//1970年1月1日，以来 过了好多天
	dt.UnixDayNumber()

	//年
	dt.Year()

	//月
	dt.Month()

	//日
	dt.Day()

	//时
	dt.Hour()

	//分
	dt.Min()

	//秒
	dt.Sec()

	//一天中第几秒
	dt.DaySecond()

	//设置时区
	dt.SetZone(datetime.Zones.E8)
