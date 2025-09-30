package main

import (
	"fmt"
	"log"
	"math/rand"
	"regexp"
	"strings"
	"time"

	"github.com/xuri/excelize/v2"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"awesomeProject/model"
	"awesomeProject/service"
	"awesomeProject/utils"
)

func main() {}

func BatchUpdateHTML() {
	dsn := "root:325523@tcp(127.0.0.1:3306)/xiyoumobile_data?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("数据库连接失败: %v", err)
	}

	var activities []model.ActivityPO
	err = db.Find(&activities).Error
	if err != nil {
		log.Fatalf("查询失败: %v", err)
	}

	for _, activity := range activities {
		if activity.Content == nil {
			continue
		}

		// 清洗 HTML
		//clean := sanitizeHTML(*activity.Content)
		clean := service.ParseImageURLS(*activity.Content)
		activity.Content = &clean

		// 更新单条记录
		if err := db.Model(&model.ActivityPO{}).
			Where("aid = ?", activity.AID).
			Update("content", activity.Content).Error; err != nil {
			log.Printf("更新失败 AID=%d: %v", activity.AID, err)
		}
	}
	log.Println("更新完成")
}

// 清洗HTML标签
func sanitizeHTML(input string) string {
	// 正则表达式：匹配被任意标签包裹的 <img> 标签
	re := regexp.MustCompile(`(?i)<[^>]+>\s*(<img[^>]*>)\s*</[^>]+>`)
	// 替换为 <p><img ... /></p>
	result := re.ReplaceAllString(input, "<p>$1</p>")
	return result
}

func BatchUpdateActivityImg() {
	dsn := "root:325523@tcp(127.0.0.1:3306)/xiyoumobile_data?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("数据库连接失败: %v", err)
	}

	var activities []model.ActivityPO

	// 查询所有活动记录（可根据实际情况加条件）
	err = db.Find(&activities).Error
	if err != nil {
		log.Fatalf("查询失败: %v", err)
	}

	var needUpdate []model.ActivityPO

	for _, act := range activities {
		if act.Img != nil && !strings.HasPrefix(*act.Img, "/") {
			newPath := "/" + *act.Img
			act.Img = &newPath
			needUpdate = append(needUpdate, act)
		}
	}

	// 批量更新
	for _, updated := range needUpdate {
		err := db.Model(&model.ActivityPO{}).
			Where("aid = ?", updated.AID).
			Update("img", updated.Img).Error
		if err != nil {
			log.Printf("更新 aid=%d 时失败: %v", updated.AID, err)
		}
	}

	fmt.Printf("共修改 %d 条记录\n", len(needUpdate))
}

func BatchUpdateMember() {
	dsn := "root:325523@tcp(127.0.0.1:3306)/xiyoumobile_data?charset=utf8mb4&parseTime=True&loc=Local"
	db, _ := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	var members []model.MemberPO

	err := db.Find(&members).Error
	if err != nil {
		log.Fatalf("查询失败: %v", err)
	}

	var needUpdate []model.MemberPO

	for _, m := range members {
		updated := false

		if m.Portrait != nil {
			if fixed := stripBeforeRes(*m.Portrait); fixed != *m.Portrait {
				m.Portrait = &fixed
				updated = true
			}
		}
		if m.MienImg != nil {
			if fixed := stripBeforeRes(*m.MienImg); fixed != *m.MienImg {
				m.MienImg = &fixed
				updated = true
			}
		}
		if m.GraduateImg != nil {
			if fixed := stripBeforeRes(*m.GraduateImg); fixed != *m.GraduateImg {
				m.GraduateImg = &fixed
				updated = true
			}
		}

		if updated {
			needUpdate = append(needUpdate, m)
		}
	}

	// 单条更新（安全）
	for _, m := range needUpdate {
		err := db.Model(&model.MemberPO{}).Where("uid = ?", m.UID).Updates(map[string]interface{}{
			"portrait":     m.Portrait,
			"mien_img":     m.MienImg,
			"graduate_img": m.GraduateImg,
		}).Error
		if err != nil {
			log.Printf("更新 id=%d 失败: %v", m.UID, err)
		}
	}

	fmt.Printf("共更新 %d 条记录\n", len(needUpdate))
}

func stripBeforeRes(input string) string {
	idx := strings.Index(input, "/res/")
	if idx == -1 {
		return input // 没找到，不修改
	}
	return input[idx:] // 从/res/开始截取
}

// 初始化全部成员的密码
func InitAllMemberPassword() {
	dsn := "root:325523@tcp(127.0.0.1:3306)/xiyoumobile_data?charset=utf8mb4&parseTime=True&loc=Local"
	db, _ := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	// 获取所有用户
	var members []model.MemberPO
	if err := db.Select("uid", "username", "name").Find(&members).Error; err != nil {
		log.Fatalf("查询用户失败: %v", err)
	}

	// 创建 Excel 文件
	f := excelize.NewFile()
	sheet := "Sheet1"
	f.SetSheetRow(sheet, "A1", &[]string{"Username", "RawPassword"})

	rand.Seed(time.Now().UnixNano()) // 随机数种子

	// 遍历用户生成密码、写入 Excel、更新数据库
	for i, m := range members {
		if m.Username == nil {
			continue
		}
		rawPassword := *m.Username + utils.Rand5Digits()
		encrypted := utils.EncryptPassword(rawPassword)

		// 更新数据库
		if err := db.Model(&model.MemberPO{}).Where("uid = ?", m.UID).Update("password", encrypted).Error; err != nil {
			log.Printf("更新用户 %s 密码失败: %v", *m.Username, err)
			continue
		}

		// 写入 Excel，第i+2行，因为第1行是表头
		cell := fmt.Sprintf("A%d", i+2)
		f.SetSheetRow(sheet, cell, &[]string{*m.Name, rawPassword})
	}

	// 保存 Excel
	if err := f.SaveAs("账号密码.xlsx"); err != nil {
		log.Fatalf("保存 Excel 文件失败: %v", err)
	}

	log.Println("完成：密码加密并写入 Excel 文件。")
}
