package dao

import (
	"fmt"

	"awesomeProject/model"
)

func InsertMember(mem *model.MemberPO) error {
	result := db.Create(mem)
	return result.Error
}

// GetMemberList retrieves members based on optional filters
func GetMemberList(team *string, isGraduate *int) ([]model.MemberPO, error) {
	var members []model.MemberPO
	query := db.Model(&model.MemberPO{})

	if team != nil {
		query = query.Where("team = ?", *team)
	}
	if isGraduate != nil {
		query = query.Where("is_graduate = ?", *isGraduate)
	}

	if err := query.Find(&members).Error; err != nil {
		return nil, fmt.Errorf("failed to query members: %v", err)
	}

	return members, nil
}

// GetMemberByUsername 根据用户名查询用户
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
	// 使用 Gorm 的 Update 方法明确更新数据
	// todo 为什么userName和uid是联合索引呢 我真是想不透彻啊 id自己不行么 这里如果直接save还会因为联合索引不完全一致而建立新的数据 我的老天鹅啊
	result := db.Model(&model.MemberPO{}).Where("uid = ?", member.UID).Updates(member)
	if result.Error != nil {
		return result.Error
	}

	// 如果未更新任何记录，返回错误
	if result.RowsAffected == 0 {
		return fmt.Errorf("no record found to update")
	}

	return nil
}
