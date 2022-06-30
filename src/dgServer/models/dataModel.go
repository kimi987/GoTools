package models

import (
	"dgServer/db"
	"dgServer/utils"
	"fmt"
	"strconv"
	"strings"
	"time"
)

var ExperCourses map[int]int //[]int{22, 23, 24, 90}

//RefreshTime 刷新时间
var RefreshTime int

//User 用户
type User struct {
	ID              int
	Nickname        string
	Mobile          string
	Email           string
	Alipay_nickname string
	Privilege       int //用户等级 10 普通用户 20 体验用户 30 付费用户
	Delete          string
	Created_at 		int64 
}

//UserMap 用户表
var UserMap map[int]map[int]*User
var UserAllMap map[int]*User
var UserWithDelete map[int]*User
var TeacherMaps map[int]*User
var StudentMap map[int]*User

type UserTag struct {
	ID   int
	Name string
}

var UserTagMap map[int]*UserTag

//UserTagRel 玩家标签表
type UserTagRel struct {
	UserID  int
	TagID   int
	Created int64
}

//UserTagRelMap 玩家标签表
var UserTagRelMap map[int][]*UserTagRel

//UserCourse 用户通道
type UserCourse struct {
	ID       int
	UserID   int
	CourseID int
}

//UserCourseMap 用户课程表
var UserCourseMap map[int]map[int]*UserCourse
var CourseUserMap map[int]map[int]*UserCourse

//Course 课程
type Course struct {
	ID    int
	Title string
}

//CourseMap 课程表
var CourseMap map[int]*Course

//Mission 任务
type Mission struct {
	ID    int
	Title string
}

//MissionMap 任务表
var MissionMap map[int]*Mission

//Work 作业
type Work struct {
	ID         int
	CourseID   int
	ChapterID  int
	MissionID  int
	UserID     int
	Images     string
	CreateTime int64
	UpdateTime int64
}

//WorkMap 作业表
var WorkMap map[int]map[int]*Work
var WorkCourseMap map[int]map[int]map[int]*Work
var WorkAllMap map[int]*Work

//UserExperience 用户体验
type UserExperience struct {
	UserID   int
	CourseID int
	Started  int64
}

//UserExperienceMap 用户体验表
var UserExperienceMap map[int]*UserExperience

var UserExperienceAllMap map[int]*UserExperience

//Order 订单
type Order struct {
	SummaryID int
	Period    int //12 年卡 6半年卡 3 季卡
	Price     int //价格
	Quantity  int //数量
	Created   int64
}

//OrderMap 订单表
var OrderMap map[int]*Order
var OrderTotalMap map[int]*Order

//OrderSummary 购买记录
type OrderSummary struct {
	ID         int
	UserID     int
	Created_at int64
}

//OrderSummaryMap 购买记录表
var OrderSummaryMap map[int]map[int]*OrderSummary

//Rating 评价结构
type Rating struct {
	ID             int
	WorkID         int
	Rank           int
	UpdateTime     int64
	MissionID      int //任务ID
	CourseID       int //通道ID
	UserID         int //点评师ID
	RatingDuration int //点评时间
}

//RatingMap 评价表
var RatingMap map[int]*Rating

//RatingUserMap 点评师课程分类
var RatingUserMap map[int]map[int]map[int]*Rating

//RatingCourseMap 课程分类
var RatingCourseMap map[int]map[int]*Rating

//Leave 请假表
type Leave struct {
	UserID    int
	StartDate int64
	EndDate   int64
}

//LeaveMap 请假表
var LeaveMap map[int][]*Leave

//HuanBi 环比表
type HuanBi struct {
	UserID int
}
//HuanBiMap 环比表map
var HuanBiMap map[int]*HuanBi

func init() {

	go func() {
		for {
			if RefreshTime > 0 {
				RefreshTime--
			}
			time.Sleep(time.Minute)
		}
	}()
}

var isLoadData = false

func LoadDataAsync() {
	if isLoadData {
		return
	}
	go func() {
		isLoadData = true
		LoadData()
		isLoadData = false
	}()
}

