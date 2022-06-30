package controllers

import (
	"dgServer/conf"
	"dgServer/models"
	"dgServer/utils"
	"fmt"
	"strings"
	"time"

	"github.com/astaxie/beego/context"

	"github.com/astaxie/beego"
	"strconv"
)

const (
	MENU_SUB_TYPE_SYSDATA = iota + 1
	MENU_SUB_TYPE_ONLINE
	MENU_SUB_TYPE_PLAYERINFO
	MENU_SUB_TYPE_HERODATAS
	MENU_SUB_TYPE_ITEMDATAS
	MENU_SUB_TYPE_DRAGONDATAS
	MENU_SUB_TYPE_FAIRYDATAS
	MENU_SUB_TYPE_RES
)

const (
	USER_TYPE_NOR  = 10
	USER_TYPE_TEST = 20 //体验学员
	USER_TYPE_PAY  = 30 //付费学员
)

type LoginController struct {
	beego.Controller
}

type MainController struct {
	beego.Controller
}

type UsersStatus struct {
	beego.Controller
}

type TaskStatus struct {
	beego.Controller
}

type TaskScore struct {
	beego.Controller
}

type PayInfo struct {
	beego.Controller
}

type PayRate struct {
	beego.Controller
}

type TeacherInfo struct {
	beego.Controller
}

type UserDetail struct {
	beego.Controller
}

var sessionMap map[string]int64

const sessionTime = 1800 //半小时不操作 超时

var FilterUser = func(ctx *context.Context) {
	models.LoadDataAsync()
	_, ok := ctx.Input.Session("username").(string)
	ok2 := strings.Contains(ctx.Request.RequestURI, "/login")
	if !ok && !ok2 {
		fmt.Println("session fail relogin")
		ctx.Redirect(302, "/login")
	}
}

func (c *LoginController) Get() {

	c.TplName = "login.html"
}

func (c *LoginController) Post() {

	username := c.GetString("username")
	password := c.GetString("password")

	fmt.Println("login process")
	fmt.Println("password = ", password)
	if conf.CheckLogin(username, password) {
		fmt.Println("登录成功")
		c.SetSession("username", username)
		// c.Data["error"] = nil
		c.Ctx.WriteString("success")
		c.TplName = "login.html"

		// c.Ctx.Redirect(302, "/")

	} else {
		fmt.Println("登录失败")
		// c.Data["error"] = "username or password not correct"
		c.Ctx.WriteString("username or password not correct")
		c.TplName = "login.html"
		// c.Ctx.Redirect(302, "/login?error=username or password not correct")
	}
}

func (c *MainController) Get() {

	if models.RefreshTime <= 0 {
		models.CourseMap, _ = models.LoadCourse()
	}
	c.Data["DataType"] = 0
	c.Data["Index"] = 1

	c.Data["CourseDatas"] = models.CourseMap
	c.Data["ChooseID"] = 0
	c.Data["Choose01"] = 1
	c.Data["SelectIndexs"] = []int {0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16}
	c.Data["Offset"] = 0
	c.TplName = "index.html"
}

//Get 获取用户信息
func (c *UsersStatus) Get() {
	searchType, err := c.GetInt("SearchType")
	c.TplName = "index.html"
	c.Data["Index"] = 1
	c.Data["Choose01"] = 1
	c.Data["SelectIndexs"] = []int {0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16}
	if err != nil {
		c.Data["error"] = err
		return
	}
	c.Data["EndTime"] = utils.GetTimeDay(time.Now())

	// if err := models.LoadData(); err != nil {
	// 	c.Data["error"] = err
	// 	return
	// }

	c.Data["DataType"] = searchType

	c.Data["CourseDatas"] = models.CourseMap
	c.Data["ChooseID"] = 0

	isDownload, _ := c.GetInt("IsDownload")
	switch searchType {
	case 1:
		//付费学员 7天内未交
		diffTimeMin := int64(0)
		diffTimeMax := int64(3600 * 24 * 7)
		if isDownload == 1 {
			fileName := models.SetXLSFile("7日未交.xlsx", "", "", GetUserStatus(USER_TYPE_PAY, diffTimeMin, diffTimeMax))
			c.Ctx.Output.Download(fileName)
		} else {
			c.Data["MainDatas"] = GetUserStatus(USER_TYPE_PAY, diffTimeMin, diffTimeMax)
		}

	case 2:
		//付费学员 一个月以内未交
		diffTimeMin := int64(3600 * 24 * 7)
		diffTimeMax := int64(3600 * 24 * 30)

		if isDownload == 1 {
			fileName := models.SetXLSFile("30日内未交.xlsx", "", "", GetUserStatus(USER_TYPE_PAY, diffTimeMin, diffTimeMax))
			c.Ctx.Output.Download(fileName)
		} else {
			c.Data["MainDatas"] = GetUserStatus(USER_TYPE_PAY, diffTimeMin, diffTimeMax)
		}
	case 3:
		//付费学员 一个月以上未交
		diffTimeMin := int64(3600 * 24 * 30)

		if isDownload == 1 {
			fileName := models.SetXLSFile("30日以上未交.xlsx", "", "", GetUserStatus(USER_TYPE_PAY, diffTimeMin, 999999999))
			c.Ctx.Output.Download(fileName)
		} else {
			c.Data["MainDatas"] = GetUserStatus(USER_TYPE_PAY, diffTimeMin, 999999999)
		}
	case 11:
		//体验学员 2天未交
		diffTimeMin := int64(0)
		diffTimeMax := int64(3600 * 24 * 2)

		if isDownload == 1 {
			fileName := models.SetXLSFile("2天体验未交.xlsx", "", "", GetUserStatus(USER_TYPE_TEST, diffTimeMin, diffTimeMax))
			c.Ctx.Output.Download(fileName)
		} else {
			c.Data["SubDatas"] = GetUserStatus(USER_TYPE_TEST, diffTimeMin, diffTimeMax)
		}

	case 12:
		//体验学员 7天未交
		diffTimeMin := int64(3600 * 24 * 2)
		diffTimeMax := int64(3600 * 24 * 7)
		if isDownload == 1 {
			fileName := models.SetXLSFile("7天体验未交.xlsx", "", "", GetUserStatus(USER_TYPE_TEST, diffTimeMin, diffTimeMax))
			c.Ctx.Output.Download(fileName)
		} else {
			c.Data["SubDatas"] = GetUserStatus(USER_TYPE_TEST, diffTimeMin, diffTimeMax)
		}

	case 13:
		//体验学员 15天未交
		diffTimeMin := int64(3600 * 24 * 7)
		diffTimeMax := int64(3600 * 24 * 15)

		if isDownload == 1 {
			fileName := models.SetXLSFile("15天体验未交.xlsx", "", "", GetUserStatus(USER_TYPE_TEST, diffTimeMin, diffTimeMax))
			c.Ctx.Output.Download(fileName)
		} else {
			c.Data["SubDatas"] = GetUserStatus(USER_TYPE_TEST, diffTimeMin, diffTimeMax)
		}
	case 7:
		//学员缴费时长
		SearchTime4 := c.GetString("SearchTime4")
		SearchTime5 := c.GetString("SearchTime5")
		NumOfDay, _ := c.GetInt("NumOfDay")

		c.Data["SearchTime4"] = SearchTime4
		c.Data["SearchTime5"] = SearchTime5
		c.Data["ChooseID"] = NumOfDay

		times1 := strings.Split(SearchTime4, "-")
		times2 := strings.Split(SearchTime5, "-")

		if len(times1) >= 2 && len(times2) >= 2 {
			timeStart1, _ := utils.ParseStringDayTime(times1[0])
			timeEnd1, _ := utils.ParseStringDayTime(times1[1])

			timeStart2, _ := utils.ParseStringDayTime(times2[0])
			timeEnd2, _ := utils.ParseStringDayTime(times2[1])

			lrRecords := GetLazeUserRecord(timeStart1.Unix(), timeEnd1.Unix(), timeStart2.Unix(), timeEnd2.Unix(), NumOfDay)

			if isDownload == 1 {
				fileName := models.SetXLSFile(fmt.Sprintf("缴费学员连续%d天未交状态", NumOfDay), times1[0], times1[1], lrRecords)
				c.Ctx.Output.Download(fileName)
			} else {
				c.Data["LazeUserRecords"] = lrRecords
			}
		}	
	case 21:
		SearchTime3 := c.GetString("SearchTime3")
		c.Data["SearchTime3"] = SearchTime3

		times := strings.Split(SearchTime3, "-")

		if len(times) >= 2 {
			timeStart, _ := utils.ParseStringDayTime(times[0])
			timeEnd, _ := utils.ParseStringDayTime(times[1])
			c.Data["WorkDatas"] = GetWorksRecord(timeStart.Unix(), timeEnd.Unix()+(3600*24))
		}
	case 101:
		//算出连续提交作业的人
		endTime := c.GetString("EndTime")
		//fmt.Println(endTime)

		eTime, _ := utils.ParseStringDayTime(endTime)

		//fmt.Println(eTime.Unix())
		c.Data["EndTime"] = endTime
		c.Data["ContinueDatas"] = GetUserStatusWithContinueSubmit(7, eTime.Unix())
	case 102:
		//算出连续提交作业的人
		endTime := c.GetString("EndTime")

		eTime, _ := utils.ParseStringDayTime(endTime)
		c.Data["EndTime"] = endTime
		c.Data["ContinueDatas"] = GetUserStatusWithContinueSubmit(30, eTime.Unix())
	case 201:
		//算出连续提交作业的人

		courseID, _ := c.GetInt("ChooseID")
		SearchTime1 := c.GetString("SearchTime1")

		c.Data["ChooseID"] = courseID
		c.Data["SearchTime1"] = SearchTime1
		times := strings.Split(SearchTime1, "-")

		if len(times) >= 2 {
			timeStart, _ := utils.ParseStringDayTime(times[0])
			timeEnd, _ := utils.ParseStringDayTime(times[1])
			c.Data["MaxDatas"] = GetUserMaxWorksUser(courseID, timeStart.Unix(), timeEnd.Unix()+(3600*24))
		}
	case 301:
		//算出最快通过的用户
		SearchTime2 := c.GetString("SearchTime2")
		courseID, _ := c.GetInt("Choose01")
		times := strings.Split(SearchTime2, "-")
		c.Data["SearchTime2"] = SearchTime2
		c.Data["Choose01"] = courseID
		if len(times) >= 2 {
			timeStart, _ := utils.ParseStringDayTime(times[0])
			timeEnd, _ := utils.ParseStringDayTime(times[1])
			c.Data["FastDatas"] = MaxStarUser(courseID, timeStart.Unix(), timeEnd.Unix()+(3600*24))
		}

	}
}

