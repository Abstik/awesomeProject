package handler

import (
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"

	"awesomeProject/utils"
)

// 定义图片存储目录
const uploadDir = "./img/"
const watermarkPath = "./watermark.png" // 固定水印图片路径

// 初始化
func init() {
	// 创建存储目录（如果不存在）
	if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
		err := os.MkdirAll(uploadDir, os.ModePerm)
		if err != nil {
			log.Fatalf("Failed to create upload directory: %v\n", err)
		}
	}
}

// UploadImgWithWaterMark 处理图片上传
func UploadImgWithWaterMark(c *gin.Context) {
	// 获取上传的文件
	file, err := c.FormFile("image")
	if err != nil {
		utils.BuildErrorResponse(c, 400, "Failed to get uploaded file")
		return
	}

	// 打开上传的文件
	srcFile, err := file.Open()
	if err != nil {
		utils.BuildErrorResponse(c, 500, "Failed to open uploaded file")
		return
	}
	defer srcFile.Close()

	// 解码上传的图片（支持 JPEG 和 PNG）
	img, format, err := image.Decode(srcFile)
	if err != nil {
		utils.BuildErrorResponse(c, 400, "Unsupported image format")
		return
	}

	// 构造目标文件名（随机生成，确保唯一性）
	filename := fmt.Sprintf("watermarked_%d.%s", time.Now().Unix(), format)
	outputPath := filepath.Join(uploadDir, filename)

	// 保存图片
	outFile, err := os.Create(outputPath)
	if err != nil {
		utils.BuildErrorResponse(c, 500, "Failed to save watermarked image")
		return
	}
	defer outFile.Close()

	// 根据图片格式保存
	if format == "jpeg" || format == "jpg" {
		err = jpeg.Encode(outFile, img, &jpeg.Options{Quality: 90})
	} else if format == "png" {
		err = png.Encode(outFile, img)
	} else {
		utils.BuildErrorResponse(c, 400, "Unsupported image format")
		return
	}
	if err != nil {
		utils.BuildErrorResponse(c, 500, "Failed to encode watermarked image")
		return
	}

	// 返回图片的 URL
	utils.BuildSuccessResponse(c, gin.H{
		"success": true,
		"url":     "/img/" + filename,
	})
}

func DeleteImg(c *gin.Context) {
	url := c.Query("url")
	if url == "" {
		utils.BuildErrorResponse(c, 400, "url 为必传参数 请传递url")
		return
	}

	err := os.Remove(filepath.Join(".", url))
	if err != nil {
		utils.BuildErrorResponse(c, 500, "删除图片失败")
		return
	}
	utils.BuildSuccessResponse(c, gin.H{
		"success": true,
	})
}
