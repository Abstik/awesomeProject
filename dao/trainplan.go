package dao

import (
	"awesomeProject/model"
)

func UpdateTrainPlan(req model.TrainPlan) error {
	return db.Model(&model.TrainPlan{}).Where("id = ?", 1).Update("content", req.Content).Error
}

func GetTrainPlan() (model.TrainPlan, error) {
	var trainPlan model.TrainPlan
	return trainPlan, db.Model(&model.TrainPlan{}).First(&trainPlan).Error
}
