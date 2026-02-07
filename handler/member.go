package handler

import (
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"

	"awesomeProject/model"
	"awesomeProject/service"
	"awesomeProject/utils"
)

// 添加成员
func Register(c *gin.Context) {
	var req *model.RegisterReq
	err := c.ShouldBindJSON(&req)
	if err != nil {
		utils.BuildErrorResponse(c, 400, "参数格式有误")
		return
	}

	err = service.Register(req)
	if err != nil {
		if err.Error() == "用户名已存在" {
			utils.BuildErrorResponse(c, 400, "用户名已存在")
			return
		}
		utils.BuildServerError(c, "注册失败", err)
		return
	}
	utils.BuildSuccessResponse(c, "注册成功")
}

// 用户登录
func Login(c *gin.Context) {
	var req *model.LoginReq
	err := c.ShouldBindJSON(&req)
	if err != nil {
		utils.BuildErrorResponse(c, 400, "参数格式有误")
		return
	}

	// 校验验证码
	if ok := VerifyCaptcha(*req.CaptchaID, *req.CaptchaData); !ok {
		utils.BuildErrorResponse(c, 400, "验证码有误")
		return
	}

	data, err := service.Login(req)
	if err != nil {
		// 如果是用户名不存在或密码错误，记录警告日志用于安全审计
		if err.Error() == "用户名不存在" || err.Error() == "密码错误" {
			zap.L().Warn("登录失败",
				zap.String("username", *req.Username),
				zap.String("ip", c.ClientIP()),
				zap.String("reason", err.Error()),
			)
			utils.BuildErrorResponse(c, 400, "用户名或密码错误")
			return
		}
		utils.BuildServerError(c, "登录失败", err)
		return
	}
	utils.BuildSuccessResponse(c, data)
}

// 获取用户列表
func GetMemberList(c *gin.Context) {
	teamStr := c.Query("team")
	isGraduateStr := c.Query("isGraduate")
	pageSizeStr := c.Query("pageSize")
	pageNumStr := c.Query("pageNum")
	yearStr := c.Query("year")

	var pageSize, pageNum, isGraduate, year *int
	var team *string
	if isGraduateStr != "" {
		var tmp int
		if isGraduateStr == "true" {
			tmp = 1
			isGraduate = &tmp
		} else if isGraduateStr == "false" {
			tmp = 0
			isGraduate = &tmp
		} else {
			utils.BuildErrorResponse(c, 400, "参数格式有误")
			return
		}
	}
	if pageSizeStr != "" {
		pageSizeInt, err := strconv.Atoi(pageSizeStr)
		if err != nil {
			utils.BuildErrorResponse(c, 400, "pageSize 参数格式有误")
			return
		}
		pageSize = &pageSizeInt
	}
	if pageNumStr != "" {
		pageNumInt, err := strconv.Atoi(pageNumStr)
		if err != nil {
			utils.BuildErrorResponse(c, 400, "pageNum 参数格式有误")
			return
		}
		pageNum = &pageNumInt
	}
	if teamStr != "" {
		team = &teamStr
	}
	if yearStr != "" {
		yearInt, err := strconv.Atoi(yearStr)
		if err != nil {
			utils.BuildErrorResponse(c, 400, "year 参数格式有误")
			return
		}
		year = &yearInt
	}

	res, total, err := service.GetMemberList(team, isGraduate, pageSize, pageNum, year)
	if err != nil {
		utils.BuildServerError(c, "查询成员列表失败", err)
		return
	}

	utils.BuildSuccessResponse(c, gin.H{
		"data":  res,
		"total": total,
	})
}

// 修改用户信息
func ChangeMemberInfo(c *gin.Context) {
	// 定义请求体结构
	var req model.UpdateMemberRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BuildErrorResponse(c, 400, "参数格式有误")
		return
	}

	status, ok := c.Get("status")
	if !ok {
		utils.BuildServerError(c, "用户状态获取失败", errors.New("status not found in context"))
		return
	}
	statusInt, ok := status.(int)
	if !ok {
		utils.BuildServerError(c, "用户状态类型异常", errors.New("status type assertion failed"))
		return
	}

	// 从 JWT 上下文获取当前登录用户名
	authUsername, _ := c.Get("username")
	authUsernameStr, _ := authUsername.(string)

	// 调用 Service 层更新用户信息
	var authPtr *string
	if statusInt != 0 {
		// 普通用户传指针，管理员为nil
		authPtr = &authUsernameStr
	}
	err := service.UpdateMember(req, authPtr)
	if err != nil {
		if err.Error() == "无权限修改他人信息" || err.Error() == "无法修改管理员账号" {
			utils.BuildErrorResponse(c, 403, err.Error())
		} else {
			utils.BuildServerError(c, "修改失败", err)
		}
		return
	}

	// 返回成功响应
	utils.BuildSuccessResponse(c, "修改成功")
}

// 根据用户名获取用户信息
func GetMemberByUserName(c *gin.Context) {
	userName := c.Query("username")
	if userName == "" {
		utils.BuildErrorResponse(c, 400, "username 为必传参数")
		return
	}

	// 调用 Service 层获取用户数据
	member, err := service.GetMemberByUsername(userName)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		utils.BuildServerError(c, "查询用户失败", err)
		return
	}

	// 返回成功响应
	utils.BuildSuccessResponse(c, member)
}

func GetMemberByName(c *gin.Context) {
	name := c.Query("name")
	if name == "" {
		utils.BuildErrorResponse(c, 400, "name 为必传参数")
		return
	}

	// 调用 Service 层获取用户数据
	members, err := service.GetMemberByName(name)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		utils.BuildServerError(c, "查询成员失败", err)
		return
	}

	// 成功响应
	utils.BuildSuccessResponse(c, members)
}

func GetYears(c *gin.Context) {
	years, err := service.GetYears()
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		utils.BuildServerError(c, "查询年份失败", err)
		return
	}
	utils.BuildSuccessResponse(c, years)
}

func DeleteMember(c *gin.Context) {
	uidStr := c.Query("uid")
	if uidStr == "" {
		utils.BuildErrorResponse(c, 400, "uid 为必传参数")
		return
	}

	uid, err := strconv.ParseInt(uidStr, 10, 64)
	if err != nil {
		utils.BuildErrorResponse(c, 400, "uid 参数格式有误")
		return
	}
	if err := service.DeleteMember(uid); err != nil {
		utils.BuildServerError(c, "删除成员失败", err)
		return
	}
	utils.BuildSuccessResponse(c, "删除成功")
}

func ResetPassword(c *gin.Context) {
	username := c.PostForm("username")
	if username == "" {
		utils.BuildErrorResponse(c, 400, "username 为必传参数")
		return
	}

	err := service.ResetPassword(username)
	if err != nil {
		if err.Error() == "用户不存在" || err.Error() == "无法修改管理员账号" {
			utils.BuildErrorResponse(c, 400, err.Error())
			return
		}
		utils.BuildServerError(c, "重置密码失败", err)
		return
	}
	utils.BuildSuccessResponse(c, gin.H{
		"message": "重置成功，密码已重置为 用户名+123",
	})
}