func (c *TaskStatus) Get() {
	searchType, err := c.GetInt("SearchType")
	c.TplName = "taskInfo.html"
	c.Data["Index"] = 2
	c.Data["ChooseID"] = 0
	//c.Data["EndTime"] = utils.GetTimeDay(time.Now())
	if err != nil {
		fmt.Println("err = ", err)
		c.Data["error"] = err
		return
	}

	// if err := models.LoadData(); err != nil {
	// 	fmt.Println("err = ", err)
	// 	c.Data["error"] = err
	// 	return
	// }

	c.Data["DataType"] = searchType

	c.Data["CourseDatas"] = models.CourseMap

	switch searchType {
	case 101:
		//算出连续提交作业的人

		courseID, _ := c.GetInt("ChooseID")
		SearchTime1 := c.GetString("SearchTime1")

		c.Data["ChooseID"] = courseID
		c.Data["SearchTime1"] = SearchTime1
		times := strings.Split(SearchTime1, "-")

		if len(times) >= 2 {
			timeStart, _ := utils.ParseStringDayTime(times[0])
			timeEnd, _ := utils.ParseStringDayTime(times[1])
			c.Data["TaskDatas"] = GetTaskStatus(courseID, timeStart.Unix(), timeEnd.Unix()+(3600*24))
		}
	}
}

func (c *TaskScore) Get() {
	searchType, err := c.GetInt("SearchType")
	c.TplName = "taskInfo.html"
	c.Data["Index"] = 2
	c.Data["ChooseID"] = 0
	//c.Data["EndTime"] = utils.GetTimeDay(time.Now())
	if err != nil {
		fmt.Println("err = ", err)
		c.Data["error"] = err
		return
	}

	// if err := models.LoadData(); err != nil {
	// 	fmt.Println("err = ", err)
	// 	c.Data["error"] = err
	// 	return
	// }

	c.Data["DataType"] = searchType

	// c.Data["CourseDatas"] = models.CourseMap

	switch searchType {
	case 301:
		//算出作业评分

		// courseID, _ := c.GetInt("ChooseID")
		SearchTime2 := c.GetString("SearchTime2")
		fmt.Println("SearchTime2 = ", SearchTime2)
		// c.Data["ChooseID"] = courseID
		c.Data["SearchTime2"] = SearchTime2
		times := strings.Split(SearchTime2, "-")

		if len(times) >= 2 {
			timeStart, _ := utils.ParseStringDayTime(times[0])
			timeEnd, _ := utils.ParseStringDayTime(times[1])
			scores := GetTaskScore(timeStart.Unix(), timeEnd.Unix()+(3600*24))

			fmt.Println("len(scores) = ", len(scores))
			c.Data["ScoresDatas"] = scores
		}
	case 302:
		// 导出数据 算出作业评分
		SearchTime2 := c.GetString("SearchTime2")
		times := strings.Split(SearchTime2, "-")
		if len(times) >= 2 {
			timeStart, _ := utils.ParseStringDayTime(times[0])
			timeEnd, _ := utils.ParseStringDayTime(times[1])
			scores := GetTaskScore(timeStart.Unix(), timeEnd.Unix()+(3600*24))
			fileName := models.SetXLSFile("作业评分", times[0], times[1], scores)
			c.Ctx.Output.Download(fileName)
		}
	}
}

func (c *PayInfo) Get() {
	searchType, err := c.GetInt("SearchType")
	c.TplName = "payInfo.html"
	c.Data["Index"] = 3

	if err != nil {
		fmt.Println("err = ", err)
		c.Data["error"] = err
		return
	}

	// if err := models.LoadData(); err != nil {
	// 	fmt.Println("err = ", err)
	// 	c.Data["error"] = err
	// 	return
	// }

	c.Data["DataType"] = searchType
	c.Data["ChooseID"] = 0
	c.Data["SelectIndexs"] = []int {0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16}

	isDownload, _ := c.GetInt("IsDownload")

	switch searchType {
	case 1:
		//下载体验期用户没有缴费的
		SearchTime5 := c.GetString("SearchTime5")
		c.Data["SearchTime5"] = SearchTime5
		times := strings.Split(SearchTime5, "-")

		if len(times) >= 2 {
			timeStart, _ := utils.ParseStringDayTime(times[0])
			timeEnd, _ := utils.ParseStringDayTime(times[1])

			fileName := models.SetXLSFile("体验用户未缴费", times[0], times[1], GetExperUserWithOutPay(timeStart.Unix(), timeEnd.Unix()+(3600*24)))
			c.Ctx.Output.Download(fileName)
		}

	case 2:
		//获取用户交作业次数
		SearchTime6 := c.GetString("SearchTime6")
		c.Data["SearchTime6"] = SearchTime6
		times := strings.Split(SearchTime6, "-")
		courses := c.GetString("inputCourses")

		if len(times) >= 2 {
			timeStart, _ := utils.ParseStringDayTime(times[0])
			timeEnd, _ := utils.ParseStringDayTime(times[1])

			fileName := models.SetXLSFile("体验用户未缴费", times[0], times[1], GetUserWorkTimes(timeStart.Unix(), timeEnd.Unix()+(3600*24), courses))
			c.Ctx.Output.Download(fileName)
		}

	case 101:
		//学员缴费到期统计
		SearchTime1 := c.GetString("SearchTime1")
		c.Data["SearchTime1"] = SearchTime1
		times := strings.Split(SearchTime1, "-")
		var payDatas []*models.PayInfo
		if len(times) >= 2 {
			timeStart, _ := utils.ParseStringDayTime(times[0])
			timeEnd, _ := utils.ParseStringDayTime(times[1])
			payDatas = GetPayInfo(timeStart.Unix(), timeEnd.Unix()+(3600*24))
		} else {
			payDatas = GetPayInfo(0, 99999999999)
		}

		if isDownload == 1 {
			fileName := models.SetXLSFile("学员缴费到期统计", times[0], times[1], payDatas)
			c.Ctx.Output.Download(fileName)
		} else {
			c.Data["PayDatas"] = payDatas
		}
	case 102:
		SearchTime1 := c.GetString("SearchTime1")
		c.Data["SearchTime1"] = SearchTime1
		times := strings.Split(SearchTime1, "-")
		var payDatas []*models.PayInfo
		if len(times) >= 2 {
			timeStart, _ := utils.ParseStringDayTime(times[0])
			// timeEnd, _ := utils.ParseStringDayTime(times[1])
			payDatas = GetPayInfoAdd(timeStart.Unix())
		} else {
			payDatas = GetPayInfoAdd(0)
		}

		fileName := models.SetXLSFile("学员缴费到期统计", times[0], times[1], payDatas)
		c.Ctx.Output.Download(fileName)
	case 301:
		//学员缴费记录
		SearchTime2 := c.GetString("SearchTime2")
		c.Data["SearchTime2"] = SearchTime2
		times := strings.Split(SearchTime2, "-")
		var orderDatas []*models.OrderRecord
		if len(times) >= 2 {
			timeStart, _ := utils.ParseStringDayTime(times[0])
			timeEnd, _ := utils.ParseStringDayTime(times[1])
			orderDatas = GetPayRecord(timeStart.Unix(), timeEnd.Unix()+(3600*24))
		}
		if isDownload == 1 {
			fileName := models.SetXLSFile("学员缴费记录", times[0], times[1], orderDatas)
			c.Ctx.Output.Download(fileName)
		} else {
			c.Data["PayRecords"] = orderDatas
		}
	case 401:
		//学员缴费时长
		SearchTime3 := c.GetString("SearchTime3")
		SearchTime4 := c.GetString("SearchTime4")
		NumOfDay, _ := c.GetInt("NumOfDay")

		c.Data["SearchTime3"] = SearchTime3
		c.Data["SearchTime4"] = SearchTime4
		c.Data["ChooseID"] = NumOfDay

		times1 := strings.Split(SearchTime3, "-")
		times2 := strings.Split(SearchTime4, "-")

		if len(times1) >= 2 && len(times2) >= 2 {
			timeStart1, _ := utils.ParseStringDayTime(times1[0])
			timeEnd1, _ := utils.ParseStringDayTime(times1[1])

			timeStart2, _ := utils.ParseStringDayTime(times2[0])
			timeEnd2, _ := utils.ParseStringDayTime(times2[1])

			lrRecords := GetLazeUserRecord(timeStart1.Unix(), timeEnd1.Unix(), timeStart2.Unix(), timeEnd2.Unix(), NumOfDay)

			if isDownload == 1 {
				fileName := models.SetXLSFile(fmt.Sprintf("缴费学员连续%d天未交状态", NumOfDay), times1[0], times1[1], lrRecords)
				c.Ctx.Output.Download(fileName)
			} else {
				c.Data["LazeUserRecords"] = lrRecords
			}
		}
	}

}

