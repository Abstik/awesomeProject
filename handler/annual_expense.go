package handler

import (
	"awesomeProject/model"
	"awesomeProject/service"
	"awesomeProject/utils"
	"errors"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"strconv"
	"strings"
)

func expenseError(c *gin.Context, err error) {
	if errors.Is(err, service.ErrInvalidAnnualExpense) {
		utils.BuildErrorResponse(c, 400, strings.TrimPrefix(err.Error(), service.ErrInvalidAnnualExpense.Error()+": "))
		return
	}
	utils.BuildServerError(c, "年度支出操作失败，请重试", err)
}
func GetAnnualExpense(c *gin.Context) {
	year, err := strconv.Atoi(c.Query("year"))
	if err != nil {
		utils.BuildErrorResponse(c, 400, "year 参数格式有误")
		return
	}
	view, err := service.GetAnnualExpense(c.Request.Context(), year)
	if err != nil {
		expenseError(c, err)
		return
	}
	utils.BuildSuccessResponse(c, view)
}
func SaveAnnualExpense(c *gin.Context) {
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, service.MaxAnnualExpenseBytes*6+4096)
	var req model.AnnualExpenseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		var tooLarge *http.MaxBytesError
		if errors.As(err, &tooLarge) || errors.Is(err, io.ErrUnexpectedEOF) {
			utils.BuildErrorResponse(c, 400, "请求内容过大或格式有误")
			return
		}
		utils.BuildErrorResponse(c, 400, "请填写有效的年份和支出内容")
		return
	}
	view, err := service.SaveAnnualExpense(c.Request.Context(), req)
	if err != nil {
		expenseError(c, err)
		return
	}
	utils.BuildSuccessResponse(c, view)
}
func DeleteAnnualExpense(c *gin.Context) {
	year, err := strconv.Atoi(c.Query("year"))
	if err != nil {
		utils.BuildErrorResponse(c, 400, "year 参数格式有误")
		return
	}
	if err := service.DeleteAnnualExpense(c.Request.Context(), year); err != nil {
		expenseError(c, err)
		return
	}
	utils.BuildSuccessResponse(c, nil)
}

func ListAnnualExpenseYears(c *gin.Context) {
	years, err := service.ListAnnualExpenseYears(c.Request.Context())
	if err != nil {
		expenseError(c, err)
		return
	}
	utils.BuildSuccessResponse(c, years)
}
