package service

import (
	"regexp"
	"time"

	"awesomeProject/dao"
	"awesomeProject/model"
	"awesomeProject/utils"
)

func GetActivityList(pageSize *int, pageNum *int) ([]*model.ActivityPO, int64, error) {
	// 分页查询活动
	res, total, err := dao.GetActivityListByPage(*pageSize, *pageNum)
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

	content := FullImageURLs(*res.Content)
	res.Content = &content

	res.Img = utils.OldFullURL(res.Img)
	return res, nil
}

// 替换html中图片的相对路径为绝对路径
func FullImageURLs(input string) string {
	re := regexp.MustCompile(`(?i)(<img[^>]+src=["'])(/res/[^"']+)(["'])`)

	return re.ReplaceAllStringFunc(input, func(match string) string {
		subMatches := re.FindStringSubmatch(match)
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
	// 匹配 src="http://域名/res/xxx" 或 poster="http://域名/res/xxx"
	re := regexp.MustCompile(`(?i)(src|poster)=["']https?://[^/]+(/res/[^"']+)["']`)

	return re.ReplaceAllStringFunc(html, func(match string) string {
		subMatches := re.FindStringSubmatch(match)
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
	err := dao.InsertActivity(activity)
	return err
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
	err := dao.UpdateActivity(activity)
	return err
}