//Get 付费比率
func (c *PayRate) Get() {
	searchType, err := c.GetInt("SearchType")
	c.TplName = "payRate.html"
	c.Data["Index"] = 4

	if err != nil {
		fmt.Println("err = ", err)
		c.Data["error"] = err
		return
	}

	// if err := models.LoadData(); err != nil {
	// 	fmt.Println("err = ", err)
	// 	c.Data["error"] = err
	// 	return
	// }

	c.Data["DataType"] = searchType

	isDownload, _ := c.GetInt("IsDownload")

	switch searchType {
	case 101:
		//算出付费转化率
		SearchTime1 := c.GetString("SearchTime1")
		c.Data["SearchTime1"] = SearchTime1
		times := strings.Split(SearchTime1, "-")

		if len(times) >= 2 {
			timeStart, _ := utils.ParseStringDayTime(times[0])
			timeEnd, _ := utils.ParseStringDayTime(times[1])
			payRateDatas, totalRate, selectRate := GetPayRate(timeStart.Unix(), timeEnd.Unix()+(3600*24))
			c.Data["PayRateDatas"] = payRateDatas
			c.Data["PayTotalRate"] = totalRate
			c.Data["PaySelectRate"] = selectRate
			c.Data["UserTagMap"] = models.UserTagMap

			if isDownload == 1 {
				fileName := models.SetXLSFile("付费转化率", times[0], times[1], payRateDatas, selectRate, totalRate, 1)
				c.Ctx.Output.Download(fileName)
			} else {
				c.Data["PayRateDatas"] = payRateDatas
				c.Data["PayTotalRate"] = totalRate
				c.Data["PaySelectRate"] = selectRate
			}
		}
	case 201:
		//算出复购率
		SearchTime2 := c.GetString("SearchTime2")
		c.Data["SearchTime2"] = SearchTime2
		times := strings.Split(SearchTime2, "-")

		if len(times) >= 2 {
			timeStart, _ := utils.ParseStringDayTime(times[0])
			timeEnd, _ := utils.ParseStringDayTime(times[1])
			payDatas, repNum, totalNum := GetRePayInfo(timeStart.Unix(), timeEnd.Unix()+(3600*24))

			var payRate float32

			if totalNum > 0 {
				payRate = (float32)(repNum) * 100.0 / (float32)(totalNum)
			}
			c.Data["RePayDatas"] = payDatas
			c.Data["RePayTotalNum"] = repNum
			c.Data["RePayRate"] = payRate

			if isDownload == 1 {
				fileName := models.SetXLSFile("复购率", times[0], times[1], payDatas, repNum, payRate, 2)
				c.Ctx.Output.Download(fileName)
			} else {
				c.Data["RePayDatas"] = payDatas
				c.Data["RePayTotalNum"] = repNum
				c.Data["RePayRate"] = payRate
			}
		}
	}
}

//Get 点评师信息
func (c *TeacherInfo) Get() {
	searchType, err := c.GetInt("SearchType")
	c.TplName = "teacherInfo.html"
	c.Data["Index"] = 5
	c.Data["CourseID"] = 0
	c.Data["UserID"] = 0

	if err != nil {
		fmt.Println("err = ", err)
		c.Data["error"] = err
		return
	}

	// if err := models.LoadData(); err != nil {
	// 	fmt.Println("err = ", err)
	// 	c.Data["error"] = err
	// 	return
	// }

	c.Data["CourseDatas"] = models.CourseMap
	c.Data["UserDatas"] = models.TeacherMaps

	c.Data["DataType"] = searchType

	isDownload, _ := c.GetInt("IsDownload")

	switch searchType {
	case 101:
		//算出点评师信息
		SearchTime1 := c.GetString("SearchTime1")
		userID, _ := c.GetInt("UserID")
		courseID, _ := c.GetInt("CourseID")
		c.Data["SearchTime1"] = SearchTime1
		c.Data["UserID"] = userID
		c.Data["CourseID"] = courseID

		times := strings.Split(SearchTime1, "-")

		if len(times) >= 2 {
			timeStart, _ := utils.ParseStringDayTime(times[0])
			timeEnd, _ := utils.ParseStringDayTime(times[1])
			teacherStatus := GetTeacherStatus(timeStart.Unix(), timeEnd.Unix()+(3600*24), courseID, userID)
			if isDownload == 1 {
				fileName := models.SetXLSFile("点评师信息", times[0], times[1], teacherStatus)
				c.Ctx.Output.Download(fileName)
			} else {
				c.Data["TeacherStatus"] = teacherStatus
			}

		}
	}
}

//Get 用户侧写
func (c *UserDetail) Get() {
	searchType, err := c.GetInt("SearchType")
	c.TplName = "userDetail.html"
	c.Data["Index"] = 6
	c.Data["UserDatas"] = models.StudentMap
	if err != nil {
		fmt.Println("err = ", err)
		c.Data["error"] = err
		return
	}

	// if err := models.LoadData(); err != nil {
	// 	fmt.Println("err = ", err)
	// 	c.Data["error"] = err
	// 	return
	// }

	searchVal := c.GetString("SearchVal")
	c.Data["YearDatas"] = utils.GetStartYearList()
	if searchVal == "" {
		//所有数据
		c.Data["UserDatas"] = make(map[int]*models.User)
	} else {
		c.Data["UserDatas"] = models.GetUserFromMapByName(searchVal)
		c.Data["SearchVal"] = searchVal
	}

	userID, _ := c.GetInt("UserId")
	var user *models.User
	if userID > 0 {
		user = models.GetUserFromMap(userID)

		if user != nil {
			c.Data["UserId"] = userID
			c.Data["UserName"] = user.Nickname
		}
	}

	year, _ := c.GetInt("Year")
	month, _ := c.GetInt("Month")
	c.Data["Year"] = year
	c.Data["Month"] = month

	switch searchType {
	case 101:
		if user != nil {
			//存在用户 获取通道数据
			detailDatas, monthDatas, monthTotal := GetUserCourseDetail(user, year, month)
			c.Data["DetailDatas"] = detailDatas
			c.Data["MonthDatas"] = monthDatas
			c.Data["MonthTotal"] = monthTotal

			isDownload := c.GetString("IsDownload")
			if isDownload == "1" {
				fileName := models.SetXLSFile(fmt.Sprintf("%s.xlsx", user.Nickname), "", "", detailDatas, monthDatas, monthTotal, user.Nickname, year, month)
				c.Ctx.Output.Download(fileName)
			}
		}
	}
}

