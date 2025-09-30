package utils

import (
	"fmt"
	"net/url"
	"os"
	"strings"
)

/*
分解url和拼接url：确保数据库中只存储/res/文件名，而前端获取到https://域名/res/文件名
*/

// 给图片的相对路径加上前缀，构成完整的图片URL
func FullURL(path *string) *string {
	if path == nil {
		return nil
	}

	if strings.HasPrefix(*path, "http") {
		return path // 已是完整路径
	}

	// path是 /res/文件名，domain是域名
	domain := os.Getenv("DOMAIN_NAME")
	if domain == "" {
		domain = "127.0.0.1:8080"
		//domain = "mobile.xupt.edu.cn"
	}

	// todo 改为https
	fullURL := fmt.Sprintf("http://%s%s", domain, *path)
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
