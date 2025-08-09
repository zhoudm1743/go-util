package types

import (
	"fmt"
	"time"
)

type XTime struct {
	t time.Time
}

// 常用时间格式常量
const (
	DateTimeFormat = "2006-01-02 15:04:05"
	DateFormat     = "2006-01-02"
	TimeFormat     = "15:04:05"
	ISO8601Format  = "2006-01-02T15:04:05Z07:00"
	RFC3339Format  = time.RFC3339
	UnixFormat     = "Mon Jan _2 15:04:05 MST 2006"
)

// Time 创建 XTime 实例
func Time(t time.Time) XTime {
	return XTime{t: t}
}

// Now 获取当前时间
func Now() XTime {
	return XTime{t: time.Now()}
}

// Date 创建指定日期时间
func Date(year, month, day, hour, min, sec int) XTime {
	return XTime{t: time.Date(year, time.Month(month), day, hour, min, sec, 0, time.Local)}
}

// ParseTime 解析时间字符串
func ParseTime(layout, value string) (XTime, error) {
	t, err := time.Parse(layout, value)
	if err != nil {
		return XTime{}, err
	}
	return XTime{t: t}, nil
}

// ParseDateTime 解析日期时间字符串 (YYYY-MM-DD HH:MM:SS)
func ParseDateTime(value string) (XTime, error) {
	return ParseTime(DateTimeFormat, value)
}

// ParseDate 解析日期字符串 (YYYY-MM-DD)
func ParseDate(value string) (XTime, error) {
	return ParseTime(DateFormat, value)
}

// FromUnix 从 Unix 时间戳创建
func FromUnix(sec int64) XTime {
	return XTime{t: time.Unix(sec, 0)}
}

// FromUnixMilli 从毫秒时间戳创建
func FromUnixMilli(msec int64) XTime {
	return XTime{t: time.Unix(msec/1000, (msec%1000)*1000000)}
}

// Time 获取原始 time.Time
func (x XTime) Time() time.Time {
	return x.t
}

// Unix 获取 Unix 时间戳
func (x XTime) Unix() int64 {
	return x.t.Unix()
}

// UnixMilli 获取毫秒时间戳
func (x XTime) UnixMilli() int64 {
	return x.t.UnixMilli()
}

// UnixNano 获取纳秒时间戳
func (x XTime) UnixNano() int64 {
	return x.t.UnixNano()
}

// Year 获取年份
func (x XTime) Year() int {
	return x.t.Year()
}

// Month 获取月份
func (x XTime) Month() int {
	return int(x.t.Month())
}

// Day 获取日期
func (x XTime) Day() int {
	return x.t.Day()
}

// Hour 获取小时
func (x XTime) Hour() int {
	return x.t.Hour()
}

// Minute 获取分钟
func (x XTime) Minute() int {
	return x.t.Minute()
}

// Second 获取秒
func (x XTime) Second() int {
	return x.t.Second()
}

// Weekday 获取星期几 (0=Sunday, 1=Monday, ...)
func (x XTime) Weekday() int {
	return int(x.t.Weekday())
}

// WeekdayName 获取星期几的名称
func (x XTime) WeekdayName() string {
	return x.t.Weekday().String()
}

// IsZero 判断是否为零时间
func (x XTime) IsZero() bool {
	return x.t.IsZero()
}

// Format 格式化时间
func (x XTime) Format(layout string) string {
	return x.t.Format(layout)
}

// FormatDateTime 格式化为日期时间字符串
func (x XTime) FormatDateTime() string {
	return x.t.Format(DateTimeFormat)
}

// FormatDate 格式化为日期字符串
func (x XTime) FormatDate() string {
	return x.t.Format(DateFormat)
}

// FormatTime 格式化为时间字符串
func (x XTime) FormatTime() string {
	return x.t.Format(TimeFormat)
}

// FormatISO8601 格式化为 ISO8601 格式
func (x XTime) FormatISO8601() string {
	return x.t.Format(ISO8601Format)
}

// String 实现 Stringer 接口
func (x XTime) String() string {
	return x.FormatDateTime()
}

// Add 添加时间间隔
func (x XTime) Add(d time.Duration) XTime {
	return XTime{t: x.t.Add(d)}
}

// AddYears 添加年数
func (x XTime) AddYears(years int) XTime {
	return XTime{t: x.t.AddDate(years, 0, 0)}
}

// AddMonths 添加月数
func (x XTime) AddMonths(months int) XTime {
	return XTime{t: x.t.AddDate(0, months, 0)}
}

// AddDays 添加天数
func (x XTime) AddDays(days int) XTime {
	return XTime{t: x.t.AddDate(0, 0, days)}
}

// AddHours 添加小时数
func (x XTime) AddHours(hours int) XTime {
	return XTime{t: x.t.Add(time.Duration(hours) * time.Hour)}
}

// AddMinutes 添加分钟数
func (x XTime) AddMinutes(minutes int) XTime {
	return XTime{t: x.t.Add(time.Duration(minutes) * time.Minute)}
}

// AddSeconds 添加秒数
func (x XTime) AddSeconds(seconds int) XTime {
	return XTime{t: x.t.Add(time.Duration(seconds) * time.Second)}
}

// Sub 计算时间差
func (x XTime) Sub(other XTime) time.Duration {
	return x.t.Sub(other.t)
}

// Before 判断是否在指定时间之前
func (x XTime) Before(other XTime) bool {
	return x.t.Before(other.t)
}

// After 判断是否在指定时间之后
func (x XTime) After(other XTime) bool {
	return x.t.After(other.t)
}

// Equal 判断时间是否相等
func (x XTime) Equal(other XTime) bool {
	return x.t.Equal(other.t)
}