//GetUserStatus 获取用户信息
func GetUserStatus(userType int, diffTimeMin, diffTimeMax int64) []*models.UserStatus {
	users := models.UserMap[userType]

	nowTime := time.Now().Unix()

	var userStatus []*models.UserStatus
	if users != nil {
	OUTLOOP:
		for _, v := range users {
			if v != nil {
				if userType == USER_TYPE_PAY {
					if !v.CheckUserIsPayUser() {
						continue
					}
				} else if userType == USER_TYPE_TEST {

					if v.CheckUserIsPayUser() {
						continue
					}
					
					if !v.CheckUserIsTestUser() {
						continue
					}
				}

				works := models.WorkMap[v.ID]
				var diff int64
				var lastTime int64
				var work *models.Work
				if works != nil {

					for _, w := range works {
						if w != nil {

							if nowTime-w.CreateTime < diffTimeMin {
								continue OUTLOOP
							}

							if nowTime-w.CreateTime > diffTimeMax {
								continue
							}
							if w.Images == "passDirectly" {
								continue
							}
							if diff == 0 {
								diff = nowTime - w.CreateTime
								lastTime = w.CreateTime
								work = w
							} else if diff > nowTime-w.CreateTime {
								diff = nowTime - w.CreateTime
								lastTime = w.CreateTime
								work = w
							}
						}
					}
					cources := "无"
					if userType != USER_TYPE_TEST {
						if work == nil {
							continue
						}
						cources = v.GetAllCourses()
					} else {
						experice := models.UserExperienceMap[v.ID]

						if experice != nil {

							if work == nil {
								if nowTime-experice.Started < diffTimeMin {
									continue OUTLOOP
								}
								if nowTime-experice.Started > diffTimeMax {
									continue
								}
								diff = nowTime - experice.Started
							}

							// lastTime = experice.Started
							ex_course := models.CourseMap[experice.CourseID]
							if ex_course != nil {
								cources = ex_course.Title
							}
						}
					}

					lastName := "无"
					teacherName := "无"
					if work != nil {
						course := models.CourseMap[work.CourseID]

						if course != nil {
							lastName = course.Title
						}
						mission := models.MissionMap[work.MissionID]

						if mission != nil {
							lastName = fmt.Sprintf("%s,%s", lastName, mission.Title)
						}

						rating := models.RatingMap[work.ID]
						if rating != nil {
							teacher := models.TeacherMaps[rating.UserID]
							if teacher != nil {
								teacherName = teacher.Nickname
							}
						}
					}

					last_submit := "空"
					if lastTime > 0 {
						last_submit = utils.GetTimeString(utils.ParseUnixTime(lastTime))
					}


					last_day := "未提交"

					if diff > 0 {
						last_day = fmt.Sprintf("%d", int(diff / 3600 / 24))
					}

					us := &models.UserStatus{
						M_id:          v.ID,
						M_name:        v.Nickname,
						M_day:         last_day,
						M_course:      cources,
						M_last_submit: last_submit,
						M_last_name:   lastName,
						M_teacher:     teacherName,
					}

					for _, l := range models.LeaveMap[v.ID] {
						if l != nil {

							// if v.ID == 503 {
							// 	fmt.Println("nowTime= ", nowTime)
							// 	fmt.Println("startTime = ", startTime)
							// 	fmt.Println("endTime = ", endTime)
							// 	fmt.Println("l.StartDate = ", l.StartDate)
							// }
							if nowTime >= l.StartDate && nowTime <= l.EndDate {
								us.M_leave = 1
								break
							}
						}
					}
					userStatus = append(userStatus, us)
				} else {

					last_day := "未提交"

					if diff > 0 {
						last_day = fmt.Sprintf("%d", int(diff / 3600 / 24))
					}

					us := &models.UserStatus{
						M_id:          v.ID,
						M_name:        v.Nickname,
						M_day:         last_day,
						M_course:      v.GetAllCourses(),
						M_last_submit: "空",
						M_last_name:   "",
						M_teacher:     "",
					}

					for _, l := range models.LeaveMap[v.ID] {
						if l != nil {

							// if v.ID == 503 {
							// 	fmt.Println("nowTime= ", nowTime)
							// 	fmt.Println("startTime = ", startTime)
							// 	fmt.Println("endTime = ", endTime)
							// 	fmt.Println("l.StartDate = ", l.StartDate)
							// }
							if nowTime >= l.StartDate && nowTime <= l.EndDate {
								us.M_leave = 1
								break
							}
						}
					}
					userStatus = append(userStatus, us)
				}
			}
		}
	}
	//fmt.Println("len(userStatus) = ", len(userStatus))
	return userStatus

}

//GetUserStatusWithContinueSubmit 获取连续提交作业的人
func GetUserStatusWithContinueSubmit(days int, nowTime int64) []*models.UserStatus {
	var userStatus []*models.UserStatus
	totalNum := 0
	for userID, ws := range models.WorkMap {
		dayMaps := make(map[int]int)
		for index := 1; index <= days; index++ {
			dayMaps[index] = index
		}
		totalNum += len(ws)
		sign := 0
		//fmt.Println("userID = ", userID)
		for _, w := range ws {

			if nowTime-w.CreateTime < 0 {
				continue
			}
			day := int(((nowTime - w.CreateTime) / 3600 / 24) + 1)
			// if day == 1 {
			// 	fmt.Println("day = ", day)
			// }

			if userID == 491 {
				if day == 1 {
					fmt.Println(utils.ParseUnixTime(w.CreateTime))
				}
			}
			if dayMaps[day] > 0 {

				delete(dayMaps, day)

				if len(dayMaps) <= 0 {
					sign = 1
					break
				}
			}
		}

		if sign == 1 {
			//满足条件
			//fmt.Println("满足")
			user := models.GetUserFromMap(userID)

			if user != nil {
				us := &models.UserStatus{
					M_id:     user.ID,
					M_name:   user.Nickname,
					M_course: user.GetAllCourses(),
				}
				userStatus = append(userStatus, us)
			}

		}
	}

	//fmt.Println("totalNum =", totalNum)

	return userStatus
}

//GetWorksRecord 获取作业记录
func GetWorksRecord(startTime, endTime int64) []*models.WorkData {
	var works []*models.WorkData

	for _, v := range models.UserExperienceAllMap {
		if v != nil && v.CheckUserInExperience() {
			w := &models.WorkData{
				M_id:     v.UserID,
				M_create: utils.GetTimeString(utils.ParseUnixTime(v.Started)),
			}

			user := models.UserAllMap[v.UserID]

			if user != nil {
				w.M_name = user.Nickname
			}

			tagRels := models.UserTagRelMap[v.UserID]

			if tagRels != nil {
				for _, rel := range tagRels {
					if rel != nil {
						tag := models.UserTagMap[rel.TagID]
						if tag != nil {
							if w.M_tag == "" {
								w.M_tag = tag.Name
							} else {
								w.M_tag = fmt.Sprintf("%s;%s", w.M_tag, tag.Name)
							}
						}
					}
				}
			}

			cs := models.UserCourseMap[v.UserID]
			
			if cs != nil {
				for _, c1 := range cs {
					if c1 != nil {
						c := models.CourseMap[c1.CourseID]
						if c != nil {
							if w.M_course == "" {
								w.M_course = c.Title
							} else {
								w.M_course = fmt.Sprintf("%s;%s", w.M_course, c.Title)
							}
						}
					}
				}
			}


			user_works := models.WorkMap[v.UserID]

			if user_works != nil {
				for _, w1 := range user_works {
					if w1.Images == "passDirectly" {
						continue
					}
					if w1.CreateTime < startTime || w1.CreateTime > endTime {
						continue
					}
				
					w.M_times = w.M_times + 1

					if w.M_work == "" {
						w.M_work = w1.GetWorkRecordString()
					} else {
						w.M_work = fmt.Sprintf("%s;%s", w.M_work, w1.GetWorkRecordString())
					}
				}
			}

			works = append(works, w)
		}
	}
	return works
}

//GetUserMaxWorksUser 获取提交作业最多的人
func GetUserMaxWorksUser(courseID int, startTime, endTime int64) []*models.UserStatus {
	var userStatus []*models.UserStatus

	for userID, ws := range models.WorkMap {
		myWorkNum := 0
		for _, w := range ws {

			if courseID != 0 {
				if courseID != w.CourseID {
					continue
				}
			}
			if w.CreateTime <= endTime && w.CreateTime >= startTime {
				myWorkNum++
			}
		}

		if myWorkNum > 0 {
			user := models.GetUserFromMap(userID)

			if user != nil {
				courseName := ""
				if courseID == 0 {
					courseName = user.GetAllCourses()
				} else {
					course := models.CourseMap[courseID]

					if course != nil {
						courseName = course.Title
					}
				}
				us := &models.UserStatus{
					M_id:     user.ID,
					M_name:   user.Nickname,
					M_course: courseName,
					M_num:    myWorkNum,
				}
				userStatus = append(userStatus, us)
			}
		}

	}

	return userStatus
}

//MaxStarUser 通过最快玩家
func MaxStarUser(in_courseID int, startTime, endTime int64) []*models.UserStatus {
	var userStatus []*models.UserStatus
	for userID, v := range models.WorkCourseMap {
		if v != nil {

			if in_courseID > 0 {
				ws := v[in_courseID]
				currentStar3 := 0
				currentName := ""
				if ws != nil {
					for _, w := range ws {
						if w != nil {
							if w.CreateTime < startTime || w.CreateTime > endTime {
								continue
							}

							rating := models.RatingMap[w.ID]

							if rating != nil && rating.Rank == 1 {
								if rating.UpdateTime < startTime || rating.UpdateTime > endTime {
									continue
								}
								currentStar3++
								mission := models.MissionMap[w.MissionID]
								if mission != nil {
									if currentName == "" {
										currentName = mission.Title
									} else {
										currentName = fmt.Sprintf("%s,%s", currentName, mission.Title)
									}
								}
							}
						}
					}

					if currentStar3 > 0 {
						user := models.GetUserFromMap(userID)

						if user != nil {
							courseName := ""
							course := models.CourseMap[in_courseID]

							if course != nil {
								courseName = course.Title
							}
							us := &models.UserStatus{
								M_id:      user.ID,
								M_name:    user.Nickname,
								M_course:  courseName,
								M_mission: currentName,
								M_num:     currentStar3,
							}
							userStatus = append(userStatus, us)
						}
					}
				}

				continue
			}

			cStar3 := 0
			cCourseID := 0
			cName := ""
			for courseID, ws := range v {
				if ws != nil {
					currentStar3 := 0
					currentName := ""
					for _, w := range ws {
						if w != nil {

							if w.CreateTime < startTime || w.CreateTime > endTime {
								continue
							}
							rating := models.RatingMap[w.ID]

							if rating != nil && rating.Rank == 1 {
								if rating.UpdateTime < startTime || rating.UpdateTime > endTime {
									continue
								}
								currentStar3++
								mission := models.MissionMap[w.MissionID]
								if mission != nil {
									if currentName == "" {
										currentName = mission.Title
									} else {
										currentName = fmt.Sprintf("%s,%s", currentName, mission.Title)
									}
								}

							}
						}
					}

					if currentStar3 > cStar3 {
						cStar3 = currentStar3
						cCourseID = courseID
						cName = currentName
					}
				}
			}

			if cCourseID > 0 {
				user := models.GetUserFromMap(userID)

				if user != nil {
					courseName := ""
					course := models.CourseMap[cCourseID]

					if course != nil {
						courseName = course.Title
					}
					us := &models.UserStatus{
						M_id:      user.ID,
						M_name:    user.Nickname,
						M_course:  courseName,
						M_mission: cName,
						M_num:     cStar3,
					}
					userStatus = append(userStatus, us)
				}
			}
		}
	}
	return userStatus
}

