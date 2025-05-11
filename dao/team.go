package dao

import "awesomeProject/model"

func InsertTeam(team model.TeamPO) error {
	return db.Create(&team).Error
}

func QueryTeams(name string, isExist *bool) ([]model.TeamPO, error) {
	var teams []model.TeamPO
	query := db.Select("tid, name, bref_info, is_exist") // 不查询 train_plan

	if name != "" {
		query = query.Where("name LIKE ?", "%"+name+"%")
	}
	if isExist != nil {
		query = query.Where("is_exist = ?", *isExist)
	}

	if err := query.Find(&teams).Error; err != nil {
		return nil, err
	}

	return teams, nil
}
