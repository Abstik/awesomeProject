package service

import (
	"awesomeProject/dao"
	"awesomeProject/model"
	"context"
	"errors"
	"fmt"
	"strings"
)

const MaxAnnualExpenseBytes = 102400

var ErrInvalidAnnualExpense = errors.New("invalid annual expense")

func ValidateAnnualExpenseYear(year int) error {
	if year < 2015 || year > 9999 {
		return fmt.Errorf("%w: 请输入 2015 至 9999 之间的四位年份", ErrInvalidAnnualExpense)
	}
	return nil
}
func GetAnnualExpense(ctx context.Context, year int) (model.AnnualExpenseVO, error) {
	view := model.AnnualExpenseVO{Year: year}
	if err := ValidateAnnualExpenseYear(year); err != nil {
		return view, err
	}
	report, err := dao.FindAnnualExpense(ctx, year)
	if err != nil {
		return view, err
	}
	if report != nil {
		view.Content = report.Content
		view.Exists = true
		view.UpdatedAt = &report.UpdatedAt
	}
	return view, nil
}
func SaveAnnualExpense(ctx context.Context, req model.AnnualExpenseRequest) (model.AnnualExpenseVO, error) {
	if err := ValidateAnnualExpenseYear(req.Year); err != nil {
		return model.AnnualExpenseVO{}, err
	}
	if len(req.Content) > MaxAnnualExpenseBytes {
		return model.AnnualExpenseVO{}, fmt.Errorf("%w: 支出信息不能超过 100 KB", ErrInvalidAnnualExpense)
	}
	content := sanitizeAnnualExpense(req.Content)
	if strings.TrimSpace(annualExpenseText(content)) == "" {
		return model.AnnualExpenseVO{}, fmt.Errorf("%w: 请填写年度支出明细", ErrInvalidAnnualExpense)
	}
	report, err := dao.UpsertAnnualExpense(ctx, req.Year, content)
	if err != nil {
		return model.AnnualExpenseVO{}, err
	}
	return model.AnnualExpenseVO{Year: report.Year, Content: report.Content, Exists: true, UpdatedAt: &report.UpdatedAt}, nil
}
func DeleteAnnualExpense(ctx context.Context, year int) error {
	if err := ValidateAnnualExpenseYear(year); err != nil {
		return err
	}
	return dao.DeleteAnnualExpense(ctx, year)
}

func ListAnnualExpenseYears(ctx context.Context) ([]int, error) {
	return dao.ListAnnualExpenseYears(ctx)
}