//GetTaskStatus 获取任务状态
func GetTaskStatus(courseID int, starTime, endTime int64) []*models.SendTaskStatus {

	tsMap := make(map[int]*models.TaskStatus)

	for userID, v := range models.WorkCourseMap {
		if v != nil {
			user := models.GetUserFromMap(userID)
			if courseID > 0 {
				ws := v[courseID]
				if ws != nil {
					for _, w := range ws {
						if w != nil {

							if w.UpdateTime < starTime || w.UpdateTime > endTime {
								continue
							}
							if tsMap[w.MissionID] == nil {

								mission := models.MissionMap[w.MissionID]
								t_name := ""
								if mission != nil {
									t_name = mission.Title
								}
								tsMap[w.MissionID] = &models.TaskStatus{
									T_id:          w.MissionID,
									T_name:        t_name,
									T_usernames:   make(map[string]int),
									T_pass_data:   make(map[int]int),
									T_course_data: make(map[int]int),
								}
							}

							tsMap[w.MissionID].T_course_num++
							tsMap[w.MissionID].T_course_data[userID]++
							if models.RatingMap[w.ID] != nil {
								if models.RatingMap[w.ID].Rank == 1 {
									tsMap[w.MissionID].T_pass_num++
									tsMap[w.MissionID].T_pass_data[userID]++
								}
							}

							if user != nil {
								tsMap[w.MissionID].T_usernames[user.Nickname] = 1
							}

						}
					}
				}

			} else {
				for _, ws := range v {
					if ws != nil {
						for _, w := range ws {
							if w != nil {

								if w.UpdateTime < starTime || w.UpdateTime > endTime {
									continue
								}
								if tsMap[w.MissionID] == nil {

									mission := models.MissionMap[w.MissionID]
									t_name := ""
									if mission != nil {
										t_name = mission.Title
									}
									tsMap[w.MissionID] = &models.TaskStatus{
										T_id:          w.MissionID,
										T_name:        t_name,
										T_usernames:   make(map[string]int),
										T_pass_data:   make(map[int]int),
										T_course_data: make(map[int]int),
									}
								}

								tsMap[w.MissionID].T_course_num++
								tsMap[w.MissionID].T_course_data[userID]++
								if models.RatingMap[w.ID] != nil {
									if models.RatingMap[w.ID].Rank == 1 {
										tsMap[w.MissionID].T_pass_num++
										tsMap[w.MissionID].T_pass_data[userID]++
									}
								}

								if user != nil {
									tsMap[w.MissionID].T_usernames[user.Nickname] = 1
								}

							}
						}
					}
				}
			}
		}
	}
	if len(tsMap) > 0 {
		returnTaskStatus := make([]*models.SendTaskStatus, len(tsMap))
		n := 0
		for _, v := range tsMap {
			if v != nil {
				returnTaskStatus[n] = v.GetSendData()
				n++
			}
		}
		return returnTaskStatus
	}
	return nil
}

//获取评分
func GetTaskScore(starTime, endTime int64) []*models.SendTaskScore {
	var returnTaskStatus []*models.SendTaskScore
	for _, v := range models.WorkCourseMap {

		for _, ws := range v {
			if ws != nil {
				for _, w := range ws {
					if w != nil {

						if w.CreateTime < starTime || w.CreateTime > endTime {
							continue
						}

						mission := models.MissionMap[w.MissionID]
						t_name := ""
						if mission != nil {
							t_name = mission.Title
						}
						course := models.CourseMap[w.CourseID]
						c_name := ""
						if course != nil {
							c_name = course.Title
						}

						user := models.GetUserFromMap(w.UserID)
						u_name := ""
						if user != nil {
							if user.Privilege == 50 {
								continue
							}
							u_name = user.Nickname
						}

						ts := &models.TaskScore{
							SD_utime:       utils.GetTimeString(utils.ParseUnixTime(w.CreateTime)),
							SD_id:          w.UserID,
							SD_name:        t_name,
							SD_username:    u_name,
							SD_course_name: c_name,
						}

						if models.RatingMap[w.ID] != nil {
							ts.SD_stime = utils.GetTimeString(utils.ParseUnixTime(models.RatingMap[w.ID].UpdateTime))
							ts.SD_rank = models.RatingMap[w.ID].Rank
						}

						returnTaskStatus = append(returnTaskStatus, ts.GetSendData())
					}
				}
			}
		}
	}
	return returnTaskStatus
}

func GetExperUserWithOutPay(startTime, endTime int64) []*models.SimpleUser {
	var susers []*models.SimpleUser
	for _, v := range models.UserExperienceAllMap {
		if v != nil {
			if v.Started < startTime {
				continue
			}
			if v.Started > endTime {
				continue
			}

			os := models.OrderSummaryMap[v.UserID]
			huanbi := models.HuanBiMap[v.UserID]
			if len(os) == 0 && huanbi != nil {
				user := models.UserAllMap[v.UserID]
				u := &models.SimpleUser{
					S_User_id: v.UserID,

				}
				if user != nil {
					u.S_name = user.Nickname
				}
				susers = append(susers, u)
			}
		}
	}

	return susers;
}

//GetUserWorkTimes 获取玩家作业次数
func GetUserWorkTimes(startTime, endTime int64, course string) map[int]map[int]map[int]*models.WorkTimes{
	wTimes := make(map[int]map[int]map[int]*models.WorkTimes)
	courses:= strings.Split(course, ",")

	for _, c := range courses {
		cID, err := strconv.Atoi(c)
		if err == nil {
			wTimes[cID] = make(map[int]map[int]*models.WorkTimes)	
		}
	}

	for _, w := range models.WorkAllMap {
		if w != nil && w.CreateTime > startTime && w.CreateTime < endTime {
			if userMap, ok := wTimes[w.CourseID];ok {
				if userMap[w.MissionID] == nil {
					userMap[w.MissionID] = make(map[int]*models.WorkTimes)	
				}
				if userMap[w.MissionID][w.UserID] != nil {
					userMap[w.MissionID][w.UserID].W_num = userMap[w.MissionID][w.UserID].W_num + 1
				} else{
					userMap[w.MissionID][w.UserID] = &models.WorkTimes{
						W_user_name: models.UserAllMap[w.UserID].Nickname,
						W_course: models.CourseMap[w.CourseID].Title,
						W_mission: models.MissionMap[w.MissionID].Title,
						W_num: 1,
					}
				}
			}
		}
	}
	return wTimes
}
 
//GetPayInfoAdd 获取要添加的人员详情
func GetPayInfoAdd(startTime int64) []*models.PayInfo {
	getPlayers := make(map[int]*models.PayInfo)

	for userID, v := range models.OrderSummaryMap {
		if v != nil {
			for sID := range v {

				order := models.OrderTotalMap[sID]

				if order != nil {
					if order.Created > startTime {
						continue
					} else {
						if getPlayers[userID] == nil {
							getPlayers[userID] = &models.PayInfo{
								P_id: userID,
							}

							user := models.UserAllMap[userID]

							if user != nil {
								getPlayers[userID].P_name = user.Nickname
								getPlayers[userID].P_phone = user.Mobile
								getPlayers[userID].P_mail = user.Email
								getPlayers[userID].P_aliPay = user.Alipay_nickname

								if user.Delete != "" {
									getPlayers[userID].P_delete = 1
								}
							}
						}
						pt := &models.PayTimes{StartTime: order.Created, EndTime: order.GetEndTime()}
						if getPlayers[userID].PayTimes == nil {
							getPlayers[userID].PayTimes = pt
							continue
						}
						node := getPlayers[userID].PayTimes
						var preNode *models.PayTimes
						for node != nil {
							if node.StartTime < pt.StartTime {
								if node.Next == nil {
									node.Next = pt
									break
								} else {
									preNode = node
									node = node.Next
								}
							} else {
								if preNode == nil {
									pt.Next = node
									getPlayers[userID].PayTimes = pt
									break
								} else {
									preNode.Next = pt
									pt.Next = node
									break
								}
							}
						}

						//for index := len(getPlayers[userID].PayTimes); index >= 1; index-- {
						//if getPlayers[userID].PayTimes[index].StartTime > pt.StartTime {
						//	getPlayers[userID].PayTimes[index+1] = getPlayers[userID].PayTimes[index]
						//	getPlayers[userID].PayTimes[index] = pt
						//} else {
						///		getPlayers[userID].PayTimes[index+1] = pt
						//	}
						//}
					}
				}
			}
		}

	}
	var sendPayInfos []*models.PayInfo
	for userID, v := range getPlayers {
		if v != nil {
			data := v.GetSendData()
			//fmt.Println("data.EndTime= ", data.EndTime)
			if data.EndTime < startTime {
				continue
			}

			// if data.EndTime < startTime {
			// 	continue
			// }
			works := models.WorkMap[userID]
			var diff int64
			var lastTime int64
			var work *models.Work
			if works != nil {
				nowTime := time.Now().Unix()
				for _, w := range works {
					if w != nil {

						if diff == 0 {
							diff = nowTime - w.CreateTime
							lastTime = w.CreateTime
							work = w
						} else if diff > nowTime-w.CreateTime {
							diff = nowTime - w.CreateTime
							lastTime = w.CreateTime
							work = w
						}
					}
				}

				lastName := ""
				if work != nil {
					course := models.CourseMap[work.CourseID]

					if course != nil {
						lastName = course.Title
					}
					mission := models.MissionMap[work.MissionID]

					if mission != nil {
						lastName = fmt.Sprintf("%s,%s", lastName, mission.Title)
					}
				}

				data.P_last_name = lastName
				data.P_last_submit = utils.GetTimeString(utils.ParseUnixTime(lastTime))
			}
			data.P_addDay = models.GetAddDay(startTime, data.EndTime)
			sendPayInfos = append(sendPayInfos, data)
		}
	}

	return sendPayInfos
}

