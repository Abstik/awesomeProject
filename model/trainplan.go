package model

type TrainPlan struct {
	ID      uint    `gorm:"primarykey" json:"-"`
	Content *string `gorm:"content" json:"content"`
}

func (*TrainPlan) TableName() string {
	return "train_plan"
}
