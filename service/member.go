package service

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"awesomeProject/dao"
	"awesomeProject/model"
	"awesomeProject/utils"
)

// 用户注册
func Register(mem *model.RegisterReq) error {
	// 检查用户名是否已存在
	if _, err := dao.GetMemberByUsername(*mem.Username); err == nil {
		// 如果查询到则返回错误
		return errors.New("用户名已存在")
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		// 如果出现错误，且错误不是未查询到
		return err
	}

	var isGraduate int
	if time.Now().Month() > 6 && time.Now().Year()-*mem.Year >= 4 {
		isGraduate = 1
	} else {
		isGraduate = 0
	}

	// 默认密码为 用户名+123，使用 bcrypt 哈希
	password, err := utils.HashPassword(*mem.Username + "123")
	if err != nil {
		return fmt.Errorf("密码加密失败: %w", err)
	}

	defaultPortrait := os.Getenv("DEFAULT_PORTRAIT")
	defaultMienImg := os.Getenv("DEFAULT_MIEN_IMG")
	defaultGraduateImg := os.Getenv("DEFAULT_GRADUATE_IMG")

	status := 1
	// 插入数据库
	memPO := &model.MemberPO{
		Username:    mem.Username,
		Name:        mem.Name,
		Year:        mem.Year,
		Team:        mem.Team,
		Password:    &password,
		IsGraduate:  &isGraduate,
		Status:      &status,
		Portrait:    &defaultPortrait,
		MienImg:     &defaultMienImg,
		GraduateImg: &defaultGraduateImg,
	}
	if err := dao.InsertMember(memPO); err != nil {
		return err
	}
	return nil
}

// 删除成员
func DeleteMember(uid int64) error {
	return dao.DeleteMember(uid)
}

// 用户登录
func Login(req *model.LoginReq) (gin.H, error) {
	// 根据用户名查询用户信息
	user, err := dao.GetMemberByUsername(*req.Username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("用户名不存在")
		}
		return nil, err
	}

	// 密码校验
	if !utils.CheckPassword(*req.Password, *user.Password) {
		return nil, errors.New("密码错误")
	}

	// 生成token
	token, err := utils.GenToken(user.UID, *user.Username, *user.Status)
	response := gin.H{
		"token":    token,
		"status":   strconv.Itoa(*user.Status),
		"username": user.Username,
	}
	return response, err
}

// fillMemberURLs 填充成员图片的完整 URL
func fillMemberURLs(m *model.MemberPO) {
	m.Portrait = utils.OldFullURL(m.Portrait)
	m.MienImg = utils.OldFullURL(m.MienImg)
	m.GraduateImg = utils.OldFullURL(m.GraduateImg)
}

// 批量查询成员
func GetMemberList(team *string, isGraduate, pageSize, pageNum, year *int) ([]model.MemberPO, int64, error) {
	res, total, err := dao.GetMemberList(team, isGraduate, pageSize, pageNum, year)
	if err != nil {
		return nil, 0, err
	}

	for i := range res {
		fillMemberURLs(&res[i])
	}
	return res, total, nil
}

// 根据用户名获取用户信息
func GetMemberByUsername(userName string) (*model.MemberPO, error) {
	member, err := dao.GetMemberByUsername(userName)
	if err != nil {
		return nil, err
	}
	fillMemberURLs(member)
	return member, nil
}

func GetMemberByName(name string) ([]model.MemberPO, error) {
	res, err := dao.GetMemberByName(name)
	if err != nil {
		return nil, err
	}
	for i := range res {
		fillMemberURLs(&res[i])
	}
	return res, nil
}

// applyMemberUpdates 将请求中的非空字段更新到成员对象（公共逻辑）
func applyMemberUpdates(member *model.MemberPO, req model.UpdateMemberRequest) {
	if req.Name != nil {
		member.Name = req.Name
	}
	if req.Tel != nil {
		member.Tel = req.Tel
	}
	if req.Gender != nil {
		member.Gender = req.Gender
	}
	if req.ClassGrade != nil {
		member.ClassGrade = req.ClassGrade
	}
	if req.Team != nil {
		member.Team = req.Team
	}
	if req.Portrait != nil {
		member.Portrait = utils.ParseURL(req.Portrait)
	}
	if req.MienImg != nil {
		member.MienImg = utils.ParseURL(req.MienImg)
	}
	if req.Company != nil {
		member.Company = req.Company
	}
	if req.GraduateImg != nil {
		member.GraduateImg = utils.ParseURL(req.GraduateImg)
	}
	if req.IsGraduate != nil {
		member.IsGraduate = req.IsGraduate
	}
	if req.Signature != nil {
		member.Signature = req.Signature
	}
	if req.Year != nil {
		member.Year = req.Year
	}
}

// UpdateMember 修改用户信息，authUsername 非空时校验权限（普通用户场景）
func UpdateMember(req model.UpdateMemberRequest, authUsername *string) error {
	// 普通用户只能修改自己的信息
	if authUsername != nil && *req.Username != *authUsername {
		return errors.New("无权限修改他人信息")
	}

	// 查询用户是否存在
	member, err := dao.GetMemberByUsername(*req.Username)
	if err != nil {
		return errors.New("用户不存在")
	}

	if *member.Status == 0 {
		return errors.New("无法修改管理员账号")
	}

	applyMemberUpdates(member, req)

	// 调用 DAO 层保存更新
	if err := dao.UpdateMember(member); err != nil {
		return fmt.Errorf("更新用户信息失败: %w", err)
	}

	return nil
}

// ResetPassword 重置用户密码为 用户名+"123"
func ResetPassword(username string) error {
	user, err := dao.GetMemberByUsername(username)
	if err != nil {
		return errors.New("用户不存在")
	}
	if user == nil || user.Status == nil || user.Username == nil {
		return errors.New("用户数据异常")
	}
	if *user.Status == 0 {
		return errors.New("无法修改管理员账号")
	}

	hashedPassword, err := utils.HashPassword(username + "123")
	if err != nil {
		return fmt.Errorf("密码加密失败: %w", err)
	}

	return dao.ResetPassword(username, hashedPassword)
}

func GetYears() ([]int, error) {
	years, err := dao.GetYears()
	return years, err
}
