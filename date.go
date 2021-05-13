package common

import "time"

/**
* 时间格式化
 */
const (
	Date        = "2006-01-02"
	shortdate   = "06-01-02"
	times       = "15:04:05"
	shorttime   = "15:04"
	Datetime    = "2006-01-02 15:04:05"
	Datetime2   = "2006-01-02T15:04:05Z07:00"
	newdatetime = "2006/01/02 15~04~05"
	newtime     = "15~04~05"
	Nanosecond  = 1
	Microsecond = 1000 * Nanosecond
	Millisecond = 1000 * Microsecond
	Second      = 1000 * Millisecond
	Minute      = 60 * Second
	Hour        = 60 * Minute
)

// Timestamp2Time 时间戳转字符串
func Timestamp2Time(timestamp int64, format string) string {
	return time.Unix(timestamp, 0).Format(format)
}

// Timestamp2Time2 时间戳转字符串
func Timestamp2Time2(timestamp int64) string {
	return time.Unix(timestamp, 0).Format(Datetime)
}

// Time2Timestamp 字符串转时间戳
func Time2Timestamp(t string) int64 {
	tt, _ := time.Parse(Datetime, t)
	return tt.Unix()
}

// Years 获取年限数组
func Years(start, end int) []int {
	var years []int
	for i := end; i >= start; i-- {
		years = append(years, i)
	}
	return years
}
