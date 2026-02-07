package service

import (
	"time"

	"awesomeProject/dao"
	"awesomeProject/model"
)

// GetDonations 查询捐款并计算总金额
func GetDonations(year string) ([]model.DonationPO, float64, error) {
	donations, err := dao.GetDonations(year)
	if err != nil {
		return nil, 0, err
	}
	var totalCount float64
	for _, d := range donations {
		if d.Money != nil {
			totalCount += *d.Money
		}
	}
	return donations, totalCount, nil
}

func AddDonations(req model.AddDonationsReq) error {
	var donations []model.DonationPO

	for _, d := range req.Donations {
		donationTime, err := time.Parse("2006", *d.Time)
		if err != nil {
			return err
		}

		donations = append(donations, model.DonationPO{
			Name:   d.Name,
			Team:   d.Team,
			Money:  d.Money,
			Time:   &donationTime,
			Remark: d.Remark,
		})
	}

	// 调用 DAO 层批量插入
	return dao.BulkInsertDonations(donations)
}

// DeleteDonation 删除捐款
func DeleteDonation(id int) error {
	return dao.DeleteDonation(id)
}
