package utils

import (
	"fmt"
	"math"
	"strings"
	"time"
)

const timeLayout = "2006-01-02 15:04:05"

// 当前毫秒数
func TimeMillisecond() int64 {
	return time.Now().UnixNano() / 1e6
}

// 当前时间戳
func TimeSecond() int64 {
	return time.Now().UnixNano() / 1e9
}

// 当天零点时间戳
func DayTimeZero() int64 {
	t := time.Now()
	_, offset := t.Zone()
	return t.Unix() - (t.Unix()+int64(offset))%86400
}

func DayTimeZero2() int64 {
	t := time.Now()
	newTime := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
	return newTime.Unix()
}

// 当周零点时间戳(周一作为一周第一天)
func WeekStartEnd() (int64, int64) {
	t := time.Now()
	newTime := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
	start := newTime.Unix() - int64((newTime.Weekday()-1)*86400)
	end := start + 7*86400
	return start, end
}

// 当月开始时间戳
func MonthStartEnd() (int64, int64) {
	t := time.Now()
	startDay := t.AddDate(0, 0, -t.Day()+1)
	startTime := time.Date(startDay.Year(), startDay.Month(), startDay.Day(), 0, 0, 0, 0, t.Location())
	endDay := startTime.AddDate(0, 1, -1)
	endTime := time.Date(endDay.Year(), endDay.Month(), endDay.Day(), 23, 59, 59, 0, t.Location())
	return startTime.Unix(), endTime.Unix()
}

// 两个时间戳相差天数
func GetDiffDays(t1 int64, t2 int64) int {
	_, offset := time.Now().Zone()
	zeroTime1 := t1 - (t1+int64(offset))%86400
	zeroTime2 := t2 - (t2+int64(offset))%86400
	return int(math.Abs(float64(zeroTime1-zeroTime2)) / 86400)
}

// 判断两个时间是否是同一天
func IsSameDay(t1 int64, t2 int64) bool {
	return GetDiffDays(t1, t2) == 0
}

func formatTimeStr(timeStr string) string {
	tmpSlices := strings.Split(timeStr, " ")
	if len(tmpSlices) != 2 {
		return ""
	}
	var dateSlice []string
	if strings.ContainsAny(tmpSlices[0], "/") {
		dateSlice = strings.Split(tmpSlices[0], "/")
	} else if strings.ContainsAny(tmpSlices[0], "-") {
		dateSlice = strings.Split(tmpSlices[0], "-")
	}
	if len(dateSlice) == 3 {
		return fmt.Sprintf("%04s-%02s-%02s %s", dateSlice[0], dateSlice[1], dateSlice[2], tmpSlices[1])
	}
	return ""
}

// 时间字符串-》时间对象
func Str2Time(timeStr string) (time.Time, error) {
	timeStr = formatTimeStr(timeStr)
	local, err := time.LoadLocation("Local")
	var t time.Time
	theTime, err := time.ParseInLocation(timeLayout, timeStr, local)
	if err != nil {
		return t, err
	}
	return theTime, nil
}

// 时间字符串-》时间戳
func Str2UnixTime(timeStr string) (int64, error) {
	t, err := Str2Time(timeStr)
	if err != nil {
		return 0, err
	}
	return t.Unix(), nil
}

// 时间戳-》时间字符串
func Time2Str(t int64) string {
	return time.Unix(t, 0).Format(timeLayout)
}

// 时间戳-》时间对象
func Time2TimeObj(t int64) time.Time {
	return time.Unix(t, 0).Local()
}

func YeadAndWeek() (y, w int) {
	return time.Now().ISOWeek()
}
