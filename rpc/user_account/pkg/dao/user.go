package dao

import (
	"encoding/json"
	"errors"
	"log"
	"time"

	"github.com/youperceive/cloudwego_instance/rpc/user_account/kitex_gen/base"
	"github.com/youperceive/cloudwego_instance/rpc/user_account/pkg/mysql"
	"gorm.io/gorm"
)

var (
	ErrUserNotFound = errors.New("user not found")
	ErrDBQuery      = errors.New("mysql query failed")
	ErrDBUpdate     = errors.New("mysql update failed")
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

func QueryUser(target string, targetType base.TargetType) (*User, error) {
	if target == "" {
		return nil, errors.New("query user failed: target is empty")
	}
	if targetType != base.TargetType_Email && targetType != base.TargetType_Phone {
		return nil, errors.New("query user failed: invalid target type")
	}

	user := &User{}
	var result *gorm.DB

	switch targetType {
	case base.TargetType_Phone:
		result = mysql.DB.Model(&User{}).Where("phone = ?", target).First(user)
	case base.TargetType_Email:
		result = mysql.DB.Model(&User{}).Where("email = ?", target).First(user)
	}

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, ErrDBQuery
	}

	return user, nil
}

func QueryUserById(id int64) (*User, error) {
	user := &User{}

	result := mysql.DB.First(user, id)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, ErrDBQuery
	}

	return user, nil
}

func UpdateUser(userId int64, info map[string]any) error {
	if len(info) == 0 {
		return nil
	}

	result := mysql.DB.Model(&User{}).Where("id = ?", userId).Updates(info)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return ErrUserNotFound
		}
		return ErrDBUpdate
	}

	return nil
}
