package product

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

// DB 全局数据库连接（仅初始化，无业务逻辑）
var DB *sql.DB

func init() {
	InitDB()
}

func InitDB() {
	// 从环境变量读取，兜底硬编码
	dsn := os.Getenv("MYSQL_DSN")
	if dsn == "" {
		dsn = "root:你的密码@tcp(127.0.0.1:3306)/你的数据库名?charset=utf8mb4&parseTime=True"
	}

	var err error
	DB, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("product: 数据库连接失败: %v", err)
	}

	if err = DB.Ping(); err != nil {
		log.Fatalf("product: 数据库 Ping 失败: %v", err)
	}
	log.Println("product: 数据库连接成功")
}
