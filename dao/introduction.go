package dao

import (
	"awesomeProject/model"
)

func UpdateIntroduction(introduction *model.IntroductionPO) error {
	return db.Model(introduction).
		Where("id = ?", introduction.Id).
		Update("content", introduction.Content).Error
}

func GetIntroduction() (*model.IntroductionPO, error) {
	var Introduction model.IntroductionPO
	err := db.First(&Introduction).Error
	if err != nil {
		return nil, err
	}
	return &Introduction, nil
}
