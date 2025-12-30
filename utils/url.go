package utils

import (
	"fmt"
	"net/url"
	"os"
	"path"
	"strings"
)

/*
分解url和拼接url：确保数据库中只存储/res/文件名，而前端获取到https://域名/res/文件名
*/

// 给图片的相对路径加上前缀，构成完整的图片URL
func FullURL(p *string) *string {
	if p == nil || *p == "" {
		return nil
	}

	s := *p

	// 如果已经完整 URL，直接返回
	if strings.HasPrefix(s, "http://") || strings.HasPrefix(s, "https://") {
		return &s
	}

	//换反斜杠为正斜杠（Windows兼容）
	s = strings.ReplaceAll(s, "\\", "/")

	// 去除重复的斜杠
	s = path.Clean(s)

	// 确保前面有 "/"
	if !strings.HasPrefix(s, "/") {
		s = "/" + s
	}

	var fullURL string
	// path是 /res/文件名，domain是域名
	domain := os.Getenv("DOMAIN_NAME")
	if domain == "" || domain == "mobile.xupt.edu.cn" {
		fullURL = fmt.Sprintf("https://mobile.xupt.edu.cn%s", s)
	} else { // 开发环境，指定服务部署机器对应的ip
		fullURL = fmt.Sprintf("http://%s%s", domain, s)
	}
	return &fullURL
}

// 用于查询旧官网上的资源
func OldFullURL(p *string) *string {
	if p == nil || *p == "" {
		return nil
	}

	s := *p

	// 如果已经完整 URL，直接返回
	if strings.HasPrefix(s, "http://") || strings.HasPrefix(s, "https://") {
		return &s
	}

	//换反斜杠为正斜杠（Windows兼容）
	s = strings.ReplaceAll(s, "\\", "/")

	// 去除重复的斜杠
	s = path.Clean(s)

	// 确保前面有 "/"
	if !strings.HasPrefix(s, "/") {
		s = "/" + s
	}

	// path是 /res/文件名，domain是域名
	domain := "mobile.xupt.edu.cn"
	fullURL := fmt.Sprintf("https://%s%s", domain, s)
	return &fullURL
}

// 解析URL，去掉https和域名，只保留 /res/文件名
func ParseURL(s *string) *string {
	if s == nil {
		return nil
	}

	parsedURL, err := url.Parse(*s)
	if err != nil {
		return nil
	}
	return &parsedURL.Path
}
