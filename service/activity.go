package service

import (
	"time"

	"awesomeProject/dao"
	"awesomeProject/model"
)

func GetActivityList(pageSize *int, pageNum *int) ([]*model.ActivityPO, int64, error) {
	// 分页查询活动
	res, total, err := dao.GetActivityListByPage(*pageSize, *pageNum)
	if err != nil {
		return nil, 0, err
	}
	return res, total, nil
}

func GetActivityByAid(aid int64) (*model.ActivityPO, error) {
	res, err := dao.GetActivityByAid(aid)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func AddActivity(req *model.AddActivityReq) error {
	activity := &model.ActivityPO{
		Title:   req.Title,
		Summary: req.Summary,
		Content: req.Content,
		Img:     req.Img,
	}
	now := time.Now()
	activity.Time = &now
	err := dao.InsertActivity(activity)
	if err != nil {
		return err
	}
	return nil
}
