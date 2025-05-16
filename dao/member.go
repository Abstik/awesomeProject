package dao

import (
	"fmt"

	"awesomeProject/model"
)

// 添加成员
func InsertMember(mem *model.MemberPO) error {
	result := db.Create(mem)
	return result.Error
}

// 根据组别和毕业状态批量查询成员
func GetMemberList(team *string, isGraduate, pageSize, pageNum, year *int) ([]model.MemberPO, int64, error) {
	var members []model.MemberPO
	var total int64
	query := db.Model(&model.MemberPO{})
	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	if team != nil {
		query = query.Where("team = ?", *team)
	}
	if isGraduate != nil {
		query = query.Where("is_graduate = ?", *isGraduate)
	}
	if pageSize != nil && pageNum != nil {
		query = query.Limit(*pageSize).Offset((*pageNum - 1) * (*pageSize))
	}
	if year != nil {
		query = query.Where("year = ?", *year)
	}

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

// 根据 UID 查询用户
func GetMemberByUID(uid int) (*model.MemberPO, error) {
	var member model.MemberPO
	result := db.Where("uid = ?", uid).First(&member)
	if result.Error != nil {
		return nil, result.Error
	}
	return &member, nil
}

// 更新用户信息
func UpdateMember(member *model.MemberPO) error {
	// todo 为什么userName和uid是联合索引呢 我真是想不透彻啊 id自己不行么 这里如果直接save还会因为联合索引不完全一致而建立新的数据 我的老天鹅啊
	result := db.Model(&model.MemberPO{}).Where("username = ?", member.Username).Updates(member)
	if result.Error != nil {
		return result.Error
	}

	// 如果未更新任何记录，返回错误
	if result.RowsAffected == 0 {
		return fmt.Errorf("no record found to update")
	}

	return nil
}

func GetYears() ([]int, error) {
	var years []int
	err := db.Model(&model.MemberPO{}).
		Distinct("year").
		Where("is_graduate = ?", 1).
		Pluck("year", &years).Error

	if err != nil {
		return nil, err
	}

	return years, nil
}
