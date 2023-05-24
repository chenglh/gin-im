package initzlize

import (
	"IM/gin_im/global"
	"log"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

// 初始化数据库
func InitDB() {
	//dsn := "root:12345678@tcp(127.0.0.1:3306)/mkshop_user_srv?charset=utf8mb4&parseTime=True&loc=Local"
	dsn := "root:12345678@unix(/tmp/mysql.sock)/chat?charset=utf8mb4&parseTime=True&loc=Local"

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold: time.Second, //慢SQL阈值
			LogLevel:      logger.Info, //Log level
			Colorful:      true,        //启用彩色打印
		},
	)

	//全局模式
	var err error
	global.DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   "im_", //表前缀
			SingularTable: true,  //表单数
		},
		Logger: newLogger,
	})
	if err != nil {
		panic(err)
	}
}