func GetPayInfo(startTime, endTime int64) []*models.PayInfo {
	getPlayers := make(map[int]*models.PayInfo)
OUTLOOP:
	for userID, v := range models.OrderSummaryMap {

		// if userID == 7 {
		// 	fmt.Println("len(v) = ", len(v))
		// }
		if v != nil {
			for sID := range v {
				order := models.OrderTotalMap[sID]
				// if userID == 7 {
				// 	fmt.Println("order = ", order)
				// }
				if order != nil {
					// if userID == 7 {
					// 	fmt.Println(" order.Start  = ", order.Created)
					// 	fmt.Println(" order.GetEndTime()  = ", order.GetEndTime())
					// 	fmt.Println(" endTime  = ", endTime)
					// }
					if order.GetEndTime() > endTime {

						if getPlayers[userID] != nil {
							delete(getPlayers, userID)
						}
						continue OUTLOOP
					} else {
						if getPlayers[userID] == nil {
							getPlayers[userID] = &models.PayInfo{
								P_id: userID,
								//PayTimes: make(map[int]*models.PayTimes),
							}

							user := models.UserAllMap[userID]

							if user != nil {
								getPlayers[userID].P_name = user.Nickname
								getPlayers[userID].P_phone = user.Mobile
								getPlayers[userID].P_mail = user.Email
								getPlayers[userID].P_aliPay = user.Alipay_nickname

								if user.Delete != "" {
									getPlayers[userID].P_delete = 1
								}
							}
						}
						pt := &models.PayTimes{StartTime: order.Created, EndTime: order.GetEndTime()}
						if getPlayers[userID].PayTimes == nil {
							getPlayers[userID].PayTimes = pt
							continue
						}
						node := getPlayers[userID].PayTimes
						var preNode *models.PayTimes
						for node != nil {
							if node.StartTime < pt.StartTime {
								if node.Next == nil {
									node.Next = pt
									break
								} else {
									preNode = node
									node = node.Next

								}
							} else {
								if preNode == nil {
									pt.Next = node
									getPlayers[userID].PayTimes = pt
									break
								} else {
									preNode.Next = pt
									pt.Next = node
									break
								}
							}
						}
						/*						pt := &models.PayTimes{StartTime: order.Created, EndTime: order.GetEndTime()}

												if getPlayers[userID].P_name == "海底捞" {
													fmt.Println(" order.GetEndTime() ", order.GetEndTime())
												}
												if len(getPlayers[userID].PayTimes) == 0 {
													getPlayers[userID].PayTimes[1] = pt
													continue
												}

												for index := len(getPlayers[userID].PayTimes); index >= 1; index-- {
													if getPlayers[userID].PayTimes[index].StartTime > pt.StartTime {
														getPlayers[userID].PayTimes[index+1] = getPlayers[userID].PayTimes[index]
														getPlayers[userID].PayTimes[index] = pt
													} else {
														getPlayers[userID].PayTimes[index+1] = pt
													}
												}*/

					}
				}
			}
		}

	}
	var sendPayInfos []*models.PayInfo
	for userID, v := range getPlayers {
		if v != nil {
			data := v.GetSendData()
			//fmt.Println("data.EndTime= ", data.EndTime)
			if data.EndTime > endTime {
				continue
			}

			if data.EndTime < startTime {
				continue
			}
			works := models.WorkMap[userID]
			var diff int64
			var lastTime int64
			var work *models.Work
			if works != nil {
				nowTime := time.Now().Unix()
				for _, w := range works {
					if w != nil {

						// if nowTime-w.CreateTime <= startTime {
						// 	continue
						// }

						// if nowTime-w.CreateTime > endTime {
						// 	continue
						// }
						if w.Images == "passDirectly" {
							continue
						}
						if diff == 0 {
							diff = nowTime - w.CreateTime
							lastTime = w.CreateTime
							work = w
						} else if diff > nowTime-w.CreateTime {
							diff = nowTime - w.CreateTime
							lastTime = w.CreateTime
							work = w
						}
					}
				}

				lastName := ""
				if work != nil {
					course := models.CourseMap[work.CourseID]

					if course != nil {
						lastName = course.Title
					}
					mission := models.MissionMap[work.MissionID]

					if mission != nil {
						lastName = fmt.Sprintf("%s,%s", lastName, mission.Title)
					}
				}

				data.P_last_name = lastName
				data.P_last_submit = utils.GetTimeString(utils.ParseUnixTime(lastTime))
			}
			sendPayInfos = append(sendPayInfos, data)
		}
	}

	return sendPayInfos
}

func GetPayRecord(startTime, endTime int64) []*models.OrderRecord {
	var result []*models.OrderRecord
	for userID, v := range models.OrderSummaryMap {

		if v != nil {
			for sID, os := range v {
				order := models.OrderTotalMap[sID]
				// if userID == 7 {
				// 	fmt.Println("order = ", order)
				// }
				if order != nil {
					if os.Created_at > endTime || os.Created_at < startTime {
						continue
					}

					user := models.UserAllMap[userID]
					uName := ""
					ali := ""
					if user != nil {
						if user.Privilege == 50 {
							continue
						}
						uName = user.Nickname
						ali = user.Alipay_nickname
					}

					userTagRels := models.UserTagRelMap[userID]
					tagRelString := ""
					for _, rel := range userTagRels {
						if rel != nil {
							tag := models.UserTagMap[rel.TagID]
							if tag != nil {
								if tagRelString == "" {
									tagRelString = tag.Name
								} else {
									tagRelString = fmt.Sprintf("%s;%s", tagRelString, tag.Name)
								}
							}
						}
					}

					or := &models.OrderRecord{
						OR_id:          userID,
						OR_name:        uName,
						OR_user_type:   models.GetUserOrderType(userID),
						OR_course:      models.GetUserExpCourses(userID),
						OR_tag:         tagRelString,
						OR_type:        order.Period,
						OR_price:       order.Price,
						OR_create_time: utils.GetTimeString(utils.ParseUnixTime(os.Created_at)),
						OR_start_time:  utils.GetTimeString(utils.ParseUnixTime(order.Created)),
						OR_aliPay:      ali,
					}
					result = append(result, or)
				}
			}
		}
	}
	return result
}

func GetLazeUserRecord(startTime1, endTime1, startTime2, endTime2 int64, numOfDay int) []*models.LazeUserRecord {
	var lrRecords []*models.LazeUserRecord
	for _,v := range models.UserAllMap {
		if v != nil && v.CheckUserIsPayUser(){
			lastOrderTime, ok := v.CheckUserLastOrderInTime(startTime2, endTime2)

			if !ok {
				continue
			}

			ws := models.WorkMap[v.ID]
			var compareWS []*models.Work
			workTimes := 0
			workDetails := ""

			lrRecord := &models.LazeUserRecord {
				LR_id : v.ID,
				LR_name : v.Nickname,
				LR_user_type: models.GetUserOrderType(v.ID),
				LR_date: utils.GetTimeString(utils.ParseUnixTime(lastOrderTime)),
				LR_leave: v.CheckUserInLeave(),
				LR_course: v.GetAllCourses(),
				// LR_last string

				// LR_times int
				// LR_work string
			}
			if ws == nil || len(ws) == 0 {
				lrRecord.LR_last = "否"
				lrRecord.LR_times = workTimes
				lrRecord.LR_work = workDetails
			} else if len(ws) == 1 {
				var maxTime int64
				for _, w := range ws {
					maxTime = time.Now().Unix() - w.CreateTime
				}
				
				if maxTime >= int64(numOfDay * 86400) {
					lrRecord.LR_last = "是"
				} else {
					lrRecord.LR_last = "否"
				}
				lrRecord.LR_times = workTimes
				lrRecord.LR_work = workDetails			
			} else {
				for _, w := range ws {
					if w.CreateTime < startTime1 || w.CreateTime > endTime1 {
						continue
					}
					workTimes++
					insertIndex := -1
					for i, cw := range compareWS {
						diff := int(w.CreateTime - cw.CreateTime)
						
						if diff > 0 {
							// if minTime == 0 {
							// 	minTime = diff
							// } else if diff < minTime {
							// 	minTime = diff
							// }
							insertIndex = i
						} else {
							break
						}

					}
					if insertIndex == -1 {
						compareWS = append([]*models.Work{w}, compareWS...)
					} else {
						compareWS = append(compareWS[:insertIndex + 1], append([]*models.Work{w}, compareWS[insertIndex + 1:]...)...)
					}
					
					// if v.ID == 2912 {
					// 	fmt.Println("insertIndex = " , insertIndex)
					// 	for _, wd := range compareWS {
					// 		if wd != nil {
					// 			fmt.Println(utils.GetTimeString(utils.ParseUnixTime(wd.CreateTime)))
								
					// 		}
					// 	}
					// 	fmt.Println("==========================")
					// }
					

				}
				var maxTime int64
				var LastTime int64
				for _, wd := range compareWS {
					if wd != nil {
						if workDetails == "" {
							workDetails = wd.GetWorkRecordString()
						} else {
							workDetails = fmt.Sprintf("%s;%s", workDetails, wd.GetWorkRecordString())
						}
						if LastTime == 0 {
							LastTime = wd.CreateTime
						} else {
							diff := wd.CreateTime - LastTime
							LastTime = wd.CreateTime
							if diff > maxTime {
								maxTime = diff
							}
						}
					}
				}

				if maxTime == 0 {
					continue
				}
				if maxTime >= int64(numOfDay * 86400) {
					lrRecord.LR_last = "是"
				} else {
					lrRecord.LR_last = "否"
				}
				lrRecord.LR_times = workTimes
				lrRecord.LR_work = workDetails
			}

			lrRecords = append(lrRecords, lrRecord)
		}
	}
	return lrRecords;
}