// StartOfDay 获取当天开始时间 (00:00:00)
func (x XTime) StartOfDay() XTime {
	year, month, day := x.t.Date()
	return XTime{t: time.Date(year, month, day, 0, 0, 0, 0, x.t.Location())}
}

// EndOfDay 获取当天结束时间 (23:59:59)
func (x XTime) EndOfDay() XTime {
	year, month, day := x.t.Date()
	return XTime{t: time.Date(year, month, day, 23, 59, 59, 999999999, x.t.Location())}
}

// StartOfWeek 获取本周开始时间 (周一)
func (x XTime) StartOfWeek() XTime {
	weekday := x.t.Weekday()
	if weekday == time.Sunday {
		weekday = 7
	}
	days := int(weekday) - 1
	return x.AddDays(-days).StartOfDay()
}

// EndOfWeek 获取本周结束时间 (周日)
func (x XTime) EndOfWeek() XTime {
	return x.StartOfWeek().AddDays(6).EndOfDay()
}

// StartOfMonth 获取本月开始时间
func (x XTime) StartOfMonth() XTime {
	year, month, _ := x.t.Date()
	return XTime{t: time.Date(year, month, 1, 0, 0, 0, 0, x.t.Location())}
}

// EndOfMonth 获取本月结束时间
func (x XTime) EndOfMonth() XTime {
	return x.StartOfMonth().AddMonths(1).AddDays(-1).EndOfDay()
}

// StartOfYear 获取本年开始时间
func (x XTime) StartOfYear() XTime {
	year := x.t.Year()
	return XTime{t: time.Date(year, 1, 1, 0, 0, 0, 0, x.t.Location())}
}

// EndOfYear 获取本年结束时间
func (x XTime) EndOfYear() XTime {
	year := x.t.Year()
	return XTime{t: time.Date(year, 12, 31, 23, 59, 59, 999999999, x.t.Location())}
}

// Age 计算年龄（基于当前时间）
func (x XTime) Age() int {
	return x.AgeAt(Now())
}

// AgeAt 计算在指定时间点的年龄
func (x XTime) AgeAt(at XTime) int {
	years := at.Year() - x.Year()
	if at.Month() < x.Month() || (at.Month() == x.Month() && at.Day() < x.Day()) {
		years--
	}
	return years
}

// DaysTo 计算到指定时间的天数
func (x XTime) DaysTo(other XTime) int {
	duration := other.StartOfDay().Sub(x.StartOfDay())
	return int(duration.Hours() / 24)
}

// IsLeapYear 判断是否为闰年
func (x XTime) IsLeapYear() bool {
	year := x.Year()
	return year%4 == 0 && (year%100 != 0 || year%400 == 0)
}

// Quarter 获取季度 (1-4)
func (x XTime) Quarter() int {
	month := x.Month()
	return (month-1)/3 + 1
}

// DayOfYear 获取一年中的第几天
func (x XTime) DayOfYear() int {
	return x.t.YearDay()
}

// WeekOfYear 获取一年中的第几周
func (x XTime) WeekOfYear() int {
	_, week := x.t.ISOWeek()
	return week
}

// FormatRelative 格式化为相对时间（如：2小时前）
func (x XTime) FormatRelative() string {
	return x.FormatRelativeAt(Now())
}

// FormatRelativeAt 基于指定时间格式化为相对时间
func (x XTime) FormatRelativeAt(at XTime) string {
	duration := at.Sub(x)

	if duration < 0 {
		duration = -duration
		return x.formatDuration(duration) + "后"
	}

	return x.formatDuration(duration) + "前"
}

// formatDuration 格式化时间间隔
func (x XTime) formatDuration(d time.Duration) string {
	if d < time.Minute {
		return "刚刚"
	}

	if d < time.Hour {
		minutes := int(d.Minutes())
		return fmt.Sprintf("%d分钟", minutes)
	}

	if d < 24*time.Hour {
		hours := int(d.Hours())
		return fmt.Sprintf("%d小时", hours)
	}

	if d < 30*24*time.Hour {
		days := int(d.Hours() / 24)
		return fmt.Sprintf("%d天", days)
	}

	if d < 365*24*time.Hour {
		months := int(d.Hours() / (24 * 30))
		return fmt.Sprintf("%d个月", months)
	}

	years := int(d.Hours() / (24 * 365))
	return fmt.Sprintf("%d年", years)
}

// In 转换时区
func (x XTime) In(loc *time.Location) XTime {
	return XTime{t: x.t.In(loc)}
}

// InUTC 转换为 UTC 时区
func (x XTime) InUTC() XTime {
	return XTime{t: x.t.UTC()}
}

// Location 获取时区
func (x XTime) Location() *time.Location {
	return x.t.Location()
}

// Between 判断时间是否在指定范围内
func (x XTime) Between(start, end XTime) bool {
	return (x.After(start) || x.Equal(start)) && (x.Before(end) || x.Equal(end))
}

// FormatChinese 格式化为中文日期格式
func (x XTime) FormatChinese() string {
	year := x.Year()
	month := x.Month()
	day := x.Day()
	hour := x.Hour()
	minute := x.Minute()
	second := x.Second()

	return fmt.Sprintf("%d年%d月%d日 %02d时%02d分%02d秒", year, month, day, hour, minute, second)
}

// TimezoneName 获取时区名称
func (x XTime) TimezoneName() string {
	name, _ := x.t.Zone()
	return name
}

// TimezoneOffset 获取时区偏移量（秒）
func (x XTime) TimezoneOffset() int {
	_, offset := x.t.Zone()
	return offset
}
