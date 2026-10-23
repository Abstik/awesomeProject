package dao

import (
	"awesomeProject/model"
	"context"
	"errors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"time"
)

func MigrateAnnualExpenses() error { return db.AutoMigrate(&model.AnnualExpensePO{}) }
func FindAnnualExpense(ctx context.Context, year int) (*model.AnnualExpensePO, error) {
	var report model.AnnualExpensePO
	err := db.WithContext(ctx).Where("year = ?", year).First(&report).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &report, err
}
func UpsertAnnualExpense(ctx context.Context, year int, content string) (*model.AnnualExpensePO, error) {
	report := model.AnnualExpensePO{Year: year, Content: content, UpdatedAt: time.Now().UTC()}
	err := db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "year"}},
		DoUpdates: clause.AssignmentColumns([]string{"content", "updated_at"}),
	}).Create(&report).Error
	return &report, err
}
func DeleteAnnualExpense(ctx context.Context, year int) error {
	return db.WithContext(ctx).Where("year = ?", year).Delete(&model.AnnualExpensePO{}).Error
}

func ListAnnualExpenseYears(ctx context.Context) ([]int, error) {
	years := make([]int, 0)
	err := db.WithContext(ctx).Model(&model.AnnualExpensePO{}).Order("year DESC").Pluck("year", &years).Error
	return years, err
}