//GetPayRate 获取转换率
func GetPayRate(startTime, endTime int64) ([]*models.RepayInfo, int, int) {
	var infos []*models.RepayInfo
	totalNum := 0
	totalPayNum := 0 
	selectTotalNum := 0
	selectPayNum := 0

	for _, v := range models.UserExperienceAllMap {
		if v != nil {
			user := models.UserWithDelete[v.UserID]

			if user == nil {
				continue
			}
			totalNum++
			if v.Started >= startTime && v.Started <= endTime {

				selectTotalNum++
				ri := &models.RepayInfo{}
				ri.P_id = user.ID
				ri.P_name = user.Nickname
				ri.P_aliPay = user.Alipay_nickname
				ri.P_E_start = utils.GetTimeString(utils.ParseUnixTime(v.Started))
				ri.P_create = utils.GetTimeString(utils.ParseUnixTime(user.Created_at))

				if models.UserCourseMap[user.ID] != nil {
					for _, v1 := range models.UserCourseMap[user.ID] {
						if v1 != nil {
							c := models.CourseMap[v1.CourseID]

							if c != nil {
								if ri.P_course == "" {
									ri.P_course = c.Title
								} else {
									ri.P_course = fmt.Sprintf("%s;%s", ri.P_course, c.Title)
								}
							}
						}
					}
				}
				// if models.CourseMap[v.CourseID] != nil {
				// 	ri.P_course = models.CourseMap[v.CourseID].Title
				// }

				userTagRels := models.UserTagRelMap[v.UserID]
				var lastestTagRel *models.UserTagRel
				for _, rel := range userTagRels {
					if lastestTagRel == nil || lastestTagRel.Created < rel.Created {
						lastestTagRel = rel
					}
				}

				if lastestTagRel != nil {
					tag := models.UserTagMap[lastestTagRel.TagID]
					if tag != nil {
						ri.P_tag = tag.Name
					}
				}

				oss := models.OrderSummaryMap[v.UserID]

				if oss != nil {
					for _, o := range oss {
						if o != nil {
							order := models.OrderTotalMap[o.ID]
							if order != nil {
								ri.P_type = order.Period
								ri.P_price = order.Price
								ri.P_start = utils.GetTimeString(utils.ParseUnixTime(order.Created))

								selectPayNum++
								totalPayNum++
								break
							}
						}
					}
				}
				infos = append(infos, ri)
			} else {
				oss := models.OrderSummaryMap[v.UserID]

				if oss != nil {
					for _, o := range oss {
						if o != nil {
							order := models.OrderMap[o.ID]
							if order != nil {
								totalPayNum++
								break
							}
						}
					}
				}
			}
		}
	}
	totalRate := 0

	if totalNum > 0 {
		totalRate = totalPayNum * 100 / totalNum
	}

	selecRate := 0

	if selectTotalNum > 0 {
		selecRate = selectPayNum * 100 / selectTotalNum
	}

	return infos, totalRate, selecRate
}

//GetRePayInfo 获取复购率
func GetRePayInfo(startTime, endTime int64) ([]*models.RepayInfo, int, int) {
	var repayInfos []*models.RepayInfo
	totalRepeatNum := 0
	totalNum := 0
	for userID, v := range models.OrderSummaryMap {
		if v != nil {
			// isRepeated := false
			// isBuy := false
			//var orders []*models.Order
			totalNum++
			// orderNum := 0
			var orders []*models.Order
			var firstOrderSummary *models.OrderSummary
			orderSummarys := make(map[int]*models.OrderSummary)
			for _, os := range v {
				order := models.OrderTotalMap[os.ID]

				if order != nil {

					if os.Created_at < startTime {
						// isRepeated = true
						firstOrderSummary = os
					} else if os.Created_at >= startTime && os.Created_at <= endTime {

						// isBuy = true
						if firstOrderSummary == nil {
							firstOrderSummary = os
						} else if firstOrderSummary.Created_at > os.Created_at {
							firstOrderSummary = os
						}
						// orderNum++
						orders = append(orders, order)
						orderSummarys[order.SummaryID] = os
						// if order1 == nil {
						// 	order1 = order
						// 	orderSummary = os
						// } else {
						// 	if order1.Created < order.Created {
						// 		order1 = order
						// 		orderSummary = os
						// 	}
						// }
					}
				}
			}

			// if orderNum > 1 {
			// 	isRepeated = true
			// }

			// if isBuy && isRepeated {
			// 	totalRepeatNum++
			// 	user := models.UserWithDelete[userID]
			// 	if user != nil {
			// 		ri := &models.RepayInfo{}
			// 		ri.P_id = user.ID
			// 		ri.P_name = user.Nickname
			// 		ri.P_type = order1.Period
			// 		ri.P_price = order1.Price
			// 		ri.P_create = utils.GetTimeString(utils.ParseUnixTime(orderSummary.Created_at))
			// 		ri.P_start = utils.GetTimeString(utils.ParseUnixTime(order1.Created))
			// 		ri.P_aliPay = user.Alipay_nickname

			// 		repayInfos = append(repayInfos, ri)
			// 	}
			// }

			for _, v1 := range orders {
				os := orderSummarys[v1.SummaryID]
				if os == firstOrderSummary {
					continue
				}
				user := models.UserWithDelete[userID]
				if user != nil {
					ri := &models.RepayInfo{}
					ri.P_id = user.ID
					ri.P_name = user.Nickname
					ri.P_type = v1.Period
					ri.P_price = v1.Price
					ri.P_create = utils.GetTimeString(utils.ParseUnixTime(os.Created_at))
					ri.P_start = utils.GetTimeString(utils.ParseUnixTime(v1.Created))
					ri.P_aliPay = user.Alipay_nickname

					repayInfos = append(repayInfos, ri)
				}
			}
		}
	}
	return repayInfos, totalRepeatNum, totalNum
}

//GetRatingAvgTimes 获取平均时间
func GetRatingAvgTimes() map[int]string {
	timeMap := make(map[int]int)
	numMap := make(map[int]int)
	returnMap := make(map[int]string)
	for _, v := range models.RatingMap {
		if v != nil {
			if v.RatingDuration > 210000000 || v.RatingDuration <= 0 {
				continue
			}
			timeMap[v.MissionID] += v.RatingDuration
			numMap[v.MissionID]++
		}
	}

	for k, v := range timeMap {
		if numMap[k] > 0 {
			timeMap[k] = v / numMap[k]
			second := timeMap[k] / 1000
			returnMap[k] = fmt.Sprintf("%d", second)
		}
	}

	return returnMap
}

