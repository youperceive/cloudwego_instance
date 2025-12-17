package mysql

import (
	"log"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func init() {
	dsn := os.Getenv("MYSQL_DSN")
	if dsn == "" {
		log.Println("ERROR: 环境变量 MYSQL_DSN 未配置，使用默认配置")
		dsn = "user_service:user123456@tcp(mysql:3306)/user_account_db?charset=utf8mb4&parseTime=True&loc=Local"
	}

	logLevel := logger.Silent
	if os.Getenv("ENV") == "dev" {
		logLevel = logger.Info
	}
	gormLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold: time.Second,
			LogLevel:      logLevel,
			Colorful:      true,
		},
	)

	rdb, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: gormLogger,
	})
	if err != nil {
		log.Fatalf("ERROR: 连接MySQL失败: %v", err)
	}

	sqlDB, err := rdb.DB()
	if err != nil {
		log.Fatalf("ERROR: 获取SQL连接池失败: %v", err)
	}
	sqlDB.SetMaxOpenConns(10)
	sqlDB.SetMaxIdleConns(20)
	sqlDB.SetConnMaxLifetime(time.Minute)
	sqlDB.SetConnMaxIdleTime(time.Minute)

	if err := sqlDB.Ping(); err != nil {
		log.Fatalf("ERROR: MySQL连接健康检查失败: %v", err)
	}

	DB = rdb
	log.Println("INFO: MySQL连接初始化成功")
}
