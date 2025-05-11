package service

import (
	"time"

	"awesomeProject/dao"
	"awesomeProject/model"
)

func AddDonation(req model.AddDonationReq) error {
	// 解析时间字符串为 time.Time
	donationTime, err := time.Parse("2006-01-02 15:04:05", *req.Time)
	if err != nil {
		return err
	}

	// 构造 DonationPO 对象
	donation := model.DonationPO{
		Name:   req.Name,
		Team:   req.Team,
		Money:  req.Money,
		Time:   &donationTime,
		Remark: req.Remark,
	}

	// 调用 DAO 层插入记录
	return dao.InsertDonation(donation)
}

func AddDonations(req model.AddDonationsReq) error {
	var donations []model.DonationPO

	for _, d := range req.Donations {
		donationTime, err := time.Parse("2006-01-02 15:04:05", *d.Time)
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
