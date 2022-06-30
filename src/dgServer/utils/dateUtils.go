package utils

import "time"

//ParseStringTime 解析时间
func ParseStringTime(timeValue string) (time.Time, error) {
	parseTime, err := time.ParseInLocation("2006-01-02 15:04:05", timeValue, time.Local)

	if err != nil {
		return time.Now(), err
	}

	return parseTime, nil
}

//ParseStringDayTypeTime 解析时间
func ParseStringDayTypeTime(timeValue string) (time.Time, error) {
	parseTime, err := time.ParseInLocation("2006-01-02", timeValue, time.Local)

	if err != nil {
		return time.Now(), err
	}

	return parseTime, nil
}

//ParseStringDayTime 解析时间
func ParseStringDayTime(timeValue string) (time.Time, error) {
	parseTime, err := time.ParseInLocation("2006/01/02", timeValue, time.Local)

	if err != nil {
		return time.Now(), err
	}

	return parseTime, nil
}

//ParseUnixTime 解析时间戳
func ParseUnixTime(timeValue int64) time.Time {
	parseTime := time.Unix(timeValue, 0)

	return parseTime
}

//GetTimeString 获取当前时间的string
func GetTimeString(t time.Time) string {
	return t.Format("2006/01/02 15:04:05")
}

//GetTimeDay 获取当前时间的string
func GetTimeDay(t time.Time) string {
	return t.Format("2006/01/02")
}

//GetMonthWeek 获取月份的周数
func GetMonthWeek(t time.Time) (int, int) {
	tMonth := time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, time.Local)

	weekDay := int(tMonth.Weekday())
	if weekDay == 0 {
		weekDay = 7
	}
	d := int(t.Sub(tMonth).Hours() / 24)

	diffDay := d - (7 - weekDay)
	weekNum := diffDay / 7
	monthNum := int(t.Month())

	if weekNum < 0 {
		weekNum = 0
	}

	if diffDay > 0 && weekDay != 1 {
		weekNum++
		monthNum = int(t.AddDate(0, -1, 0).Month())
	}

	if weekDay > 1 {
		if weekNum <= 0 {
			tLastMonth := tMonth.AddDate(0, 0, -1)
			return GetMonthWeek(tLastMonth)
		}
	}

	if weekDay != 1 {
		monthNum = int(t.Month())
		weekNum--
	}

	if weekNum > 5 {
		weekNum = 5
	}

	return monthNum, weekNum
}

//GetCurrentYear 获取当前年份
func GetCurrentYear() int {
	return time.Now().Year()
}

type yearData struct {
	Index int
	Value int
}

//GetStartYearList 获取年份列表
func GetStartYearList() []*yearData {
	cYear := GetCurrentYear()
	var years []*yearData
	for index := cYear; index >= 2017; index-- {
		years = append(years, &yearData{
			Index: cYear - index,
			Value: index,
		})
	}
	return years
}
