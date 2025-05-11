package model

import "time"

type ActivityPO struct {
	AID     int64      `gorm:"primaryKey;column:aid" json:"aid"`
	Title   *string    `gorm:"column:" json:"title"`
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

type ActivityListReq struct {
	PageSize int64 `json:"page_size,omitempty"`
	PageNum  int64 `json:"page_num,omitempty"`
}

type AddActivityReq struct {
	Title   *string `gorm:"column:" json:"title"`
	Img     *string `gorm:"column:img" json:"img"`
	Summary *string `gorm:"column:summary" json:"summary"`
	Content *string `gorm:"column:content" json:"content"`
}
