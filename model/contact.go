package model

type ContactPO struct {
	ID          int    `gorm:"primaryKey;column:id" json:"-"`
	ContactInfo string `gorm:"column:contact_info" json:"qqnumber"`
}

func (*ContactPO) TableName() string {
	return "contact"
}
