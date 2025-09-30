package handler

import (
	"fmt"
	"net/http"
	"os"
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

// ------------------- 上传视频 -------------------
func UploadOrUpdateVideo(c *gin.Context) {
	// 获取类别名称
	category := c.PostForm("category")
	if category == "" {
		utils.BuildErrorResponse(c, http.StatusBadRequest, "category 为必传参数")
		return
	}

	// 获取上传的文件
	file, err := c.FormFile("video")
	if err != nil {
		utils.BuildErrorResponse(c, http.StatusBadRequest, "无法获取上传的视频文件")
		return
	}

	// 打开目录
	categoryDir := filepath.Join(videoUploadDir, category)
	if _, err := os.Stat(categoryDir); os.IsNotExist(err) {
		if err := os.MkdirAll(categoryDir, os.ModePerm); err != nil {
			utils.BuildErrorResponse(c, http.StatusInternalServerError, "创建类别目录失败")
			return
		}
	}

	// 获取文件名
	filename := file.Filename
	// 文件的相对路径
	outputPath := filepath.Join(categoryDir, filename)
	if err := c.SaveUploadedFile(file, outputPath); err != nil {
		utils.BuildErrorResponse(c, http.StatusInternalServerError, "保存视频文件失败")
		return
	}
	// 去除文件路径的.前缀
	fileUrl := strings.TrimPrefix(outputPath, ".")

	utils.BuildSuccessResponse(c, gin.H{
		"category": category,
		"name":     strings.TrimSuffix(filename, filepath.Ext(filename)),
		"url":      utils.FullURL(&fileUrl),
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

	url := utils.ParseURL(&videoURL)
	localPath := "." + *url

	// 检查文件是否存在
	if _, err := os.Stat(localPath); os.IsNotExist(err) {
		utils.BuildErrorResponse(c, http.StatusNotFound, "文件不存在")
		return
	}

	// 删除文件
	if err := os.Remove(localPath); err != nil {
		utils.BuildErrorResponse(c, http.StatusInternalServerError, "删除文件失败")
		return
	}

	utils.BuildSuccessResponse(c, nil)
}

// ------------------- 获取所有类别视频 -------------------
func GetAllVideos(c *gin.Context) {
	categories, _ := os.ReadDir(videoUploadDir)
	result := []gin.H{}

	for _, cat := range categories {
		if !cat.IsDir() {
			continue
		}

		categoryDir := filepath.Join(videoUploadDir, cat.Name())
		files, _ := os.ReadDir(categoryDir)
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
