package models

import (
	"dgServer/utils"
	"fmt"
	"time"
)

var Months = []string{"1月", "2月", "3月", "4月", "5月", "6月", "7月", "8月", "9月", "10月", "11月", "12月"}
var WeekDays = []string{"周一", "周二", "周三", "周四", "周五", "周六", "周日"}
var MonthWeeks = []string{"第一周", "第二周", "第三周", "第四周", "第五周", "第六周"}

type SimpleUser struct {
	S_User_id int 
	S_name string
}

type WorkTimes struct {
	W_user_name string
	W_course string 
	W_mission string 
	W_num int
}


//UserStatus 用户状态
type UserStatus struct {
	M_id          int
	M_name        string
	M_day         string
	M_course      string
	M_last_submit string
	M_last_name   string
	M_num         int
	M_mission     string
	M_leave       int    //是否请假中,0 未请假 1.请假
	M_teacher     string //点评人
}

type WorkData struct {
	M_id     int
	M_name   string
	M_tag    string
	M_course string
	M_create string
	M_times  int
	M_work   string
}

//TaskStatus 任务状态
type TaskStatus struct {
	T_id           int
	T_name         string
	T_course_num   int
	T_usernames    map[string]int
	T_pass_num     int
	T_pass_percent int

	T_pass_data   map[int]int
	T_course_data map[int]int
}

type TaskScore struct {
	SD_utime       string
	SD_id          int
	SD_username    string
	SD_course_name string
	SD_name        string
	SD_stime       string
	SD_rank        int
}

type PayInfo struct {
	P_id       int
	P_name     string
	P_delete   int
	P_mail     string
	P_phone    string
	P_buy_time string
	P_end_time string

	P_last_submit string
	P_last_name   string

	EndTime     int64
	ReleaseTime string

	P_aliPay string

	PayTimes *PayTimes

	P_addDay int //添加天数，特殊需求
}

type OrderRecord struct {
	OR_id          int
	OR_name        string
	OR_user_type   string
	OR_course      string
	OR_tag         string
	OR_type        int
	OR_price       int
	OR_create_time string
	OR_start_time  string
	OR_aliPay      string
}

//LazeUserRecord 未缴费用户记录
type LazeUserRecord struct {
	LR_id int
	LR_name string
	LR_user_type string
	LR_date string
	LR_last string
	LR_leave string
	LR_course string
	LR_times int
	LR_work string
}

//PayTimes 时间
type PayTimes struct {
	StartTime int64
	EndTime   int64
	Next      *PayTimes
}

//SendTaskStatus 任务状态
type SendTaskStatus struct {
	T_id           int
	T_name         string
	T_course_num   int
	T_usernames    string
	T_pass_num     int
	T_pass_percent int
}

//SendTaskScore 任务评分状态
type SendTaskScore struct {
	SD_utime       string
	SD_id          int
	SD_username    string
	SD_course_name string
	SD_name        string
	SD_stime       string
	SD_rank        int
}

//RepayInfo 复购状态
type RepayInfo struct {
	P_id      int
	P_name    string
	P_type    int
	P_price   int
	P_create  string
	P_start   string
	P_aliPay  string
	P_E_start string
	P_course  string
	P_tag     string
}

//TeacherInfo 点评师信息
type TeacherInfo struct {
	T_createDate   string //日期
	T_createTime   string //时间
	T_month        string //月份
	T_monthWeek    string
	T_weekDay      string
	T_name         string
	T_level        int
	T_time         string
	T_user_name    string //用户名
	T_course_name  string //通道名字
	T_mission_name string //任务名字
	T_avg_time     string //平均时间

	T_delete int //是否归档
}

type CourseData struct {
	Month int
	Num   int
}

type UserDetail struct {
	ID          int
	CourseName  string
	TotalNum    int
	CourseDatas []*CourseData
}

func CreateMonthCourseData() []*CourseData {
	returnDatas := make([]*CourseData, 12)

	for index := 0; index < 12; index++ {
		returnDatas[index] = &CourseData{
			Month: index + 1,
			Num:   0,
		}
	}

	return returnDatas
}

