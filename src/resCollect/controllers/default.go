package controllers

import (
	"strings"
	"fmt"
	"resCollect/models"
	"resCollect/utils"

	"github.com/astaxie/beego"
)

const (
	IGNORE_CSV_FILE_NAME = "AssetBundleIgnore.csv"
	CMD_GAME_INIT        = 1000 //初始化
	CMD_RESET_RESOURCE   = 1001 //重置资源
	CMD_DOWNLOAD_TABLE   = 1002 //下载资源
	CMD_COL_RESOURCE     = 1003 //收集资源
	CMD_DEL_RESOURCE     = 1004 //删除资源
)

type MainController struct {
	beego.Controller
}

//ResourceController 资源收集controller
type ResourceController struct {
	beego.Controller
}

//ResourceShowController 资源展示页面
type ResourceShowController struct {
	beego.Controller
}

func (c *MainController) Get() {
	c.Data["Website"] = "beego.me"
	c.Data["Email"] = "astaxie@gmail.com"
	c.TplName = "index.tpl"
}

//Get 资源收集
func (c *ResourceController) Get() {

	cmd, _ := c.GetInt("cmd")
	fmt.Println("<CMD> = ", cmd)
	switch cmd {
	case CMD_GAME_INIT:
		quality := c.GetString("quality")
		platform := c.GetString("platform")
		clientVersion, _ := c.GetInt("clientVersion")
		// clientABVersion, _ := c.GetInt("clientABVersion")
		isReset, _ := c.GetInt("isReset")

		models.GlobalResCol.ResetWithParams(quality, platform, clientVersion, isReset)

		//添加默认需要的资源
		models.GlobalResCol.AddRes(0, "Assets/Resources/Prefabs", "UI_loading_shalong2")
		models.GlobalResCol.AddRes(0, "Assets/Resources/UIImage/Summon", "CommonBG")
		models.GlobalResCol.AddRes(0, "Assets/Resources/High/Avatar", "黑人_黑人")
	case CMD_COL_RESOURCE:
		resName := c.GetString("resName")
		resPath := c.GetString("resPath")
		clientABVersion, _ := c.GetInt("clientABVersion")
		if resName == "" || resPath == "" {
			return
		}
		if strings.Contains(resName, "Low/Avatar") {
			resPathHigh := strings.ReplaceAll(resPath, "Low/Avatar", "High/Avatar")
			models.GlobalResCol.AddRes(clientABVersion, resPathHigh, resName)
		} else if strings.Contains(resName, "High/Avatar") {
			resPathLow := strings.ReplaceAll(resPath, "High/Avatar", "Low/Avatar")
			models.GlobalResCol.AddRes(clientABVersion, resPathLow, resName)
		}
		
		models.GlobalResCol.AddRes(clientABVersion, resPath, resName)
	}
	c.TplName = "index.tpl"
	c.Ctx.WriteString("OK")
}

//Get 展示资源列表
func (c *ResourceShowController) Get() {
	isDownload, _ := c.GetInt("isDownload")
	if isDownload == 1 {
		fmt.Println("isDownload = ", isDownload)
		filePath := utils.CreateCSVFile(IGNORE_CSV_FILE_NAME, models.GlobalResCol.GetReses())
		fmt.Println("filePath = ", filePath)
		c.Ctx.Output.Download(filePath)
	}

	c.Data["totalSize"] = models.GlobalResCol.GetTotalResSize()
	c.Data["clientVerison"] = models.GlobalResCol.GetClientVerison()
	c.Data["resources"] = models.GlobalResCol.GetReses()
	c.Data["Index"] = 1
	c.TplName = "resourceShow.html"
}

//Post 通过cmd来执行对应的逻辑
func (c *ResourceShowController) Post() {
	cmd, err := c.GetInt("cmd")

	if err != nil {
		fmt.Println("err = ", err)
		return
	}

	switch cmd {
	case CMD_RESET_RESOURCE:
		models.GlobalResCol.Reset()
		c.Data["Index"] = 1
		c.Data["resources"] = models.GlobalResCol.GetReses()
		c.TplName = "resourceShow.html"
	case CMD_DOWNLOAD_TABLE:
		// filePath := utils.CreateCSVFile(IGNORE_CSV_FILE_NAME, models.GlobalResCol.GetReses())
		// fmt.Println("filePath = ", filePath)
		// // c.Ctx.Output.Download(filePath)
		// c.Data["filePath"] = filePath

		// c.Ctx.Redirect(302, "/")
		c.Data["resources"] = models.GlobalResCol.GetReses()
		c.Data["Index"] = 1
		c.TplName = "resourceShow.html"
	case CMD_DEL_RESOURCE:
		id := c.GetString("id")
		// fmt.Println("id = ", id)
		models.GlobalResCol.RemoveRes(id)
		totalSize := models.GlobalResCol.GetTotalResSize()
		c.Ctx.WriteString(fmt.Sprintf("%s|%d", id, totalSize))
		// c.Data["id"] = id
		// c.Data["Index"] = 1
		// c.TplName = "resourceShow.html"
	}
}
