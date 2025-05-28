package handler

import (
	"fmt"
	"image"
	"image/draw"
	"image/jpeg"
	"image/png"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/disintegration/imaging" // 用于简化图像操作
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

	/*// 加载水印图片
	watermarkImg, err := loadWatermark()
	if err != nil {
		utils.BuildErrorResponse(c, 500, "Failed to load watermark image")
		return
	}*/

	/*// 将水印叠加到图片
	watermarkedImg := addWatermark(img, watermarkImg)*/

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

// 加载水印图片
func loadWatermark() (image.Image, error) {
	// 打开水印图片文件
	file, err := os.Open(watermarkPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open watermark file: %v", err)
	}
	defer file.Close()

	// 解码水印图片
	watermarkImg, _, err := image.Decode(file)
	if err != nil {
		return nil, fmt.Errorf("failed to decode watermark image: %v", err)
	}

	return watermarkImg, nil
}

// 添加水印到图片右下角
func addWatermark(baseImg image.Image, watermarkImg image.Image) image.Image {
	// 获取图片的宽高
	baseBounds := baseImg.Bounds()
	baseWidth, baseHeight := baseBounds.Dx(), baseBounds.Dy()

	// 克隆原始图片为 RGBA 格式
	rgba := imaging.Clone(baseImg)

	// 获取水印宽高
	watermarkBounds := watermarkImg.Bounds()
	watermarkWidth, watermarkHeight := watermarkBounds.Dx(), watermarkBounds.Dy()

	// 计算水印放置的位置（右下角）
	offsetX := baseWidth - watermarkWidth - 10   // 距离右边缘 10 像素
	offsetY := baseHeight - watermarkHeight - 10 // 距离下边缘 10 像素

	// 在原图上叠加水印
	draw.Draw(
		rgba, // 目标图像
		image.Rect(offsetX, offsetY, offsetX+watermarkWidth, offsetY+watermarkHeight), // 水印区域
		watermarkImg,      // 水印图像
		image.Point{0, 0}, // 水印图像的起始点
		draw.Over,         // 绘制模式（叠加）
	)

	return rgba
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
