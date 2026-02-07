package handler

import (
	"fmt"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"

	"awesomeProject/utils"
)

// 视频存储根目录
const videoUploadDir = "./videos/"

func init() {
	if _, err := os.Stat(videoUploadDir); os.IsNotExist(err) {
		if err := os.MkdirAll(videoUploadDir, os.ModePerm); err != nil {
			panic(fmt.Sprintf("无法创建视频存储目录: %v", err))
		}
	}
}

// isSafePath 检查路径是否安全（不包含路径穿越字符）
func isSafePath(name string) bool {
	cleaned := filepath.Clean(name)
	if strings.Contains(cleaned, "..") || strings.ContainsAny(name, `/\`) {
		return false
	}
	return true
}

// ------------------- 上传视频 -------------------
func UploadOrUpdateVideo(c *gin.Context) {
	// 获取类别名称
	category := c.PostForm("category")
	if category == "" {
		utils.BuildErrorResponse(c, http.StatusBadRequest, "category 为必传参数")
		return
	}

	// 防止路径穿越
	if !isSafePath(category) {
		utils.BuildErrorResponse(c, http.StatusBadRequest, "category 包含非法字符")
		return
	}

	// 获取上传的文件
	file, err := c.FormFile("video")
	if err != nil {
		utils.BuildErrorResponse(c, http.StatusBadRequest, "无法获取上传的视频文件")
		return
	}

	// 防止文件名路径穿越
	if !isSafePath(file.Filename) {
		utils.BuildErrorResponse(c, http.StatusBadRequest, "文件名包含非法字符")
		return
	}

	// 打开目录
	categoryDir := filepath.Join(videoUploadDir, category)
	if _, err := os.Stat(categoryDir); os.IsNotExist(err) {
		if err := os.MkdirAll(categoryDir, os.ModePerm); err != nil {
			utils.BuildServerError(c, "创建类别目录失败", err)
			return
		}
	}

	// 获取文件名
	filename := file.Filename
	// 文件的相对路径
	outputPath := filepath.Join(categoryDir, filename)
	if err := c.SaveUploadedFile(file, outputPath); err != nil {
		utils.BuildServerError(c, "保存视频文件失败", err)
		return
	}
	// 去除文件路径的.前缀
	urlPath := path.Join("/videos", category, filename) // 注意使用 path.Join（URL 路径，使用 '/'）
	utils.BuildSuccessResponse(c, gin.H{
		"category": category,
		"name":     strings.TrimSuffix(filename, filepath.Ext(filename)),
		"url":      utils.FullURL(&urlPath),
	})
}

// ------------------- 删除视频 -------------------
func DeleteVideoByURL(c *gin.Context) {
	// 前端传来的 URL
	videoURL := c.Query("url")
	if videoURL == "" {
		utils.BuildErrorResponse(c, http.StatusBadRequest, "url 为必传参数")
		return
	}

	parsedPath := utils.ParseURL(&videoURL)
	if parsedPath == nil || *parsedPath == "" {
		utils.BuildErrorResponse(c, http.StatusBadRequest, "非法的 URL")
		return
	}

	// 仅允许 /videos/ 开头的路径，防止路径穿越
	if !strings.HasPrefix(*parsedPath, "/videos/") {
		utils.BuildErrorResponse(c, http.StatusBadRequest, "非法的资源路径")
		return
	}

	localPath := filepath.Join(".", *parsedPath)

	// 解析后再次检查是否在 videos 目录内
	absLocal, err := filepath.Abs(localPath)
	if err != nil {
		utils.BuildErrorResponse(c, http.StatusBadRequest, "路径解析失败")
		return
	}
	absVideos, _ := filepath.Abs(videoUploadDir)
	if !strings.HasPrefix(absLocal, absVideos) {
		utils.BuildErrorResponse(c, http.StatusBadRequest, "非法的资源路径")
		return
	}

	// 检查文件是否存在
	if _, err := os.Stat(localPath); os.IsNotExist(err) {
		utils.BuildErrorResponse(c, http.StatusNotFound, "文件不存在")
		return
	}

	// 删除文件
	if err := os.Remove(localPath); err != nil {
		utils.BuildServerError(c, "删除视频失败", err)
		return
	}

	utils.BuildSuccessResponse(c, nil)
}

// ------------------- 获取所有类别视频 -------------------
func GetAllVideos(c *gin.Context) {
	categories, err := os.ReadDir(videoUploadDir)
	if err != nil {
		utils.BuildServerError(c, "读取视频目录失败", err)
		return
	}
	result := []gin.H{}

	for _, cat := range categories {
		if !cat.IsDir() {
			continue
		}

		categoryDir := filepath.Join(videoUploadDir, cat.Name())
		files, err := os.ReadDir(categoryDir)
		if err != nil {
			continue
		}
		if len(files) == 0 {
			continue
		}

		// 用来装该类别下的所有视频
		videos := []gin.H{}
		for _, f := range files {
			if f.IsDir() {
				continue // 如果还有子目录就跳过
			}

			videoFile := f.Name()
			filePath := filepath.Join(categoryDir, videoFile)
			fileUrl := strings.TrimPrefix(filePath, ".")

			videos = append(videos, gin.H{
				"name": strings.TrimSuffix(videoFile, filepath.Ext(videoFile)),
				"url":  utils.FullURL(&fileUrl),
			})

		}
		result = append(result, gin.H{
			"category": cat.Name(),
			"videos":   videos,
		})
	}

	utils.BuildSuccessResponse(c, result)
}
