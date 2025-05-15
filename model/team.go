package model

type TeamPO struct {
	Tid       int     `gorm:"column:tid;primaryKey;autoIncrement" json:"tid,omitempty"`
	Name      *string `gorm:"column:name;type:varchar(255);not null" json:"name"`
	BrefInfo  *string `gorm:"column:bref_info;type:text" json:"brefInfo"`
	TrainPlan *string `gorm:"column:train_plan;type:text" json:"trainPlan,omitempty"`
	IsExist   *bool   `gorm:"column:is_exist;not null" json:"isExist"`
	Delay     int64   `gorm:"column:delay" json:"delay,omitempty"`
}

func (*TeamPO) TableName() string {
	return "team"
}

type AddTeamReq struct {
	Name      *string `json:"name,omitempty"`      // 团队名称
	BrefInfo  *string `json:"brefInfo,omitempty"`  // 简介
	TrainPlan *string `json:"trainPlan,omitempty"` // 培训计划
	IsExist   *bool   `json:"isExist,omitempty"`   // 是否存在
	QQ        *string `json:"qq,omitempty"`
}
