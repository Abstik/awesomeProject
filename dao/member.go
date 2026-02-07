package dao

import (
	"fmt"

	"gorm.io/gorm"

	"awesomeProject/model"
)

// 添加成员
func InsertMember(mem *model.MemberPO) error {
	result := db.Create(mem)
	return result.Error
}

// 根据条件批量查询成员
func GetMemberList(team *string, isGraduate, pageSize, pageNum, year *int) ([]model.MemberPO, int64, error) {
	var members []model.MemberPO
	var total int64

	// 构造基础查询条件
	baseQuery := db.Model(&model.MemberPO{})
	if team != nil {
		baseQuery = baseQuery.Where("team = ?", *team)
	}
	if isGraduate != nil {
		baseQuery = baseQuery.Where("is_graduate = ?", *isGraduate)
	}
	if year != nil {
		baseQuery = baseQuery.Where("year = ?", *year)
	}

	// 先 count 总数
	if err := baseQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 使用独立会话避免 Count 的 SELECT 污染后续 Find
	query := baseQuery.Session(&gorm.Session{})
	if pageSize != nil && pageNum != nil {
		query = query.Limit(*pageSize).Offset((*pageNum - 1) * (*pageSize))
	}

	// 获取数据
	if err := query.Find(&members).Error; err != nil {
		return nil, 0, err
	}

	return members, total, nil
}

// 根据用户名查询用户
func GetMemberByUsername(userName string) (*model.MemberPO, error) {
	var member model.MemberPO
	result := db.Where("username = ?", userName).First(&member)
	if result.Error != nil {
		return nil, result.Error
	}
	return &member, nil
}

func GetMemberByName(name string) ([]model.MemberPO, error) {
	var members []model.MemberPO
	result := db.Where("name LIKE ? AND status != 0", "%"+name+"%").Find(&members)
	if result.Error != nil {
		return nil, result.Error
	}
	return members, nil
}

// 更新用户信息
func UpdateMember(member *model.MemberPO) error {
	result := db.Model(&model.MemberPO{}).Where("username = ?", member.Username).Updates(member)
	if result.Error != nil {
		return result.Error
	}

	// 如果未更新任何记录，返回错误
	if result.RowsAffected == 0 {
		return fmt.Errorf("未找到可更新的记录")
	}

	return nil
}

func GetYears() ([]int, error) {
	var years []int
	err := db.Model(&model.MemberPO{}).
		Distinct("year").
		Where("is_graduate = ?", 1).
		Order("year DESC").
		Pluck("year", &years).Error

	if err != nil {
		return nil, err
	}

	return years, nil
}

func DeleteMember(uid int64) error {
	return db.Where("uid = ?", uid).Delete(&model.MemberPO{}).Error
}

func ResetPassword(username, password string) error {
	result := db.Model(&model.MemberPO{}).Where("username = ?", username).Update("password", password)
	if result.Error != nil {
		return result.Error
	}

	// 如果未更新任何记录，返回错误
	if result.RowsAffected == 0 {
		return fmt.Errorf("未找到可更新的记录")
	}

	return nil
}
