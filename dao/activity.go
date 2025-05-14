package dao

import "awesomeProject/model"

func GetActivityListByPage(pageSize, pageNum int) ([]*model.ActivityPO, int64, error) {
	var res []*model.ActivityPO
	offset := (pageSize - 1) * pageNum
	var total int64
	result := db.Model(&model.ActivityPO{}).Count(&total)
	if result.Error != nil {
		return nil, 0, result.Error
	}

	dbRes := db.Model(&model.ActivityPO{}).Order("time DESC").Offset(offset).Limit(pageNum).Find(&res)
	if dbRes.Error != nil {
		return nil, 0, dbRes.Error
	}

	return res, total, nil
}

func GetActivityList() ([]*model.ActivityPO, error) {
	var res []*model.ActivityPO
	dbRes := db.Model(&model.ActivityPO{}).Order("time DESC").Find(&res)
	if dbRes.Error != nil {
		return nil, dbRes.Error
	}
	return res, nil
}

func GetActivityByAid(aid int64) (*model.ActivityPO, error) {
	var res *model.ActivityPO
	dbRes := db.Model(&model.ActivityPO{}).Where("aid = ?", aid).Find(&res)
	if dbRes.Error != nil {
		return nil, dbRes.Error
	}
	return res, nil
}

func InsertActivity(activity *model.ActivityPO) error {
	result := db.Create(activity)
	return result.Error
}
