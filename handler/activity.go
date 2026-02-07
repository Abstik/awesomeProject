package handler

import (
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"awesomeProject/model"
	"awesomeProject/service"
	"awesomeProject/utils"
)

func GetActivityById(c *gin.Context) {
	aid := c.Query("aid")
	if aid == "" {
		utils.BuildErrorResponse(c, 400, "aid 为必传参数 请传递aid")
		return
	}
	aidInt, err := strconv.ParseInt(aid, 10, 64)
	if err != nil {
		utils.BuildErrorResponse(c, 400, "aid 参数格式有误")
		return
	}

	res, err := service.GetActivityByAid(aidInt)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.BuildErrorResponse(c, 404, "活动不存在")
			return
		}
		utils.BuildServerError(c, "查询活动失败", err)
		return
	}
	utils.BuildSuccessResponse(c, res)
}

// 上传活动
func AddActivity(c *gin.Context) {
	var addActivityReq *model.ActivityReq
	err := c.ShouldBindJSON(&addActivityReq)
	if err != nil {
		utils.BuildErrorResponse(c, 400, "参数格式有误")
		return
	}

	err = service.AddActivity(addActivityReq)
	if err != nil {
		utils.BuildServerError(c, "添加活动失败", err)
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
			utils.BuildErrorResponse(c, 400, "pageSize 参数格式有误")
			return
		}
		pageNumInt, err := strconv.Atoi(pageNum)
		if err != nil {
			utils.BuildErrorResponse(c, 400, "pageNum 参数格式有误")
			return
		}

		res, total, err := service.GetActivityList(pageSizeInt, pageNumInt)
		if err != nil {
			utils.BuildServerError(c, "查询活动列表失败", err)
			return
		}
		utils.BuildSuccessResponse(c, gin.H{
			"activities": res,
			"total":      total,
		})
		return
	} else {
		utils.BuildErrorResponse(c, 400, "pageSize 和 pageNum 为必传参数")
		return
	}
}

func UpdateActivity(c *gin.Context) {
	var activityReq *model.ActivityReq
	err := c.ShouldBindJSON(&activityReq)
	if err != nil {
		utils.BuildErrorResponse(c, 400, "参数格式有误")
		return
	}
	if activityReq.AID == 0 {
		utils.BuildErrorResponse(c, 400, "aid 参数无效")
		return
	}

	err = service.UpdateActivity(activityReq)
	if err != nil {
		utils.BuildServerError(c, "更新活动失败", err)
		return
	}
	utils.BuildSuccessResponse(c, "更新成功")
}

func DeleteActivity(c *gin.Context) {
	aidStr := c.Query("aid")
	if aidStr == "" {
		utils.BuildErrorResponse(c, 400, "aid 为必传参数 请传递aid")
		return
	}
	aid, err := strconv.ParseInt(aidStr, 10, 64)
	if err != nil {
		utils.BuildErrorResponse(c, 400, "aid 参数格式有误")
		return
	}
	if err := service.DeleteActivity(aid); err != nil {
		utils.BuildServerError(c, "删除活动失败", err)
		return
	}
	utils.BuildSuccessResponse(c, "删除成功")
}
