package main

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/xuri/excelize/v2"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"awesomeProject/model"
	"awesomeProject/utils"
)

func main() {
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
