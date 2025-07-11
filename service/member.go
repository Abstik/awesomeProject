package service

import (
	"errors"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"awesomeProject/dao"
	"awesomeProject/model"
	"awesomeProject/utils"
)

// 用户注册
func Register(mem *model.MemberRequest) error {
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

	password := utils.EncryptPassword(*mem.Username + "123")

	status := 1
	// 插入数据库
	memPO := &model.MemberPO{
		Username:   mem.Username,
		Name:       mem.Name,
		Year:       mem.Year,
		Team:       mem.Team,
		Password:   &password,
		IsGraduate: &isGraduate,
		Status:     &status,
	}
	err := dao.InsertMember(memPO)
	if err != nil {
		return err
	}
	return nil
}

// 用户登录
func Login(req *model.MemberRequest) (gin.H, error) {
	// 根据用户名查询用户信息
	user, err := dao.GetMemberByUsername(*req.Username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("用户名不存在")
		}
		return nil, err
	}

	// 密码校验
	if *user.Password != utils.EncryptPassword(*req.Password) {
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

// 批量查询成员
func GetMemberList(team *string, isGraduate, pageSize, pageNum, year *int) ([]model.MemberPO, int64, error) {
	res, total, err := dao.GetMemberList(team, isGraduate, pageSize, pageNum, year)
	if err != nil {
		return nil, 0, err
	}
	return res, total, nil
}

// 根据用户名获取用户信息
func GetMemberByUsername(userName string) (*model.MemberPO, error) {
	member, err := dao.GetMemberByUsername(userName)
	if err != nil {
		return nil, err
	}
	return member, nil
}

func GetMemberByName(name string) ([]model.MemberPO, error) {
	members, err := dao.GetMemberByName(name)
	return members, err
}

func UpdateMember(req model.UpdateMemberRequest, statusInt int) error {
	// 查询用户是否存在
	member, err := dao.GetMemberByUsername(*req.Username)
	if err != nil {
		return errors.New("user not found")
	}

	// 判断请求用户是否和修改用户是否一致
	if *member.Username != *req.Username {
		if statusInt != 0 {
			return errors.New("无权限")
		}
	}

	// 更新用户字段（仅更新非空字段）
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
		member.Portrait = req.Portrait
	}
	if req.MienImg != nil {
		member.MienImg = req.MienImg
	}
	if req.Company != nil {
		member.Company = req.Company
	}
	if req.GraduateImg != nil {
		member.GraduateImg = req.GraduateImg
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

	// 调用 DAO 层保存更新
	if err := dao.UpdateMember(member); err != nil {
		return errors.New("failed to update user in database")
	}

	return nil
}

func GetYears() ([]int, error) {
	years, err := dao.GetYears()
	return years, err
}
