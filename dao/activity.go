package dao

import "awesomeProject/model"

func GetActivityListByPage(pageSize, pageNum int) ([]*model.ActivityPO, error) {
	var res []*model.ActivityPO
	offset := (pageSize - 1) * pageNum
	dbRes := db.Model(&model.ActivityPO{}).Order("time DESC").Offset(offset).Limit(pageNum).Find(&res)
	if dbRes.Error != nil {
		return nil, dbRes.Error
	}
	return res, nil
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
