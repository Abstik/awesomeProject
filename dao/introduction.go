package dao

import "awesomeProject/model"

func GetIntroduction() (*model.IntroductionPO, error) {
	var introduction model.IntroductionPO
	err := db.First(&introduction).Error
	if err != nil {
		return nil, err
	}
	return &introduction, nil
}

func UpdateIntroduction(introduction *model.IntroductionPO) error {
	return db.Model(&model.IntroductionPO{}).Where("id = ?", 1).Update("content", introduction.Content).Error
}
