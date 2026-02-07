package model

import "time"

type ActivityPO struct {
	AID     int64      `gorm:"primaryKey;column:aid" json:"aid"`
	Title   *string    `gorm:"column:title" json:"title"`
	Img     *string    `gorm:"column:img" json:"img"`
	Summary *string    `gorm:"column:summary" json:"summary"`
	Content *string    `gorm:"column:content" json:"content"`
	Time    *time.Time `gorm:"column:time" json:"time"`
	View    *int64     `gorm:"column:view" json:"view"`
	Status  *int64     `gorm:"column:status" json:"status"`
}

func (*ActivityPO) TableName() string {
	return "activity"
}

type ActivityReq struct {
	AID     int64   `json:"aid"`
	Title   *string `json:"title"`
	Img     *string `json:"img"`
	Summary *string `json:"summary"`
	Content *string `json:"content"`
}
