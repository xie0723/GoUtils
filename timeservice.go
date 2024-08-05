package GoUtils

import "time"

// NowFmtStr 返回格式化的时间
func NowFmtStr() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

// TimeDiff 2个字符串时间差值
func TimeDiff(start, end string) int64 {
	startTime, _ := time.ParseInLocation("2006-01-02 15:04:05", start, time.Local)
	endTime, _ := time.ParseInLocation("2006-01-02 15:04:05", end, time.Local)
	return endTime.Unix() - startTime.Unix()
}

// TimeStampS 获取时间戳-秒
func TimeStampS() int64 {
	return time.Now().Unix()
}

// TimeStampMs 获取时间戳-毫秒
func TimeStampMs() int64 {
	return time.Now().UnixNano() / 1e6
}

// TimeStamp2FmtStr 时间戳转为格式化时间
func TimeStamp2FmtStr(timeStamp int64) string {
	timeObj := time.Unix(timeStamp, 0)
	return timeObj.Format("2006-01-02 15:04:05")
}

// FmtStr2TimeStamp 格式化时间转为时间戳
func FmtStr2TimeStamp(timeStr string) int64 {
	timeObj, _ := time.ParseInLocation("2006-01-02 15:04:05", timeStr, time.Local)
	return timeObj.Unix()
}

// Utc2Local utc时间转为本地时间
func Utc2Local(utcTimeStr string) string {
	utcTime, _ := time.ParseInLocation("2006-01-02 15:04:05", utcTimeStr, time.UTC)
	localTime := utcTime.Add(8 * time.Hour)
	return localTime.Format("2006-01-02 15:04:05")
}