//LoadData 读取数据
func LoadData() error {
	if RefreshTime <= 0 {
		var err error
		UserMap, err = LoadUser()

		if err != nil {
			return err
		}

		UserCourseMap, err = LoadUserCourse()

		if err != nil {
			return err
		}

		CourseMap, err = LoadCourse()

		if err != nil {
			return err
		}

		MissionMap, err = LoadMission()

		if err != nil {
			return err
		}

		WorkMap, err = LoadWork()

		if err != nil {
			return err
		}

		UserExperienceMap, err = LoadUserExperiences()

		if err != nil {
			return err
		}

		OrderMap, err = LoadOrder()

		if err != nil {
			return err
		}

		OrderSummaryMap, err = LoadOrderSummary()

		if err != nil {
			return err
		}

		RatingMap, err = LoadRating()

		if err != nil {
			return err
		}

		LeaveMap, err = LoadLeaves()

		if err != nil {
			return err
		}

		HuanBiMap, err = LoadHuanBi()

		if err != nil {
			return err
		}


		RefreshTime = 15
	}
	return nil
}

//LoadUser 读取用户
func LoadUser() (map[int]map[int]*User, error) {

	userTags, err := LoadUserTag()

	if err != nil {
		return nil, err
	}

	UserTagMap = userTags

	tags, err := LoadUserTagRel()

	if err != nil {
		return nil, err
	}

	UserTagRelMap = tags

	req := db.CreateRequestDB("Select id,nickname,mobile,email,privilege,created_at, deleted_at,alipay_nickname from user ", nil)
	res := db.CreateResponseDB()
	err = db.QuerySelect(req, res)

	if err != nil {
		return nil, err
	}
	users := make(map[int]map[int]*User)
	userAllMap := make(map[int]*User)
	studentMap := make(map[int]*User)
	userWithDelete := make(map[int]*User)
	if res.Rows != nil {
		for _, v := range res.Rows {
			if v != nil {
				id, _ := strconv.Atoi(v["id"])
				nickname := v["nickname"]
				mobile := v["mobile"]
				email := v["email"]
				privilege, _ := strconv.Atoi(v["privilege"])
				created_at, _ := utils.ParseStringTime(v["created_at"])
				delete_at := v["deleted_at"]
				alipay_nickname := v["alipay_nickname"]
				
				user := &User{ID: id, Nickname: nickname, Mobile: mobile, Email: email, Privilege: privilege, Alipay_nickname: alipay_nickname, Delete: delete_at}
				user.Created_at = created_at.Unix()
				userAllMap[user.ID] = user
				isPass := false
				if UserTagRelMap[id] != nil {
					for _, v := range UserTagRelMap[id] {
						if v.TagID == 4 || v.TagID == 13 {
							isPass = true
							break
						}
					}
				}

				if isPass {
					continue
				}

				userWithDelete[user.ID] = user

				if delete_at != "" {
					continue
				}

				if users[privilege] == nil {
					users[privilege] = make(map[int]*User)
				}
				if privilege != 50 {
					studentMap[id] = user
				}

				users[privilege][id] = user
			}
		}
	}

	UserAllMap = userAllMap
	UserWithDelete = userWithDelete
	StudentMap = studentMap

	return users, nil
}

func LoadUserTag() (map[int]*UserTag, error) {
	req := db.CreateRequestDB("Select id,name from user_tag", nil)
	res := db.CreateResponseDB()
	err := db.QuerySelect(req, res)

	if err != nil {
		return nil, err
	}
	userTags := make(map[int]*UserTag)
	if res.Rows != nil {
		for _, v := range res.Rows {
			if v != nil {
				id, _ := strconv.Atoi(v["id"])
				name := v["name"]

				userTag := &UserTag{ID: id, Name: name}
				userTags[id] = userTag
			}
		}
	}

	return userTags, nil
}

//LoadUserTagRel 读取玩家标签
func LoadUserTagRel() (map[int][]*UserTagRel, error) {
	req := db.CreateRequestDB("Select user_id,tag_id,created_at from user_tag_rel", nil)
	res := db.CreateResponseDB()
	err := db.QuerySelect(req, res)

	if err != nil {
		return nil, err
	}
	userTagRels := make(map[int][]*UserTagRel)
	if res.Rows != nil {
		for _, v := range res.Rows {
			if v != nil {
				userID, _ := strconv.Atoi(v["user_id"])
				tagID, _ := strconv.Atoi(v["tag_id"])

				createTime := v["created_at"]

				cTime, _ := utils.ParseStringTime(createTime)

				userTagRel := &UserTagRel{UserID: userID, TagID: tagID, Created: cTime.Unix()}
				userTagRels[userID] = append(userTagRels[userID], userTagRel)
			}
		}
	}

	return userTagRels, nil
}

