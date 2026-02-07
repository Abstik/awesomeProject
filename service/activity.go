package service

import (
	"regexp"
	"time"

	"awesomeProject/dao"
	"awesomeProject/model"
	"awesomeProject/utils"
)

// 预编译正则，避免每次调用重新编译
var (
	imgSrcRe  = regexp.MustCompile(`(?i)(<img[^>]+src=["'])(/res/[^"']+)(["'])`)
	fullURLRe = regexp.MustCompile(`(?i)(src|poster)=["']https?://[^/]+(/res/[^"']+)["']`)
)

func GetActivityList(pageSize, pageNum int) ([]*model.ActivityPO, int64, error) {
	// 分页查询活动
	res, total, err := dao.GetActivityListByPage(pageSize, pageNum)
	if err != nil {
		return nil, 0, err
	}
	for _, v := range res {
		v.Img = utils.OldFullURL(v.Img)
	}
	return res, total, nil
}

func GetActivityByAid(aid int64) (*model.ActivityPO, error) {
	res, err := dao.GetActivityByAid(aid)
	if err != nil {
		return nil, err
	}

	if res.Content != nil {
		content := FullImageURLs(*res.Content)
		res.Content = &content
	}

	res.Img = utils.OldFullURL(res.Img)
	return res, nil
}

// 替换html中图片的相对路径为绝对路径
func FullImageURLs(input string) string {
	return imgSrcRe.ReplaceAllStringFunc(input, func(match string) string {
		subMatches := imgSrcRe.FindStringSubmatch(match)
		if len(subMatches) != 4 {
			return match
		}

		fullURL := utils.OldFullURL(&subMatches[2])
		url := subMatches[1] + *fullURL + subMatches[3]
		return url
	})
}

// 替换html中图片绝对路径为相对路径
func ParseImageURLS(html string) string {
	return fullURLRe.ReplaceAllStringFunc(html, func(match string) string {
		subMatches := fullURLRe.FindStringSubmatch(match)
		if len(subMatches) != 3 {
			return match
		}
		attr := subMatches[1]    // src 或 poster
		relPath := subMatches[2] // /res/xxx

		// 返回： src="/res/xxx" 或 poster="/res/xxx"
		return attr + `="` + relPath + `"`
	})
}

func AddActivity(req *model.ActivityReq) error {
	content := ParseImageURLS(*req.Content)

	activity := &model.ActivityPO{
		Title:   req.Title,
		Summary: req.Summary,
		Content: &content,
		Img:     utils.ParseURL(req.Img),
	}
	now := time.Now()
	activity.Time = &now
	return dao.InsertActivity(activity)
}

func UpdateActivity(req *model.ActivityReq) error {
	activity := &model.ActivityPO{AID: req.AID}
	if req.Img != nil {
		activity.Img = utils.ParseURL(req.Img)
	}
	if req.Title != nil {
		activity.Title = req.Title
	}
	if req.Summary != nil {
		activity.Summary = req.Summary
	}
	if req.Content != nil {
		content := ParseImageURLS(*req.Content)
		activity.Content = &content
	}
	return dao.UpdateActivity(activity)
}

// DeleteActivity 删除活动
func DeleteActivity(aid int64) error {
	return dao.DeleteActivity(aid)
}
