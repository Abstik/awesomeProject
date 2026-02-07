package model

type IntroductionPO struct {
	ID      int    `gorm:"primaryKey;column:id" json:"-"`
	Content string `gorm:"column:content" json:"introduction"`
}

func (*IntroductionPO) TableName() string {
	return "introduction"
}
