package model

type TeamPO struct {
	Tid       int     `gorm:"column:tid;primaryKey;autoIncrement" json:"tid"`
	Name      *string `gorm:"column:name;type:varchar(255);not null" json:"name"`
	BrefInfo  *string `gorm:"column:bref_info;type:text" json:"bref_info"`
	TrainPlan *string `gorm:"column:train_plan;type:text" json:"train_plan"`
	IsExist   *bool   `gorm:"column:is_exist;not null" json:"is_exist"`
}

func (*TeamPO) TableName() string {
	return "team"
}

type AddTeamReq struct {
	Name      *string `json:"name,omitempty"`       // 团队名称
	BrefInfo  *string `json:"bref_info,omitempty"`  // 简介
	TrainPlan *string `json:"train_plan,omitempty"` // 培训计划
	IsExist   *bool   `json:"is_exist,omitempty"`   // 是否存在
}
