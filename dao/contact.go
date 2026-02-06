package dao

import "awesomeProject/model"

func GetContact() (*model.ContactPO, error) {
	var contact model.ContactPO
	err := db.First(&contact).Error
	if err != nil {
		return nil, err
	}
	return &contact, nil
}

func UpdateContact(contact *model.ContactPO) error {
	return db.Select("contact_info").Save(contact).Error
}
