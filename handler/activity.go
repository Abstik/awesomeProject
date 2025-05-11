package handler

import (
	"awesomeProject/model"
	"awesomeProject/service"
	"awesomeProject/utils"
	"github.com/gin-gonic/gin"
	"strconv"
)

func GetActivityList(c *gin.Context) {
	pageSize := c.Query("pageSize")
	pageNum := c.Query("pageNum")
	var res []*model.ActivityPO
	if pageSize != "" && pageNum != "" {
		pageSizeInt, err := strconv.Atoi(pageSize)
		if err != nil {
			utils.BuildErrorResponse(c, 500, "pageSize not valid")
			return
		}
		pageNumInt, err := strconv.Atoi(pageNum)
		if err != nil {
			utils.BuildErrorResponse(c, 500, "pageNum not valid")
			return
		}
		res, err = service.GetActivityList(&pageSizeInt, &pageNumInt)
		if err != nil {
			utils.BuildErrorResponse(c, 500, "GetActivityList err err is: "+err.Error())
			return
		}
	} else {
		var err error
		res, err = service.GetActivityList(nil, nil)
		if err != nil {
			utils.BuildErrorResponse(c, 500, "GetActivityList err err is: "+err.Error())
			return
		}
	}
	utils.BuildSuccessResponse(c, res)
}

func GetActivityById(c *gin.Context) {
	aid := c.Query("id")
	if aid == "" {
		utils.BuildErrorResponse(c, 500, "aid 为必传参数 请传递aid")
		return
	}
	aidInt, err := strconv.ParseInt(aid, 10, 64)
	if err != nil {
		utils.BuildErrorResponse(c, 500, "aid is not valid")
		return
	}
	res, err := service.GetActivityByAid(aidInt)
	if err != nil {
		utils.BuildErrorResponse(c, 500, "GetActivityById err err is: "+err.Error())
		return
	}
	utils.BuildSuccessResponse(c, res)
}

func AddActivity(c *gin.Context) {
	var addActivityReq *model.AddActivityReq
	err := c.ShouldBindJSON(&addActivityReq)
	if err != nil {
		utils.Logger.Errorf("handler.AddActivity format error, error is %s", err.Error())
		utils.BuildErrorResponse(c, 500, "AddActivity format error, error is "+err.Error())
		return
	}
	err = service.AddActivity(addActivityReq)
	if err != nil {
		utils.BuildErrorResponse(c, 500, "AddActivity Failed error is "+err.Error())
		return
	}
	utils.BuildSuccessResponse(c, "添加成功")
}
