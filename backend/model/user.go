package model

import "github.com/dromara/carbon/v2"

type UserStatus int8

const (
	UserStatusDefault UserStatus = 0 // 正常
	UserStatusBan     UserStatus = 2 // 封禁
)

type User struct {
	ID     uint       `json:"-" gorm:"primarykey;"`
	Name   string     `gorm:"size:255;not null;comment:用户名"` // 名字
	Status UserStatus `gorm:"tinyint(1); not null;default:0;comment:用户状态"`

	CreatedAt carbon.DateTime `json:"created_at" gorm:"not null;default:CURRENT_TIMESTAMP(3);index"`
	UpdatedAt carbon.DateTime `json:"updated_at" gorm:"not null;default:CURRENT_TIMESTAMP(3);index"`
}
