package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"awesomeProject/model"
	"awesomeProject/service"
	"awesomeProject/utils"
)

func AddDonations(c *gin.Context) {
	var req model.AddDonationsReq
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BuildErrorResponse(c, 400, "参数格式有误")
		return
	}

	// 调用 service 层处理
	if err := service.AddDonations(req); err != nil {
		utils.BuildServerError(c, "添加捐款失败", err)
		return
	}

	utils.BuildSuccessResponse(c, "添加成功")
}

// 根据year查询捐款信息
func GetDonations(c *gin.Context) {
	year := c.Query("year")
	if year == "" {
		utils.BuildErrorResponse(c, 400, "year 为必传参数")
		return
	}
	donations, totalCount, err := service.GetDonations(year)
	if err != nil {
		utils.BuildServerError(c, "查询捐款失败", err)
		return
	}
	utils.BuildSuccessResponse(c, gin.H{
		"donations":  donations,
		"totalCount": totalCount,
	})
}

func DeleteDonation(c *gin.Context) {
	// 从查询参数中获取 id
	idStr := c.Query("id")
	if idStr == "" {
		utils.BuildErrorResponse(c, 400, "id 为必传参数")
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.BuildErrorResponse(c, 400, "id 参数格式有误")
		return
	}
	if err := service.DeleteDonation(id); err != nil {
		utils.BuildServerError(c, "删除捐款失败", err)
		return
	}
	utils.BuildSuccessResponse(c, "删除成功")
}
