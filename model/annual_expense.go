package model

import "time"

type AnnualExpensePO struct {
	Year      int       `gorm:"primaryKey;autoIncrement:false" json:"year"`
	Content   string    `gorm:"type:mediumtext;not null" json:"content"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func (AnnualExpensePO) TableName() string { return "annual_expense" }

type AnnualExpenseRequest struct {
	Year    int    `json:"year" binding:"required"`
	Content string `json:"content" binding:"required"`
}
type AnnualExpenseVO struct {
	Year      int        `json:"year"`
	Content   string     `json:"content"`
	Exists    bool       `json:"exists"`
	UpdatedAt *time.Time `json:"updatedAt"`
}
