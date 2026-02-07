package model

type TeamPO struct {
	Tid       int     `gorm:"column:tid;primaryKey" json:"tid"`
	Name      *string `gorm:"column:name" json:"name"`
	BrefInfo  *string `gorm:"column:bref_info" json:"brefInfo"`
	TrainPlan *string `gorm:"column:train_plan" json:"trainPlan"`
	IsExist   *bool   `gorm:"column:is_exist" json:"isExist"`
}

func (*TeamPO) TableName() string {
	return "team"
}

type TeamVO struct {
	TeamPO
	Delay int `json:"delay"`
}

type AddTeamReq struct {
	Tid       int     `json:"tid"`       // 团队ID
	Name      *string `json:"name"`      // 团队名称
	BrefInfo  *string `json:"brefInfo"`  // 简介
	TrainPlan *string `json:"trainPlan"` // 培训计划
	IsExist   *bool   `json:"isExist"`   // 是否存在
}
