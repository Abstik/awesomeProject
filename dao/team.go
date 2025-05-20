package dao

import (
	"awesomeProject/model"
)

func InsertTeam(team model.TeamPO) error {
	return db.Create(&team).Error
}

func QueryTeams(name string, isExist *bool, isUser bool) ([]model.TeamPO, error) {
	var teams []model.TeamPO

	// 游客只查询部分字段，且不加任何查询条件
	if !isUser {
		if err := db.
			Select("tid, name, bref_info, is_exist").
			Find(&teams).Error; err != nil {
			return nil, err
		}
		return teams, nil
	}

	// 用户或管理员：添加查询条件和完整字段
	query := db.Select("tid, name, bref_info, train_plan, is_exist")

	if name != "" {
		query = query.Where("name LIKE ?", name)
	}
	if isExist != nil {
		query = query.Where("is_exist = ?", *isExist)
	}

	if err := query.Find(&teams).Error; err != nil {
		return nil, err
	}

	return teams, nil
}

func UpdateTeam(req model.AddTeamReq) error {
	// 更新非零值字段
	return db.Model(&model.TeamPO{}).Where("tid = ?", req.Tid).Updates(model.TeamPO{
		Name:      req.Name,
		BrefInfo:  req.BrefInfo,
		TrainPlan: req.TrainPlan,
		IsExist:   req.IsExist,
	}).Error
}
