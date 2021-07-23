package util

import (
	"sync"
	"time"
)

var (
	timeUtilInstance *Time
	timeUtilOnce     sync.Once
)

func TimeUtil() *Time {
	timeUtilOnce.Do(func() {
		timeUtilInstance = new(Time)
	})
	return timeUtilInstance
}

// Time 时间工具库
type Time struct {
}

// GetCurrentSecond 获取当前时间戳(秒)
func (u Time) GetCurrentSecond() int64 {
	return time.Now().Unix()
}

// GetCurrentMilliTime 获取当前时间戳(毫秒)
func (u Time) GetCurrentMilliTime() int64 {
	return time.Now().UnixNano() / 1e6
}

// GetCurrentDateTime 获取当前日期时间
// @param layout 时间格式
func (u Time) GetCurrentDateTime(layout string) string {
	return time.Now().Format(layout)
}

// StringTimeToSecond 字符串时间转时间戳（秒）
// @param dateTime 日期时间
// @param layout 时间格式
func (u Time) StringTimeToSecond(dateTime string, layout string) (int64, error) {
	time, err := time.ParseInLocation(layout, dateTime, time.Local)
	if err != nil {
		return 0, err
	}

	return time.Unix(), nil
}

// MillisecondToDateTime 时间戳(毫秒)转日期时间
// @param millisecond 毫秒时间戳
// @param layout 时间格式
// @param string 日期时间
func (u Time) MillisecondToDateTime (millisecond int64, layout string) string {
	return time.Unix(0, millisecond * int64(time.Millisecond)).Format(layout)
}

// SecondToDateTime 时间戳(秒)转日期时间
// @param second 秒时间戳
// @param layout 时间格式
// @param string 日期时间
func (u Time) SecondToDateTime (second int64, layout string) string {
	return time.Unix(second, 0).Format(layout)
}