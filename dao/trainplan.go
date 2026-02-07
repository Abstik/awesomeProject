package dao

import "awesomeProject/model"

func GetTrainPlan() (*model.TrainPlanPO, error) {
	var trainPlan model.TrainPlanPO
	err := db.First(&trainPlan).Error
	if err != nil {
		return nil, err
	}
	return &trainPlan, nil
}

func UpdateTrainPlan(trainPlan *model.TrainPlanPO) error {
	return db.Model(&model.TrainPlanPO{}).Where("id = ?", 1).Update("content", trainPlan.Content).Error
}
