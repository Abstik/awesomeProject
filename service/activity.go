package service

import (
	"regexp"
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

	content := sanitizeHTML(*res.Content)
	res.Content = &content
	return res, nil
}

func sanitizeHTML(input string) string {
	/*// 匹配 <img ... src="..." ...> 标签
	imgTagRegex := regexp.MustCompile(`<img[^>]*src=["']([^"']+)["'][^>]*>`)

	// 替换为 <p style="text-align:center;"><img src="..." /></p>
	safeHTML := imgTagRegex.ReplaceAllStringFunc(input, func(imgTag string) string {
		matches := imgTagRegex.FindStringSubmatch(imgTag)
		if len(matches) < 2 {
			return "" // 找不到 src 就不保留了
		}
		src := matches[1]
		return fmt.Sprintf(`<p style="text-align:center;"><img src="%s" /></p>`, src)
	})*/

	// 正则表达式：匹配被任意标签包裹的 <img> 标签
	re := regexp.MustCompile(`(?i)<[^>]+>\s*(<img[^>]*>)\s*</[^>]+>`)
	// 替换为 <p><img ... /></p>
	result := re.ReplaceAllString(input, "<p>$1</p>")

	return result
}

func AddActivity(req *model.ActivityReq) error {
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

func UpdateActivity(req *model.ActivityReq) error {
	activity := &model.ActivityPO{AID: req.AID}
	if req.Img != nil {
		activity.Img = req.Img
	}
	if req.Title != nil {
		activity.Title = req.Title
	}
	if req.Summary != nil {
		activity.Summary = req.Summary
	}
	if req.Content != nil {
		activity.Content = req.Content
	}
	err := dao.UpdateActivity(activity)
	return err
}
