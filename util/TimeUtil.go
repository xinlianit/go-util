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

// 时间工具库
type Time struct {
}

// 获取当前毫秒时间
func (u Time) GetCurrentMilliTime() int64 {
	return time.Now().UnixNano() / 1e6
}

// 当前时间
// @param layout 时间格式
func (u Time) GetCurrentDateTime(layout string) string {
	return time.Now().Format(layout)
}

// 字符串时间转时间戳（秒）
// @param dateTime 日期时间
// @param layout 时间格式
func (u Time) StringTimeToSecond(dateTime string, layout string) (int64, error) {
	time, err := time.ParseInLocation(layout, dateTime, time.Local)
	if err != nil {
		return 0, err
	}

	return time.Unix(), nil
}