//LoadUserCourse 读取用户课程
func LoadUserCourse() (map[int]map[int]*UserCourse, error) {
	req := db.CreateRequestDB("Select id,user_id,course_id from user_course", nil)
	res := db.CreateResponseDB()
	err := db.QuerySelect(req, res)

	if err != nil {
		return nil, err
	}
	userCourses := make(map[int]map[int]*UserCourse)
	courseUserMap := make(map[int]map[int]*UserCourse)
	if res.Rows != nil {
		for _, v := range res.Rows {
			if v != nil {
				id, _ := strconv.Atoi(v["id"])
				userID, _ := strconv.Atoi(v["user_id"])
				courseID, _ := strconv.Atoi(v["course_id"])
				userCourse := &UserCourse{ID: id, UserID: userID, CourseID: courseID}
				if userCourses[userID] == nil {
					userCourses[userID] = make(map[int]*UserCourse)
				}

				userCourses[userID][courseID] = userCourse

				if courseUserMap[courseID] == nil {
					courseUserMap[courseID] = make(map[int]*UserCourse)
				}
				courseUserMap[courseID][userID] = userCourse
				
			}
		}
	}
	CourseUserMap = courseUserMap
	return userCourses, nil
}

//LoadCourse 读取课程
func LoadCourse() (map[int]*Course, error) {
	req := db.CreateRequestDB("Select id,title from course where deleted_at is NULL", nil)
	res := db.CreateResponseDB()
	err := db.QuerySelect(req, res)

	if err != nil {
		return nil, err
	}
	courses := make(map[int]*Course)
	if res.Rows != nil {
		for _, v := range res.Rows {
			if v != nil {
				id, _ := strconv.Atoi(v["id"])
				title := v["title"]

				course := &Course{ID: id, Title: title}
				courses[id] = course
			}
		}
	}

	return courses, nil
}

//LoadMission 读取课程
func LoadMission() (map[int]*Mission, error) {
	req := db.CreateRequestDB("Select id,title from mission", nil)
	res := db.CreateResponseDB()
	err := db.QuerySelect(req, res)

	if err != nil {
		return nil, err
	}
	missions := make(map[int]*Mission)
	if res.Rows != nil {
		for _, v := range res.Rows {
			if v != nil {
				id, _ := strconv.Atoi(v["id"])
				title := v["title"]

				mission := &Mission{ID: id, Title: title}
				missions[id] = mission
			}
		}
	}
	return missions, nil
}

//LoadWork 读取作业
func LoadWork() (map[int]map[int]*Work, error) {
	req := db.CreateRequestDB("Select id,course_id, chapter_id, mission_id, user_id, images, created_at, updated_at from work", nil)
	res := db.CreateResponseDB()
	err := db.QuerySelect(req, res)

	if err != nil {
		return nil, err
	}
	works := make(map[int]map[int]*Work)
	workCourses := make(map[int]map[int]map[int]*Work)
	workAllMap := make(map[int]*Work)
	if res.Rows != nil {
		for _, v := range res.Rows {
			if v != nil {
				id, _ := strconv.Atoi(v["id"])
				courseID, _ := strconv.Atoi(v["course_id"])
				chapterID, _ := strconv.Atoi(v["chapter_id"])
				missionID, _ := strconv.Atoi(v["mission_id"])
				userID, _ := strconv.Atoi(v["user_id"])
				images := v["images"]
				createTime := v["created_at"]
				updateTime := v["updated_at"]

				cTime, err := utils.ParseStringTime(createTime)
				var cTimeSecond int64
				if err == nil {
					cTimeSecond = cTime.Unix()
				}
				uTime, err := utils.ParseStringTime(updateTime)
				var uTimeSecond int64
				if err == nil {
					uTimeSecond = uTime.Unix()
				}

				work := &Work{ID: id, CourseID: courseID, ChapterID: chapterID, MissionID: missionID, UserID: userID, CreateTime: cTimeSecond, UpdateTime: uTimeSecond}

				if images == "passDirectly" {
					work.Images = images
				}
				if works[userID] == nil {
					works[userID] = make(map[int]*Work)
				}

				if workCourses[userID] == nil {
					workCourses[userID] = make(map[int]map[int]*Work)
				}

				if workCourses[userID][work.CourseID] == nil {
					workCourses[userID][work.CourseID] = make(map[int]*Work)
				}

				works[userID][work.ID] = work
				workCourses[userID][work.CourseID][work.ID] = work
				workAllMap[work.ID] = work
			}
		}
	}
	WorkCourseMap = workCourses
	WorkAllMap = workAllMap
	return works, nil
}

