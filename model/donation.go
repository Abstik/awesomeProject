package model

import "time"

type AddDonationReq struct {
	Name   *string  `json:"name" binding:"required"`               // 捐赠者姓名，必填
	Team   *string  `json:"team" binding:"required"`               // 团队名称，必填
	Money  *float64 `json:"money" binding:"required"`              // 捐款金额，必填
	Time   *string  `json:"time" binding:"required,datetime=2006"` // 捐款时间，必填
	Remark *string  `json:"remark"`                                // 备注信息，可选
}

type AddDonationsReq struct {
	Donations []AddDonationReq `json:"donations" binding:"required,dive"` // 必填，且校验数组中的每个对象
}

// DonationPO 捐款数据库实体
type DonationPO struct {
	ID     int        `gorm:"column:id;primaryKey" json:"id"` // 主键，自增
	Name   *string    `gorm:"column:name" json:"name"`        // 捐赠者姓名
	Team   *string    `gorm:"column:team" json:"team"`        // 团队名称
	Money  *float64   `gorm:"column:money" json:"money"`      // 捐款金额
	Time   *time.Time `gorm:"column:time" json:"time"`        // 捐款时间
	Remark *string    `gorm:"column:remark" json:"remark"`    // 备注信息
}

func (*DonationPO) TableName() string {
	return "donation"
}
