package service

import (
	"errors"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"awesomeProject/dao"
	"awesomeProject/model"
	"awesomeProject/utils"
)

// 预编译正则，避免每次调用重新编译
var (
	imgSrcRe   = regexp.MustCompile(`(?i)(<img[^>]+src=["'])(/res/[^"']+)(["'])`)
	fullURLRe  = regexp.MustCompile(`(?i)(src|poster)=["']https?://[^/]+(/res/[^"']+)["']`)
	imageURLRe = regexp.MustCompile(`(?i)\b(?:src|poster)\s*=\s*["']([^"']+)["']`)
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
	activity, err := dao.GetActivityByAid(aid)
	if err != nil {
		return err
	}
	activities, err := dao.GetActivityList()
	if err != nil {
		return err
	}

	deletedImages := unsharedActivityImagePaths(activity, activities)

	if err := dao.DeleteActivity(aid); err != nil {
		return err
	}
	return removeActivityImages(deletedImages)
}

func removeActivityImages(images map[string]struct{}) error {
	var cleanupErrors []error
	for imagePath := range images {
		info, err := os.Lstat(imagePath)
		if errors.Is(err, os.ErrNotExist) {
			continue
		}
		if err != nil {
			cleanupErrors = append(cleanupErrors, err)
			continue
		}
		if !info.Mode().IsRegular() {
			continue
		}
		if err := os.Remove(imagePath); err != nil {
			cleanupErrors = append(cleanupErrors, err)
		}
	}
	return errors.Join(cleanupErrors...)
}

func unsharedActivityImagePaths(activity *model.ActivityPO, activities []*model.ActivityPO) map[string]struct{} {
	paths := activityImagePaths(activity)
	for _, other := range activities {
		if other.AID == activity.AID {
			continue
		}
		for path := range activityImagePaths(other) {
			delete(paths, path)
		}
	}
	return paths
}

// activityImagePaths returns local files in ./res referenced by an activity.
// Paths outside /res and nested paths are ignored so database content cannot
// cause arbitrary files to be removed.
func activityImagePaths(activity *model.ActivityPO) map[string]struct{} {
	paths := make(map[string]struct{})
	if activity == nil {
		return paths
	}
	configuredHost := ""
	if configuredDomain, err := url.Parse("//" + os.Getenv("DOMAIN_NAME")); err == nil {
		configuredHost = configuredDomain.Hostname()
	}
	addURL := func(raw string) {
		parsed, err := url.Parse(raw)
		if err != nil {
			return
		}
		if parsed.Host != "" {
			host := parsed.Hostname()
			if !strings.EqualFold(host, "mobile.xupt.edu.cn") && !strings.EqualFold(host, configuredHost) {
				return
			}
		}
		resourcePath := parsed.Path
		if !strings.HasPrefix(resourcePath, "/res/") {
			return
		}
		name := strings.TrimPrefix(resourcePath, "/res/")
		if name == "" || strings.ContainsAny(name, `/\`) || name == "." || name == ".." {
			return
		}
		switch strings.ToLower(filepath.Ext(name)) {
		case ".jpg", ".jpeg", ".png", ".gif", ".webp", ".avif":
		default:
			return
		}
		paths[filepath.Join("res", name)] = struct{}{}
	}
	if activity.Img != nil {
		addURL(*activity.Img)
	}
	if activity.Content != nil {
		for _, match := range imageURLRe.FindAllStringSubmatch(*activity.Content, -1) {
			if len(match) > 1 {
				addURL(match[1])
			}
		}
	}
	return paths
}