//LoadUserExperiences 读取用户体验记录
func LoadUserExperiences() (map[int]*UserExperience, error) {
	req := db.CreateRequestDB("Select user_id,course_id,started_at from user_experience", nil)
	res := db.CreateResponseDB()
	err := db.QuerySelect(req, res)

	if err != nil {
		return nil, err
	}
	nowTime := time.Now().Unix()
	userExperiences := make(map[int]*UserExperience)
	UserExperienceAllMap = make(map[int]*UserExperience)
	if res.Rows != nil {
		for _, v := range res.Rows {
			if v != nil {
				userID, _ := strconv.Atoi(v["user_id"])
				courseID, _ := strconv.Atoi(v["course_id"])
				started := v["started_at"]
				cTime, _ := utils.ParseStringTime(started)
				userExperience := &UserExperience{UserID: userID, CourseID: courseID, Started: cTime.Unix()}
				UserExperienceAllMap[userID] = userExperience
				if nowTime-cTime.Unix() > 15*24*3600 {
					continue
				}

				userExperiences[userID] = userExperience
			}
		}
	}
	return userExperiences, nil
}

//LoadOrder 读取订单
func LoadOrder() (map[int]*Order, error) {
	req := db.CreateRequestDB("Select order_summary_id,period,price,quantity,created_at from `order` ", nil)
	res := db.CreateResponseDB()
	err := db.QuerySelect(req, res)

	if err != nil {
		return nil, err
	}
	orders := make(map[int]*Order)
	orderTotalMap := make(map[int]*Order)
	if res.Rows != nil {
		for _, v := range res.Rows {
			if v != nil {
				summaryID, _ := strconv.Atoi(v["order_summary_id"])
				period, _ := strconv.Atoi(v["period"])
				price, _ := strconv.Atoi(v["price"])
				quantity, _ := strconv.Atoi(v["quantity"])
				created := v["created_at"]

				cTime, _ := utils.ParseStringTime(created)

				nowTime := time.Now().Unix()

				order := &Order{SummaryID: summaryID, Period: period, Price: price, Quantity: quantity, Created: cTime.Unix()}

				orderTotalMap[summaryID] = order

				if period == 12 {
					//年卡
					if nowTime-cTime.Unix() > int64(3600*24*365*quantity) {
						continue
					}
				} else if period == 6 {
					//半年卡
					if nowTime-cTime.Unix() > int64(3600*24*365*quantity/2) {
						continue
					}
				} else if period == 3 {
					//季卡
					if nowTime-cTime.Unix() > int64(3600*24*365*quantity/4) {
						continue
					}
				}

				orders[summaryID] = order
			}
		}
	}

	OrderTotalMap = orderTotalMap
	return orders, nil
}

//LoadOrderSummary 读取订单
func LoadOrderSummary() (map[int]map[int]*OrderSummary, error) {
	req := db.CreateRequestDB("Select id,user_id,created_at from orderSummary where status = 2 and deleted_at is NULL ", nil)
	res := db.CreateResponseDB()
	err := db.QuerySelect(req, res)

	if err != nil {
		return nil, err
	}
	orderSummarys := make(map[int]map[int]*OrderSummary)
	if res.Rows != nil {
		for _, v := range res.Rows {
			if v != nil {
				id, _ := strconv.Atoi(v["id"])
				userID, _ := strconv.Atoi(v["user_id"])
				create_time := v["created_at"]
				cTime, _ := utils.ParseStringTime(create_time)
				if orderSummarys[userID] == nil {
					orderSummarys[userID] = make(map[int]*OrderSummary)
				}

				orderSummary := &OrderSummary{ID: id, UserID: userID, Created_at: cTime.Unix()}
				orderSummarys[userID][id] = orderSummary
			}
		}
	}
	return orderSummarys, nil
}

