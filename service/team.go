package service

import (
	"awesomeProject/dao"
	"awesomeProject/model"
)

func AddTeam(req model.AddTeamReq) error {
	team := model.TeamPO{
		Name:      req.Name,
		BrefInfo:  req.BrefInfo,
		TrainPlan: req.TrainPlan,
		IsExist:   req.IsExist,
	}

	// 调用 DAO 层的方法
	return dao.InsertTeam(team)
}

func GetTeams(name string, isExist *bool, isUser bool) ([]model.TeamPO, error) {
	// 调用 DAO 层查询
	return dao.QueryTeams(name, isExist, isUser)
}
