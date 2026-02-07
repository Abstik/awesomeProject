package handler

import (
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"awesomeProject/utils"
)

// 定义图片存储目录
const uploadDir = "./res/"

// 初始化
func init() {
	// 创建存储目录（如果不存在）
	if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
		if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
			panic(fmt.Sprintf("无法创建图片存储目录: %v", err))
		}
	}
}

// UploadImgWithWaterMark 处理图片上传
func UploadImgWithWaterMark(c *gin.Context) {
	// 获取上传的文件
	file, err := c.FormFile("image")
	if err != nil {
		utils.BuildErrorResponse(c, 400, "无法获取上传的图片文件")
		return
	}

	// 打开上传的文件
	srcFile, err := file.Open()
	if err != nil {
		utils.BuildServerError(c, "打开上传文件失败", err)
		return
	}
	defer srcFile.Close()

	// 解码上传的图片（支持 JPEG 和 PNG）
	img, format, err := image.Decode(srcFile)
	if err != nil {
		utils.BuildErrorResponse(c, 400, "不支持的图片格式")
		return
	}

	// 构造目标文件名（使用纳秒时间戳+随机数，确保唯一性）
	filename := fmt.Sprintf("%d_%s.%s", time.Now().UnixNano(), utils.Rand5Digits(), format)
	outputPath := filepath.Join(uploadDir, filename)

	// 保存图片
	outFile, err := os.Create(outputPath)
	if err != nil {
		utils.BuildServerError(c, "保存图片失败", err)
		return
	}
	defer outFile.Close()

	// 根据图片格式保存
	if format == "jpeg" || format == "jpg" {
		err = jpeg.Encode(outFile, img, &jpeg.Options{Quality: 90})
	} else if format == "png" {
		err = png.Encode(outFile, img)
	} else {
		utils.BuildErrorResponse(c, 400, "不支持的图片格式")
		return
	}
	if err != nil {
		utils.BuildServerError(c, "编码图片失败", err)
		return
	}

	filePath := uploadDir[1:] + filename

	// 返回图片的 URL
	utils.BuildSuccessResponse(c, gin.H{
		"success": true,
		"url":     utils.FullURL(&filePath),
	})
}

func DeleteImg(c *gin.Context) {
	imgUrl := c.Query("url")
	if imgUrl == "" {
		utils.BuildErrorResponse(c, 400, "url 为必传参数 请传递url")
		return
	}

	// 解析 URL
	parsedURL, err := url.Parse(imgUrl)
	if err != nil || parsedURL.Path == "" {
		utils.BuildErrorResponse(c, 400, "非法的 URL")
		return
	}

	// 仅允许 /res/ 开头的路径
	if !strings.HasPrefix(parsedURL.Path, "/res/") {
		utils.BuildErrorResponse(c, 400, "非法的资源路径")
		return
	}

	// 拼接本地文件路径
	localPath := filepath.Join(".", parsedURL.Path) // 即 ./res/xxx.jpg

	// 绝对路径二次校验，防止路径穿越
	absLocal, err := filepath.Abs(localPath)
	if err != nil {
		utils.BuildErrorResponse(c, 400, "路径解析失败")
		return
	}
	absRes, _ := filepath.Abs(uploadDir)
	if !strings.HasPrefix(absLocal, absRes) {
		utils.BuildErrorResponse(c, 400, "非法的资源路径")
		return
	}

	// 删除文件
	err = os.Remove(localPath)
	if err != nil {
		utils.BuildServerError(c, "删除图片失败", err)
		return
	}

	utils.BuildSuccessResponse(c, gin.H{
		"success": true,
	})
}