//LoadRating 读取评价
func LoadRating() (map[int]*Rating, error) {
	req := db.CreateRequestDB("Select * from rating ", nil)
	res := db.CreateResponseDB()
	err := db.QuerySelect(req, res)

	if err != nil {
		return nil, err
	}
	ratings := make(map[int]*Rating)
	ratingUserMap := make(map[int]map[int]map[int]*Rating)
	ratingCourseMap := make(map[int]map[int]*Rating)
	teacherMaps := make(map[int]*User)
	if res.Rows != nil {
		for _, v := range res.Rows {
			if v != nil {
				id, _ := strconv.Atoi(v["id"])
				workID, _ := strconv.Atoi(v["work_id"])
				rank, _ := strconv.Atoi(v["rank"])
				missionID, _ := strconv.Atoi(v["mission_id"])
				courseID, _ := strconv.Atoi(v["course_id"])
				userID, _ := strconv.Atoi(v["user_id"])
				ratingDuration, _ := strconv.Atoi(v["rating_duration"])
				updateTime := v["updated_at"]

				uTime, _ := utils.ParseStringTime(updateTime)

				rating := &Rating{ID: id, WorkID: workID, Rank: rank, MissionID: missionID, CourseID: courseID, UserID: userID, RatingDuration: ratingDuration, UpdateTime: uTime.Unix()}

				if ratingUserMap[userID] == nil {
					ratingUserMap[userID] = make(map[int]map[int]*Rating)
				}

				if ratingUserMap[userID][courseID] == nil {
					ratingUserMap[userID][courseID] = make(map[int]*Rating)
				}

				ratingUserMap[userID][courseID][workID] = rating

				if ratingCourseMap[courseID] == nil {
					ratingCourseMap[courseID] = make(map[int]*Rating)
				}

				ratingCourseMap[courseID][workID] = rating

				ratings[workID] = rating
				if ratingDuration > 0 && ratingDuration <= 21000000000 {

					if UserAllMap[userID] != nil {

						if UserAllMap[userID].Nickname == "CDwalker" {
							continue
						}

						if UserAllMap[userID].Nickname == "几大王" || UserAllMap[userID].Nickname == "随喜小王子" {
							continue
						}

						if UserAllMap[userID].Nickname == "cd404" {
							continue
						}

						teacherMaps[userID] = UserAllMap[userID]

					}
				}

			}
		}
	}

	RatingUserMap = ratingUserMap
	RatingCourseMap = ratingCourseMap
	TeacherMaps = teacherMaps
	return ratings, nil
}

//LoadLeaves 读取离开表
func LoadLeaves() (map[int][]*Leave, error) {
	req := db.CreateRequestDB("Select user_id,start_date,period from `leave` ", nil)
	res := db.CreateResponseDB()
	err := db.QuerySelect(req, res)

	if err != nil {
		return nil, err
	}
	leaves := make(map[int][]*Leave)
	if res.Rows != nil {
		for _, v := range res.Rows {
			if v != nil {
				userID, _ := strconv.Atoi(v["user_id"])
				startDate := v["start_date"]
				period, _ := strconv.Atoi(v["period"])

				sTime, _ := utils.ParseStringDayTypeTime(startDate)

				leave := &Leave{UserID: userID, StartDate: sTime.Unix(), EndDate: sTime.Unix() + int64(period*3600*24)}
				leaves[userID] = append(leaves[userID], leave)
			}
		}
	}
	return leaves, nil
}

//LoadHuanBi 读取环比
func LoadHuanBi() (map[int]*HuanBi, error) {
	req := db.CreateRequestDBCRM("Select userId from `huanbi` ", nil)
	res := db.CreateResponseDB()
	err := db.QuerySelect(req, res)

	if err != nil {
		return nil, err
	}
	huanbis := make(map[int]*HuanBi)
	if res.Rows != nil {
		for _, v := range res.Rows {
			if v != nil {
				userID, _ := strconv.Atoi(v["userId"])


				huanbi := &HuanBi{UserID: userID}
				huanbis[userID] = huanbi
			}
		}
	}
	return huanbis, nil
}

