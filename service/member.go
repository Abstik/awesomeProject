package service

import (
	"errors"

	"awesomeProject/dao"
	"awesomeProject/model"
)

func GetMemberList(team *string, isGraduate *int) ([]model.MemberPO, error) {
	res, err := dao.GetMemberList(team, isGraduate)
	// 按照team排序一下

	if err != nil {
		return nil, err
	}
	return res, nil
}

// 根据用户名获取用户信息
func GetMemberByUsername(userName string) (*model.MemberPO, error) {
	// 调用 DAO 层获取用户数据
	member, err := dao.GetMemberByUsername(userName)
	if err != nil {
		return nil, err
	}
	return member, nil
}

func UpdateMember(req model.UpdateMemberRequest) error {
	// 查询用户是否存在
	member, err := dao.GetMemberByUID(req.UID)
	if err != nil {
		return errors.New("user not found")
	}

	// 更新用户字段（仅更新非空字段）
	if req.Portrait != nil {
		member.Portrait = req.Portrait
	}
	if req.ClassGrade != nil {
		member.ClassGrade = req.ClassGrade
	}
	if req.Phone != nil {
		member.Tel = req.Phone
	}
	if req.Username != nil {
		member.Username = req.Username
	}
	if req.Name != nil {
		member.Name = req.Name
	}
	if req.Team != nil {
		member.Team = req.Team
	}
	if req.MienImg != nil {
		member.MienImg = req.MienImg
	}
	if req.Signature != nil {
		member.Signature = req.Signature
	}

	// 调用 DAO 层保存更新
	if err := dao.UpdateMember(member); err != nil {
		return errors.New("failed to update user in database")
	}

	return nil
}
