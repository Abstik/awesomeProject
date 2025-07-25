package dao

import "awesomeProject/model"

func InsertDonation(donation model.DonationPO) error {
	return db.Create(&donation).Error
}

func BulkInsertDonations(donations []model.DonationPO) error {
	return db.Create(&donations).Error
}

func GetDonations(year string) ([]model.DonationPO, error) {
	var donations []model.DonationPO
	result := db.Where("YEAR(time) = ?", year).Find(&donations)
	if result.Error != nil {
		return nil, result.Error
	}
	return donations, nil
}

func DeleteDonation(id int) error {
	return db.Delete(&model.DonationPO{}, id).Error
}
