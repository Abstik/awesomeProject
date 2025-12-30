package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"awesomeProject/dao"
	"awesomeProject/model"
	"awesomeProject/service"
	"awesomeProject/utils"
)

func AddDonations(c *gin.Context) {
	var req model.AddDonationsReq
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BuildErrorResponse(c, 500, "AddDonations parse failed err is: "+err.Error())
		return
	}

	// 调用 service 层处理
	if err := service.AddDonations(req); err != nil {
		utils.BuildErrorResponse(c, 500, "AddDonations insert failed err is: "+err.Error())
		return
	}

	utils.BuildSuccessResponse(c, "添加成功")
}

// 根据year查询捐款信息
func GetDonations(c *gin.Context) {
	year := c.Query("year")
	if year == "" {
		utils.BuildErrorResponse(c, 400, "GetDonations year is empty")
		return
	}
	donations, err := dao.GetDonations(year)
	if err != nil {
		utils.BuildErrorResponse(c, 500, "GetDonations err is: "+err.Error())
		return
	}
	// 统计捐款总金额
	var totalCount float64
	for _, donation := range donations {
		totalCount += *donation.Money
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
		utils.BuildErrorResponse(c, 400, "DeleteDonation id is empty")
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.BuildErrorResponse(c, 400, "DeleteDonation id is not a number")
		return
	}
	if err := dao.DeleteDonation(id); err != nil {
		utils.BuildErrorResponse(c, 500, "DeleteDonation err is: "+err.Error())
		return
	}
	utils.BuildSuccessResponse(c, "删除成功")
}
