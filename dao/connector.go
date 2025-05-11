package dao

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var db *gorm.DB

func InitDatabaseConnector() error {
	// TODO 参数后面考虑配置化一下把，这里先简单写，另外放到Dao还是放到Util的合理性有待商榷
	dsn := "root:325523@tcp(127.0.0.1:3306)/xiyoumobile_data?charset=utf8mb4&parseTime=True&loc=Local"
	var err error
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}
	return nil
}
