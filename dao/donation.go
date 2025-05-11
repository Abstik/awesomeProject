package dao

import "awesomeProject/model"

// InsertDonation adds a new donation record to the database
func InsertDonation(donation model.DonationPO) error {
	return db.Create(&donation).Error
}

func BulkInsertDonations(donations []model.DonationPO) error {
	return db.Create(&donations).Error
}
