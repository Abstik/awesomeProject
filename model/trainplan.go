package model

type TrainPlanPO struct {
	ID      int    `gorm:"primaryKey;column:id" json:"-"`
	Content string `gorm:"column:content" json:"content" binding:"required"`
}

func (*TrainPlanPO) TableName() string {
	return "train_plan"
}