//GetTeacherStatus 获取点评师状态
func GetTeacherStatus(startTime, endTime int64, courseID, userID int) []*models.TeacherInfo {

	avgMap := GetRatingAvgTimes()
	var sendTeacherInfos []*models.TeacherInfo
	if courseID == 0 && userID == 0 {
		//全部的
		for _, v := range models.RatingMap {
			if v != nil {
				if v.UpdateTime < startTime || v.UpdateTime > endTime {
					continue
				}
				if v.RatingDuration <= 0 {
					continue
				}
				if v.RatingDuration > 2100000000 {
					continue
				}
				teacher := models.TeacherMaps[v.UserID]

				if teacher == nil {
					continue
				}

				ti := &models.TeacherInfo{}

				ti.T_name = teacher.Nickname
				uTime := utils.ParseUnixTime(v.UpdateTime)
				createTimes := strings.Split(utils.GetTimeString(uTime), " ")
				ti.T_createDate = createTimes[0]
				if len(createTimes) > 1 {
					ti.T_createTime = createTimes[1]
				}
				month, monthWeek := utils.GetMonthWeek(uTime)
				ti.T_month = models.GetMonthString(month)
				ti.T_monthWeek = models.GetMonthWeekString(monthWeek)

				ti.T_weekDay = models.GetWeekDayString(int(uTime.Weekday()))

				ti.T_level = 1
				if v.RatingDuration > 300000 {
					ti.T_level = 2
				}

				if v.RatingDuration > 1800000 {
					ti.T_level = 3
				}
				ratingSecond := v.RatingDuration / 1000

				ti.T_time = fmt.Sprintf("%d", ratingSecond)

				work := models.WorkAllMap[v.WorkID]

				if work != nil {
					user := models.UserAllMap[work.UserID]

					if user != nil {
						ti.T_user_name = user.Nickname
						if user.Delete != "" {
							ti.T_delete = 1
						}
					}
				}

				course := models.CourseMap[v.CourseID]

				if course != nil {
					ti.T_course_name = course.Title
				}

				mission := models.MissionMap[v.MissionID]

				if mission != nil {
					ti.T_mission_name = mission.Title
				}

				ti.T_avg_time = avgMap[v.MissionID]

				sendTeacherInfos = append(sendTeacherInfos, ti)
			}
		}
	} else if courseID == 0 {
		//全部课程
		ratingMap := models.RatingUserMap[userID]
		teacher := models.TeacherMaps[userID]

		if ratingMap != nil && teacher != nil {
			for c, v1 := range models.RatingUserMap[userID] {
				if v1 != nil {

					course := models.CourseMap[c]

					if course == nil {
						continue
					}

					for _, v := range v1 {
						if v != nil {
							if v.UpdateTime < startTime || v.UpdateTime > endTime {
								continue
							}
							if v.RatingDuration <= 0 {
								continue
							}
							if v.RatingDuration > 2100000000 {
								continue
							}

							ti := &models.TeacherInfo{}

							ti.T_name = teacher.Nickname
							uTime := utils.ParseUnixTime(v.UpdateTime)
							createTimes := strings.Split(utils.GetTimeString(uTime), " ")
							ti.T_createDate = createTimes[0]
							if len(createTimes) > 1 {
								ti.T_createTime = createTimes[1]
							}
							month, monthWeek := utils.GetMonthWeek(uTime)
							ti.T_month = models.GetMonthString(month)
							ti.T_monthWeek = models.GetMonthWeekString(monthWeek)
							ti.T_weekDay = models.GetWeekDayString(int(uTime.Weekday()))

							ti.T_level = 1
							if v.RatingDuration > 300000 {
								ti.T_level = 2
							}

							if v.RatingDuration > 1800000 {
								ti.T_level = 3
							}
							ratingSecond := v.RatingDuration / 1000

							ti.T_time = fmt.Sprintf("%d", ratingSecond)

							work := models.WorkAllMap[v.WorkID]

							if work != nil {
								user := models.UserAllMap[work.UserID]

								if user != nil {
									ti.T_user_name = user.Nickname
									if user.Delete != "" {
										ti.T_delete = 1
									}
								}
							}
							ti.T_course_name = course.Title
							mission := models.MissionMap[v.MissionID]

							if mission != nil {
								ti.T_mission_name = mission.Title
							}

							ti.T_avg_time = avgMap[v.MissionID]

							sendTeacherInfos = append(sendTeacherInfos, ti)
						}
					}
				}
			}
		}

	} else if userID == 0 {
		//所有点评师
		ratingMap := models.RatingCourseMap[courseID]
		course := models.CourseMap[courseID]

		if course != nil && ratingMap != nil {
			for _, v := range ratingMap {
				if v != nil {
					if v.UpdateTime < startTime || v.UpdateTime > endTime {
						continue
					}
					if v.RatingDuration <= 0 {
						continue
					}
					if v.RatingDuration > 2100000000 {
						continue
					}

					teacher := models.TeacherMaps[v.UserID]

					if teacher == nil {
						continue
					}

					ti := &models.TeacherInfo{}

					ti.T_name = teacher.Nickname
					uTime := utils.ParseUnixTime(v.UpdateTime)
					createTimes := strings.Split(utils.GetTimeString(uTime), " ")
					ti.T_createDate = createTimes[0]
					if len(createTimes) > 1 {
						ti.T_createTime = createTimes[1]
					}
					month, monthWeek := utils.GetMonthWeek(uTime)
					ti.T_month = models.GetMonthString(month)
					ti.T_monthWeek = models.GetMonthWeekString(monthWeek)
					ti.T_weekDay = models.GetWeekDayString(int(uTime.Weekday()))

					ti.T_level = 1
					if v.RatingDuration > 300000 {
						ti.T_level = 2
					}

					if v.RatingDuration > 1800000 {
						ti.T_level = 3
					}
					ratingSecond := v.RatingDuration / 1000

					ti.T_time = fmt.Sprintf("%d", ratingSecond)

					work := models.WorkAllMap[v.WorkID]

					if work != nil {
						user := models.UserAllMap[work.UserID]

						if user != nil {
							ti.T_user_name = user.Nickname
							if user.Delete != "" {
								ti.T_delete = 1
							}
						}
					}
					ti.T_course_name = course.Title

					mission := models.MissionMap[v.MissionID]

					if mission != nil {
						ti.T_mission_name = mission.Title
					}

					ti.T_avg_time = avgMap[v.MissionID]

					sendTeacherInfos = append(sendTeacherInfos, ti)
				}
			}
		}
	} else {
		//都是特定
		ratingMap := models.RatingUserMap[userID]

		if ratingMap != nil {
			ratingCourseMap := ratingMap[courseID]
			teacher := models.UserAllMap[userID]

			course := models.CourseMap[courseID]
			if ratingCourseMap != nil && teacher != nil && course != nil {

				for _, v := range ratingCourseMap {
					if v != nil {
						if v.UpdateTime < startTime || v.UpdateTime > endTime {
							continue
						}
						if v.RatingDuration <= 0 {
							continue
						}
						if v.RatingDuration > 2100000000 {
							continue
						}

						ti := &models.TeacherInfo{}

						ti.T_name = teacher.Nickname
						uTime := utils.ParseUnixTime(v.UpdateTime)
						createTimes := strings.Split(utils.GetTimeString(uTime), " ")
						ti.T_createDate = createTimes[0]
						if len(createTimes) > 1 {
							ti.T_createTime = createTimes[1]
						}
						month, monthWeek := utils.GetMonthWeek(uTime)
						ti.T_month = models.GetMonthString(month)
						ti.T_monthWeek = models.GetMonthWeekString(monthWeek)
						ti.T_weekDay = models.GetWeekDayString(int(uTime.Weekday()))

						ti.T_level = 1
						if v.RatingDuration > 300000 {
							ti.T_level = 2
						}

						if v.RatingDuration > 1800000 {
							ti.T_level = 3
						}
						ratingSecond := v.RatingDuration / 1000

						ti.T_time = fmt.Sprintf("%d", ratingSecond)

						work := models.WorkAllMap[v.WorkID]

						if work != nil {
							user := models.UserAllMap[work.UserID]

							if user != nil {
								ti.T_user_name = user.Nickname
								if user.Delete != "" {
									ti.T_delete = 1
								}
							}
						}
						ti.T_course_name = course.Title

						mission := models.MissionMap[v.MissionID]

						if mission != nil {
							ti.T_mission_name = mission.Title
						}

						ti.T_avg_time = avgMap[v.MissionID]

						sendTeacherInfos = append(sendTeacherInfos, ti)
					}
				}
			}
		}
	}

	return sendTeacherInfos
}

func GetUserCourseDetail(user *models.User, year, month int) ([]*models.UserDetail, []*models.CourseData, int) {

	var returnData []*models.UserDetail
	courseDatas := models.CreateMonthCourseData()
	monthTotal := 0
	userCs := models.WorkCourseMap[user.ID]

	year = utils.GetCurrentYear() - year
	isAll := false
	if month == 0 {
		isAll = true
	}

	for k, v := range userCs {
		cour := models.CourseMap[k]

		if cour == nil {
			continue
		}
		ud := &models.UserDetail{
			ID:         k,
			CourseName: cour.Title,
		}
		if isAll {
			monthDatas := make([]*models.CourseData, 12)

			for _, v1 := range v {
				if v1 != nil {
					var date1 = utils.ParseUnixTime(v1.CreateTime)
					var m = int(date1.Month())
					var y = date1.Year()
					if y != year {
						continue
					}
					if monthDatas[m-1] == nil {
						monthDatas[m-1] = &models.CourseData{
							Month: m,
							Num:   1,
						}
					} else {
						monthDatas[m-1].Num = monthDatas[m-1].Num + 1
					}
					ud.TotalNum = ud.TotalNum + 1
					courseDatas[m-1].Num = courseDatas[m-1].Num + 1
					monthTotal = monthTotal + 1
				}
			}

			var datas []*models.CourseData

			for _, v := range monthDatas {
				if v != nil {
					datas = append(datas, v)
				}
			}
			ud.CourseDatas = datas
		} else {
			monthDatas := make([]*models.CourseData, 1)

			for _, v1 := range v {
				if v1 != nil {
					var date1 = utils.ParseUnixTime(v1.CreateTime)
					var m = int(date1.Month())
					var y = date1.Year()
					if y != year || m != month {
						continue
					}
					if monthDatas[0] == nil {
						monthDatas[0] = &models.CourseData{
							Month: m,
							Num:   1,
						}
					} else {
						monthDatas[0].Num = monthDatas[0].Num + 1
					}
					ud.TotalNum = ud.TotalNum + 1
					courseDatas[m-1].Num = courseDatas[m-1].Num + 1
					monthTotal = monthTotal + 1
				}
			}
			var datas []*models.CourseData

			for _, v := range monthDatas {
				if v != nil {
					datas = append(datas, v)
				}
			}
			ud.CourseDatas = datas
		}
		returnData = append(returnData, ud)
	}
	// fmt.Println("returnData = ", len(returnData))
	// for _, v := range returnData {
	// 	for _, v1 := range v.CourseDatas {
	// 		if v1 == nil {
	// 			continue
	// 		}
	// 		fmt.Println("v1.Month = ", v1.Month)
	// 	}
	// }
	return returnData, courseDatas, monthTotal
}
