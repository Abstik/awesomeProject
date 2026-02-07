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

func GetTeams(name string, isExist *bool, isUser bool) ([]model.TeamVO, error) {
	teamPOs, err := dao.QueryTeams(name, isExist, isUser)
	if err != nil {
		return nil, err
	}

	teamVOs := make([]model.TeamVO, len(teamPOs))
	delay := 0
	for i := range teamPOs {
		teamVOs[i].TeamPO = teamPOs[i]
		if teamPOs[i].IsExist != nil && *teamPOs[i].IsExist {
			teamVOs[i].Delay = delay
			delay += 100
		}
	}

	return teamVOs, nil
}

func UpdateTeam(req model.AddTeamReq) error {
	team := &model.TeamPO{
		Tid:       req.Tid,
		Name:      req.Name,
		BrefInfo:  req.BrefInfo,
		TrainPlan: req.TrainPlan,
		IsExist:   req.IsExist,
	}
	return dao.UpdateTeam(team)
}
