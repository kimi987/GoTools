package models

//SetXLSFile 设置文件
import (
	"dgServer/utils"
	"fmt"
	"strings"

	"github.com/tealeg/xlsx"
)

//SetXLSFile 设置下载excel
func SetXLSFile(fileName, startTime, endTime string, fileDatas ...interface{}) string {
	file := xlsx.NewFile()
	sheet, err := file.AddSheet("Sheet1")
	if err != nil {
		fmt.Printf(err.Error())
	}
	if startTime != "" && endTime != "" {
		fileName = fmt.Sprintf("%s-%s-%s.xlsx", fileName, startTime, endTime)
		fileName = strings.ReplaceAll(fileName, "/", "")
		row := sheet.AddRow()
		addCell(row, startTime)
		addCell(row, endTime)
	}
	datas := fileDatas[0]
	if userStatus, ok := datas.([]*UserStatus); ok {
		for n, v := range userStatus {
			if v == nil {
				continue
			}

			if n == 0 {
				n++

				row := sheet.AddRow()
				row.SetHeightCM(1)
				addCell(row, "编号")
				addCell(row, "用户名")
				addCell(row, "几天未交作业")
				addCell(row, "所开通道")
				addCell(row, "最后一次提交时间")
				addCell(row, "最后一次通道任务名字")
				addCell(row, "点评师")
			}

			row := sheet.AddRow()
			row.SetHeightCM(1)
			addCell(row, fmt.Sprintf("%d", v.M_id))
			if v.M_leave == 1 {
				addCell(row, fmt.Sprintf("%s%s", v.M_name, "(请假中)"))
			} else {
				addCell(row, v.M_name)
			}

			addCell(row, fmt.Sprintf("%d", v.M_day))
			addCell(row, v.M_course)
			addCell(row, v.M_last_submit)
			addCell(row, v.M_last_name)
			addCell(row, v.M_teacher)
		}
	} else if simpleUsers, ok := datas.([]*SimpleUser); ok {
		for n, v := range simpleUsers {
			if v == nil {
				continue
			}

			if n == 0 {
				n++
				row := sheet.AddRow()
				row.SetHeightCM(1)
				addCell(row, "编号")
				addCell(row, "用户名")
			}

			row := sheet.AddRow()
			row.SetHeightCM(1)
			addCell(row, fmt.Sprintf("%d", v.S_User_id))
			addCell(row, v.S_name)
		}
	} else if WorkTimes, ok := datas.(map[int]map[int]map[int]*WorkTimes); ok {
		n := 0
		for _, v:= range WorkTimes {
			if v == nil {
				continue
			}
			for _, v1 := range v {
				if v1 == nil {
					continue
				}
				for _,v2 := range v1 {
					if v2 == nil {
						continue
					}
					if n == 0 {
						n++
						row := sheet.AddRow()
						row.SetHeightCM(1)
						addCell(row, "用户名")
						addCell(row, "通道")
						addCell(row, "任务")
						addCell(row, "数量")
					}
					row := sheet.AddRow()
					row.SetHeightCM(1)
					addCell(row, v2.W_user_name)
					addCell(row, v2.W_course)
					addCell(row, v2.W_mission)
					addCell(row, fmt.Sprintf("%d", v2.W_num))
				}
				
			}
		}
	} else if taskScores, ok := datas.([]*SendTaskScore); ok {
		for n, v := range taskScores {
			if v == nil {
				continue
			}
			if n == 0 {
				n++

				row := sheet.AddRow()
				row.SetHeightCM(1)
				addCell(row, "提交时间")
				addCell(row, "编号")
				addCell(row, "用户名")
				addCell(row, "通道名")
				addCell(row, "任务名")
				addCell(row, "评分时间")
				addCell(row, "评分等级")
			}

			row := sheet.AddRow()
			row.SetHeightCM(1)
			addCell(row, v.SD_utime)
			addCell(row, fmt.Sprintf("%d", v.SD_id))
			addCell(row, v.SD_username)
			addCell(row, v.SD_course_name)
			addCell(row, v.SD_name)
			addCell(row, v.SD_stime)
			addCell(row, fmt.Sprintf("%d", v.SD_rank))
		}
	} else if payInfos, ok := datas.([]*PayInfo); ok {
		for n, v := range payInfos {
			if v == nil {
				continue
			}
			if n == 0 {
				n++

				row := sheet.AddRow()
				row.SetHeightCM(1)
				addCell(row, "编号")
				addCell(row, "用户名")
				addCell(row, "归档状态")
				addCell(row, "邮箱")
				addCell(row, "电话")
				addCell(row, "到期时间")
				addCell(row, "当前请假")
				addCell(row, "支付宝账号")
				addCell(row, "其他")
			}

			row := sheet.AddRow()
			row.SetHeightCM(1)
			addCell(row, fmt.Sprintf("%d", v.P_id))
			addCell(row, v.P_name)

			if v.P_delete == 0 {
				addCell(row, "正常")
			} else {
				addCell(row, "归档")
			}

			addCell(row, v.P_mail)
			addCell(row, v.P_phone)
			addCell(row, v.P_end_time)
			addCell(row, v.ReleaseTime)
			addCell(row, v.P_aliPay)
			addCell(row, fmt.Sprintf("%d", v.P_addDay))
		}
	} else if orderRecords, ok := datas.([]*OrderRecord); ok {
		for n, v := range orderRecords {
			if v == nil {
				continue
			}
			if n == 0 {
				n++

				row := sheet.AddRow()
				row.SetHeightCM(1)
				addCell(row, "编号")
				addCell(row, "用户名")
				addCell(row, "用户类型")
				addCell(row, "所开主线通道")
				addCell(row, "用户标签")
				addCell(row, "缴费类型")
				addCell(row, "缴费金额")
				addCell(row, "付费时间")
				addCell(row, "开始时间")
				addCell(row, "支付宝账号")
			}

			row := sheet.AddRow()
			row.SetHeightCM(1)
			addCell(row, fmt.Sprintf("%d", v.OR_id))
			addCell(row, v.OR_name)
			addCell(row, v.OR_user_type)
			addCell(row, v.OR_course)
			addCell(row, v.OR_tag)

			if v.OR_type == 3 {
				addCell(row, "季卡")
			} else if v.OR_type == 6 {
				addCell(row, "半年卡")
			} else if v.OR_type == 12 {
				addCell(row, "年卡")
			} else {
				addCell(row, "无效类型")
			}

			addCell(row, fmt.Sprintf("%d", v.OR_price))
			addCell(row, v.OR_create_time)
			addCell(row, v.OR_start_time)
			addCell(row, v.OR_aliPay)
		}
	} else if lrRecords, ok := datas.([]*LazeUserRecord); ok {
		for n, v := range lrRecords {
			if v == nil {
				continue
			}
			if n == 0 {
				n++

				row := sheet.AddRow()
				row.SetHeightCM(1)
				addCell(row, "编号")
				addCell(row, "用户名")
				addCell(row, "用户类型")
				addCell(row, "最近缴费日期")
				addCell(row, "连续未交")
				addCell(row, "请假状态")
				addCell(row, "所开主线通道")
				addCell(row, "提交次数")
				addCell(row, "提交作业详情")
			}

			row := sheet.AddRow()
			row.SetHeightCM(1)
			addCell(row, fmt.Sprintf("%d", v.LR_id))
			addCell(row, v.LR_name)
			addCell(row, v.LR_user_type)
			addCell(row, v.LR_date)
			addCell(row, v.LR_last)
			addCell(row, v.LR_leave)
			addCell(row,  v.LR_course)
			addCell(row, fmt.Sprintf("%d", v.LR_times))
			addCell(row, v.LR_work)
		}
	} else if repayInfos, ok := datas.([]*RepayInfo); ok {

		repayType, _ := fileDatas[3].(int)
		for n, v := range repayInfos {
			if v == nil {
				continue
			}

			if n == 0 {
				rowStart := sheet.AddRow()
				if repayType == 1 {
					selectRate, _ := fileDatas[1].(int)
					totalRate, _ := fileDatas[2].(int)
					addCell(rowStart, fmt.Sprintf("%d%%所选时段转化率", selectRate))
					addCell(rowStart, fmt.Sprintf("%d%%总转化率", totalRate))
				} else if repayType == 2 {
					repNum, _ := fileDatas[1].(int)
					payRate, _ := fileDatas[2].(float32)
					addCell(rowStart, fmt.Sprintf("%f%%复购率", payRate))
					addCell(rowStart, fmt.Sprintf("%d复购数", repNum))
				}

				row := sheet.AddRow()
				row.SetHeightCM(1)
				addCell(row, "编号")
				addCell(row, "用户名")
				addCell(row, "开启时间")
				addCell(row, "开启通道")
				addCell(row, "用户标签")
				if repayType == 1 {
					addCell(row, "是否购买")
				}
				addCell(row, "年卡类型")
				addCell(row, "年卡金额")
				addCell(row, "付费时间")
				addCell(row, "年卡起始日期")
				addCell(row, "支付宝账号")
			}

			row := sheet.AddRow()
			row.SetHeightCM(1)
			addCell(row, fmt.Sprintf("%d", v.P_id))
			addCell(row, v.P_name)
			addCell(row, v.P_E_start)
			addCell(row, v.P_course)
			addCell(row, v.P_tag)
			if v.P_type == 6 {
				if repayType == 1 {
					addCell(row, "是")
				}

				addCell(row, "半年卡")
			} else if v.P_type == 12 {
				if repayType == 1 {
					addCell(row, "是")
				}
				addCell(row, "一年卡")
			} else if v.P_type == 3 {
				if repayType == 1 {
					addCell(row, "是")
				}
				addCell(row, "季卡")
			} else {
				if repayType == 1 {
					addCell(row, "否")
				}
				addCell(row, "")
			}

			addCell(row, fmt.Sprintf("%d", v.P_price))
			addCell(row, v.P_create)
			addCell(row, v.P_start)
			addCell(row, v.P_aliPay)
		}
	} else if teacherInfos, ok := datas.([]*TeacherInfo); ok {
		for n, v := range teacherInfos {
			if v == nil {
				continue
			}
			if n == 0 {
				n++

				row := sheet.AddRow()
				row.SetHeightCM(1)
				addCell(row, "点评日期")
				addCell(row, "点评时间")
				addCell(row, "第几月")
				addCell(row, "第几周")
				addCell(row, "周几")
				addCell(row, "点评师")
				addCell(row, "异常/超标/达标")
				addCell(row, "用时(秒)")
				addCell(row, "用户名")
				addCell(row, "通道")
				addCell(row, "任务")
				addCell(row, "任务平均时长(秒)")
			}

			row := sheet.AddRow()
			row.SetHeightCM(1)
			addCell(row, v.T_createDate)
			addCell(row, v.T_createTime)
			addCell(row, v.T_month)
			addCell(row, v.T_monthWeek)
			addCell(row, v.T_weekDay)
			addCell(row, v.T_name)

			if v.T_level == 1 {
				addCell(row, "达标")
			} else if v.T_level == 2 {
				addCell(row, "超标")
			} else {
				addCell(row, "异常")
			}

			addCell(row, v.T_time)
			addCell(row, v.T_user_name)
			addCell(row, v.T_course_name)
			addCell(row, v.T_mission_name)
			addCell(row, v.T_avg_time)
		}
	} else if userDetails, ok := datas.([]*UserDetail); ok {
		if courseDatas, ok := fileDatas[1].([]*CourseData); ok {
			if totalNum, ok := fileDatas[2].(int); ok {
				name := fileDatas[3].(string)
				year := fileDatas[4].(int)
				month := fileDatas[5].(int)
				row := sheet.AddRow()
				row.SetHeightCM(1)
				addCell(row, "用户名")
				addCell(row, name)
				row = sheet.AddRow()
				row.SetHeightCM(1)
				addCell(row, "日期")
				if month > 0 {
					addCell(row, fmt.Sprintf("%d年%d月", utils.GetCurrentYear()-year, month))
				} else {
					addCell(row, fmt.Sprintf("%d年", utils.GetCurrentYear()-year))
				}
				sheet.AddRow()
				courseIndex := 0
				for index := 0; index < len(userDetails); index++ {
					row = sheet.AddRow()
					row.SetHeightCM(1)
					addCell(row, fmt.Sprintf("通道名:%s", userDetails[index].CourseName))
					addCell(row, fmt.Sprintf("总提交数:%d", userDetails[index].TotalNum))
					addCell(row, "")
					if courseIndex < len(courseDatas) {
						cd := courseDatas[courseIndex]
						addCell(row, fmt.Sprintf("%d月", cd.Month))
						addCell(row, fmt.Sprintf("总提交数:%d", cd.Num))
						courseIndex++
					} else if courseIndex == len(courseDatas) {
						addCell(row, fmt.Sprintf("年度总提交:%d", totalNum))
						courseIndex++
					}

					for _, v := range userDetails[index].CourseDatas {
						if v != nil {
							row = sheet.AddRow()
							row.SetHeightCM(1)
							addCell(row, fmt.Sprintf("%d月", v.Month))
							addCell(row, fmt.Sprintf("总提交数:%d", v.Num))
							addCell(row, "")
							if courseIndex < len(courseDatas) {
								cd := courseDatas[courseIndex]
								addCell(row, fmt.Sprintf("%d月", cd.Month))
								addCell(row, fmt.Sprintf("总提交数:%d", cd.Num))
								courseIndex++
							} else if courseIndex == len(courseDatas) {
								addCell(row, fmt.Sprintf("年度总提交:%d", totalNum))
								courseIndex++
							}
						}
					}
				}

				for index := courseIndex; index <= len(courseDatas); index++ {
					row = sheet.AddRow()
					addCell(row, "")
					addCell(row, "")
					addCell(row, "")
					if index < len(courseDatas) {
						cd := courseDatas[index]
						addCell(row, fmt.Sprintf("%d月", cd.Month))
						addCell(row, fmt.Sprintf("总提交数:%d", cd.Num))
					} else {
						addCell(row, fmt.Sprintf("年度总提交:%d", totalNum))
					}
				}
			}
		}
	}
	fileName = fmt.Sprintf("%s/%s", "download", fileName)
	err = file.Save(fileName)
	if err != nil {
		fmt.Printf(err.Error())
	}

	return fileName
}

func addCell(row *xlsx.Row, value string) {
	cell := row.AddCell()
	cell.Value = value
}
