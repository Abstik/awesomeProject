package dao

import (
	"fmt"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var db *gorm.DB

func InitDatabaseConnector() error {
	// dsn := "root:325523@tcp(127.0.0.1:3306)/xiyoumobile_data?charset=utf8mb4&parseTime=True&loc=Local"
	// todo：部署到服务器后，使用以下配置
	dsn := fmt.Sprintf("root:%s@tcp(%s:%s)/xiyoumobile_data?charset=utf8mb4&parseTime=True",
		os.Getenv("DB_PASSWORD"), // 密码
		"mysql",                  // 主机（必须改为 "mysql"）
		os.Getenv("DB_PORT"))
	var err error
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}
	return nil
}
