package util

// ConvertSeconds 将秒时间戳转换为天、小时、分钟和秒的格式
func ConvertSeconds(seconds int64) (days, hours, minutes, secs int64) {
	minutes = seconds / 60
	hours = minutes / 60
	days = hours / 24

	secs = seconds % 60
	minutes = minutes % 60
	hours = hours % 24

	return
}