//GetSendData 发送数据
func (ts *TaskStatus) GetSendData() *SendTaskStatus {

	nameString := ""

	for k := range ts.T_usernames {
		if k != "" {
			if nameString == "" {
				nameString = k
			} else {
				nameString = fmt.Sprintf("%s,%s", nameString, k)
			}
		}
	}
	totalNum := 0
	totalPercent := 0
	for useID, v := range ts.T_pass_data {
		if total := ts.T_course_data[useID]; total > 0 {
			totalNum++
			totalPercent += (v * 100 / total)
		}
	}

	percent := 0

	if totalNum > 0 {
		percent = totalPercent / totalNum
	}

	return &SendTaskStatus{
		T_id:           ts.T_id,
		T_name:         ts.T_name,
		T_course_num:   ts.T_course_num,
		T_usernames:    nameString,
		T_pass_num:     ts.T_pass_num,
		T_pass_percent: percent,
	}
}

//GetSendData 发送数据
func (ts *TaskScore) GetSendData() *SendTaskScore {

	return &SendTaskScore{
		SD_utime:       ts.SD_utime,
		SD_id:          ts.SD_id,
		SD_username:    ts.SD_username,
		SD_course_name: ts.SD_course_name,
		SD_name:        ts.SD_name,
		SD_stime:       ts.SD_stime,
		SD_rank:        ts.SD_rank,
	}
}

//GetSendData 统计数据
func (pi *PayInfo) GetSendData() *PayInfo {

	var endTime int64
	var releaseTime int64

	var leaves []*Leave
	if LeaveMap[pi.P_id] != nil {
		leaves = append(leaves, LeaveMap[pi.P_id]...)
	}

	// if pi.P_id == 480 {
	// 	fmt.Println(pi)
	// }
	node := pi.PayTimes

	for node != nil {

		for k, v := range leaves {
			if v != nil {
				if node.StartTime <= v.StartDate && node.EndTime >= v.StartDate {
					node.EndTime += (v.EndDate - v.StartDate)
					leaves[k] = nil
					releaseTime += (v.EndDate - v.StartDate)
				}
			}
		}

		if pi.P_name == "海底捞" {
			fmt.Println("pi.PayTimes[i] = ", node)
		}
		if endTime == 0 {
			endTime = node.EndTime
		} else {

			if endTime > node.StartTime {
				endTime += (node.EndTime - node.StartTime)
			} else {
				releaseTime = 0
				endTime = node.EndTime
			}
		}
		node = node.Next
	}
	/*for index := 1; index <= len(pi.PayTimes); index++ {

		for k, v := range leaves {
			if v != nil {
				if pi.PayTimes[index].StartTime <= v.StartDate && pi.PayTimes[index].EndTime >= v.StartDate {
					pi.PayTimes[index].EndTime += (v.EndDate - v.StartDate)
					leaves[k] = nil
					releaseTime += (v.EndDate - v.StartDate)
				}
			}
		}

		if pi.P_name == "海底捞" {
			fmt.Println("pi.PayTimes[i] = ", pi.PayTimes[index])
		}
		if index == 1 {
			endTime = pi.PayTimes[1].EndTime
		} else {

			if endTime > pi.PayTimes[index].StartTime {
				endTime += (pi.PayTimes[index].EndTime - pi.PayTimes[index].StartTime)
			} else {
				releaseTime = 0
				endTime = pi.PayTimes[index].EndTime
			}
		}
	}*/

	if pi.P_name == "海底捞" {
		fmt.Println("endTime = ", endTime)
	}
	pi.EndTime = endTime
	pi.P_end_time = utils.GetTimeString(utils.ParseUnixTime(endTime))
	pi.ReleaseTime = fmt.Sprintf("%d天", releaseTime/86400)

	return pi
}

func (or *OrderRecord) GetSendData() *OrderRecord {
	return or
}

func GetAddDay(startTime, endTime int64) int {
	timeSet := time.Unix(startTime, 0)
	n := 0
	for endTime >= timeSet.AddDate(0, n, 0).Unix() && n < 12 {
		n++

		timeSet = timeSet.Add(time.Second * time.Duration(n))

	}

	return n * 4
}

//GetWeekDayString 获取周日期
func GetWeekDayString(weekNum int) string {
	if weekNum > 7 || weekNum <= 0 {
		weekNum = 7
	}

	return WeekDays[weekNum-1]
}

//GetMonthWeekString 获取周数
func GetMonthWeekString(monthWeek int) string {

	return MonthWeeks[monthWeek]
}

//GetMonthString 获取周数
func GetMonthString(month int) string {

	return Months[month-1]
}
