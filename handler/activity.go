package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"awesomeProject/model"
	"awesomeProject/service"
	"awesomeProject/utils"
)

func GetActivityById(c *gin.Context) {
	aid := c.Query("id")
	if aid == "" {
		utils.BuildErrorResponse(c, 400, "aid 为必传参数 请传递aid")
		return
	}
	aidInt, err := strconv.ParseInt(aid, 10, 64)
	if err != nil {
		utils.BuildErrorResponse(c, 400, "aid is not valid")
		return
	}

	res, err := service.GetActivityByAid(aidInt)
	if err != nil {
		utils.BuildErrorResponse(c, 500, "GetActivityById err err is: "+err.Error())
		return
	}
	utils.BuildSuccessResponse(c, res)
}

// 上传活动
func AddActivity(c *gin.Context) {
	var addActivityReq *model.AddActivityReq
	err := c.ShouldBindJSON(&addActivityReq)
	if err != nil {
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

// 查询活动列表
func GetActivityList(c *gin.Context) {
	pageSize := c.Query("pageSize")
	pageNum := c.Query("pageNum")

	if pageSize != "" && pageNum != "" {
		// 如果指定了分页参数
		pageSizeInt, err := strconv.Atoi(pageSize)
		if err != nil {
			utils.BuildErrorResponse(c, 400, "pageSize not valid")
			return
		}
		pageNumInt, err := strconv.Atoi(pageNum)
		if err != nil {
			utils.BuildErrorResponse(c, 400, "pageNum not valid")
			return
		}

		res, total, err := service.GetActivityList(&pageSizeInt, &pageNumInt)
		if err != nil {
			utils.BuildErrorResponse(c, 500, "GetActivityList err err is: "+err.Error())
			return
		}
		utils.BuildSuccessResponse(c, gin.H{
			"activities": res,
			"total":      total,
		})
		return
	} else {
		utils.BuildErrorResponse(c, 400, "pageSize or pageNum is not valid")
		return
	}
}
