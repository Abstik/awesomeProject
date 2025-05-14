package model

type ContactPO struct {
	Tid         int    `gorm:"primaryKey;column:tid" json:"-"`
	ContactInfo string `gorm:"column:contactInfo" json:"qqnumber"`
}

func (*ContactPO) TableName() string {
	return "contact"
}