//GetNewestCourse 获取学员最新通道
func (user *User) GetNewestCourse() *Course {
	courses := UserCourseMap[user.ID]

	if courses == nil {
		return nil
	}
	courseID := 0
	for _, v := range courses {
		if v != nil {
			if courseID == 0 {
				courseID = v.CourseID
			} else if courseID < v.CourseID {
				courseID = v.CourseID
			}
		}
	}

	return CourseMap[courseID]
}

//GetAllCourses 获取所有的通道名字
func (user *User) GetAllCourses() string {
	courses := UserCourseMap[user.ID]
	courseName := "无"
	if courses == nil {
		return courseName
	}

	for _, v := range courses {
		if v != nil {
			course := CourseMap[v.CourseID]
			if course == nil {
				continue
			}
			if courseName == "无" {
				courseName = course.Title
			} else {
				courseName = fmt.Sprintf("%s,%s", courseName, course.Title)
			}
		}
	}

	return courseName
}

//GetUserExpCourses 获取学院体验期通道的名字
func GetUserExpCourses(userID int) string {
	if ExperCourses == nil {
		ExperCourses = make(map[int]int)
		ExperCourses[22] = 0
		ExperCourses[23] = 0
		ExperCourses[24] = 0
		ExperCourses[90] = 0
	}
	var expCourseString = ""
	for _, v := range UserCourseMap[userID] {
		if v != nil && ExperCourses[v.CourseID] == 0 {
			c := CourseMap[v.CourseID]
			if c != nil {
				if expCourseString == "" {
					expCourseString = c.Title
				} else {
					expCourseString = fmt.Sprintf("%s;%s", expCourseString, c.Title)
				}
			}

		}
	}

	return expCourseString
}

//判断当前用户是否正在请假
func (user *User) CheckUserInLeave() string {
	nowTime := time.Now().Unix()
	for _, l := range LeaveMap[user.ID] {
		if l != nil {
			if l.StartDate <= nowTime && l.EndDate >= nowTime {
				return "是"
			}
		}
	}
	return "否"
}

//CheckUserIsPayUser 检查用户是否是付费
func (user *User) CheckUserIsPayUser() bool {
	if OrderSummaryMap[user.ID] == nil {
		return false
	}

	if len(OrderSummaryMap[user.ID]) <= 0 {
		return false
	}
	//fmt.Println("user.ID = ", user.ID)
	nowTime := time.Now().Unix()

	for _, v := range OrderSummaryMap[user.ID] {
		if v != nil {
			if OrderTotalMap[v.ID] != nil {

				if OrderTotalMap[v.ID].Period == 12 {
					//年卡
					var addTime int64
					for _, l := range LeaveMap[user.ID] {
						if l != nil {
							if l.StartDate >= OrderTotalMap[v.ID].Created && l.StartDate <= OrderTotalMap[v.ID].Created+int64(3600*24*365*OrderTotalMap[v.ID].Quantity) {
								addTime += l.EndDate - l.StartDate
							}
						}
					}
					if nowTime-OrderTotalMap[v.ID].Created > int64(3600*24*365*OrderTotalMap[v.ID].Quantity)+addTime {

						continue
					}
				} else if OrderTotalMap[v.ID].Period == 6 {
					//半年卡
					var addTime int64
					for _, l := range LeaveMap[user.ID] {
						if l != nil {
							if l.StartDate >= OrderTotalMap[v.ID].Created && l.StartDate <= OrderTotalMap[v.ID].Created+int64(3600*24*365*OrderTotalMap[v.ID].Quantity/2) {
								addTime += l.EndDate - l.StartDate
							}
						}
					}
					if nowTime-OrderTotalMap[v.ID].Created > int64(3600*24*365*OrderTotalMap[v.ID].Quantity/2)+addTime {

						continue
					}
				} else if OrderTotalMap[v.ID].Period == 3 {
					//季卡
					var addTime int64
					for _, l := range LeaveMap[user.ID] {
						if l != nil {
							if l.StartDate >= OrderTotalMap[v.ID].Created && l.StartDate <= OrderTotalMap[v.ID].Created+int64(3600*24*365*OrderTotalMap[v.ID].Quantity/4) {
								addTime += l.EndDate - l.StartDate
							}
						}
					}
					if nowTime-OrderTotalMap[v.ID].Created > int64(3600*24*365*OrderTotalMap[v.ID].Quantity/4)+addTime {

						continue
					}
				}

				return true
			}
		}
	}

	return false
}

