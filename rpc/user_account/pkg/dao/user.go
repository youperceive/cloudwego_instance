package dao

import (
	"encoding/json"
	"log"
	"time"

	"github.com/cloudwego_instance/rpc/user_account/pkg/mysql"
)

type User struct {
	ID           int64           `gorm:"primarykey;column:id;type:bigint unsigned;autoIncrement"`
	Username     string          `gorm:"column:username;type:varchar(64);default:''"`
	Email        *string         `gorm:"column:email;type:varchar(128);uniqueIndex:uk_email"`
	Phone        *string         `gorm:"column:phone;type:varchar(20);uniqueIndex:uk_phone"`
	Password     string          `gorm:"column:password;type:varchar(255);not null"`
	RegisterType int8            `gorm:"column:register_type;type:tinyint;not null"`
	UserType     int8            `gorm:"column:user_type;type:tinyint;not null;default:1"`
	Status       int8            `gorm:"column:status;type:tinyint;not null;default:1"`
	Ext          json.RawMessage `gorm:"column:ext;type:json"`
	CreatedAt    int64           `gorm:"column:created_at;type:bigint;not null"`
	UpdatedAt    int64           `gorm:"column:updated_at;type:bigint;not null"`
}

func (u *User) TableName() string {
	return "user"
}

func CreateUser(user *User) (int64, error) {
	now := time.Now().Unix()
	user.CreatedAt = now
	user.UpdatedAt = now
	if user.Ext == nil {
		user.Ext = json.RawMessage("{}")
	}

	if err := mysql.DB.Create(user).Error; err != nil {
		log.Println(err)
		return 0, err
	}

	return user.ID, nil
}
