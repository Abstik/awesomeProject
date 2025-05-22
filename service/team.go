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
	teamPOs, err := dao.QueryTeams(name, isExist, isUser)
	if err != nil {
		return nil, err
	}

	delay := 0
	for i := range teamPOs {
		if *teamPOs[i].IsExist {
			teamPOs[i].Delay = delay
			delay += 100
		}
	}

	return teamPOs, nil
}

func UpdateTeam(req model.AddTeamReq) error {
	return dao.UpdateTeam(req)
}