//CheckUserLastOrderInTime 判断用户最后一个订单在时间内
func (user *User)CheckUserLastOrderInTime(startTime, endTime int64) (int64, bool) {
	oss := OrderSummaryMap[user.ID]
	var newestOS *OrderSummary
	if oss != nil {
		for _, v := range oss {
			if v != nil {
				if newestOS == nil {
					newestOS = v
				} else if newestOS.Created_at < v.Created_at {
					newestOS = v
				}
			}
		}
	}
	if newestOS == nil {
		return  0, false
	}

	if newestOS.Created_at > endTime || newestOS.Created_at < startTime {
		return 0, false
	}
	return newestOS.Created_at, true
}

//GetEndTime 获取结束时间
func (order *Order) GetEndTime() int64 {
	if order.Period == 6 {
		return order.Created + (365 / 2 * 3600 * 24)
	} else if order.Period == 12 {
		return order.Created + (365 * 3600 * 24)
	} else if order.Period == 3 {
		return order.Created + (365 / 4 * 3600 * 24)
	}
	return order.Created
}

//GetPeriodTime 获取付费时长 月
func (order *Order) GetPeriodTime() int {
	switch order.Period {
	case 3:
		return 4
	case 6:
		return 6
	case 12:
		return 12
	}
	return 0
}

//CheckUserIsTestUser 检查用户是否是体验用户
func (user *User) CheckUserIsTestUser() bool {
	if UserExperienceMap[user.ID] == nil {
		return false
	}

	return true
}

//GetUserOrderType 获取用户类型 (体验期用户，无缴费用户，新人用户，2年以下复购用户，2年以上复购用户)
func GetUserOrderType(userID int) string {
	inExper := false
	if UserExperienceMap[userID] != nil {
		inExper = UserExperienceMap[userID].CheckUserInExperience()
	}

	if OrderSummaryMap[userID] == nil {
		if inExper {
			return "体验期用户"
		}
		return "无缴费用户"
	}
	payNum := 0
	payTotalTime := 0
	for _, v := range OrderSummaryMap[userID] {
		if v != nil {
			// if userID == 520 {
			// 	fmt.Println("v.ID = ", v.ID)
			// }
			payNum++
			o := OrderTotalMap[v.ID]
			if o != nil {


				// if userID == 520 {
				// 	fmt.Println("o.Period = ", o.Period)
				// }
				payTotalTime += o.GetPeriodTime()
			}
		}
	}
	// if userID == 520 {
	// 	fmt.Println("payTotalTime = " , payTotalTime)
	// }

	
	if payNum == 0 {
		return "无缴费用户"
	}

	if payNum == 1 {
		return "新人用户"
	}

	if payTotalTime < 24 {
		return "2年以下复购用户"
	}
	return "2年以上复购用户"
}

//CheckUserInExperience 判断用户在体验期内
func (userExperience *UserExperience) CheckUserInExperience() bool {
	if time.Now().Unix() - userExperience.Started > 15 * 86400 {
		return false
	}

	if time.Now().Unix()-userExperience.Started > 14*86400 {
		if time.Now().Unix()%86400 < userExperience.Started%86400 {
			return false
		}
	}

	return true
}

//GetWorkRecordString 获取作业记录
func (work *Work) GetWorkRecordString() string {
	result := utils.GetTimeString(utils.ParseUnixTime(work.CreateTime))

	c := CourseMap[work.CourseID]

	if c != nil {
		result = fmt.Sprintf("%s-%s", result, c.Title)
	}

	m := MissionMap[work.MissionID]

	if m != nil {
		result = fmt.Sprintf("%s-%s", result, m.Title)
	}

	return result
}

//GetUserFromMap 获取用户
func GetUserFromMap(userID int) *User {
	for _, v := range UserMap {
		if v != nil {
			if v[userID] != nil {
				return v[userID]
			}
		}
	}
	return nil
}

//GetUserFromMapByName 通过用户名获取用户
func GetUserFromMapByName(userName string) map[int]*User {
	returnMap := make(map[int]*User)
	for _, v := range StudentMap {
		if v != nil && strings.Contains(v.Nickname, userName) {
			returnMap[v.ID] = v
		}
	}
	return returnMap
}
