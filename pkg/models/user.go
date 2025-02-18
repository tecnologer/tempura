package models

import "gorm.io/gorm"

type User struct {
	*gorm.Model
	TelegramChatID      int64                `json:"telegram_chat_id"     gorm:"unique"`
	Username            string               `json:"username"             gorm:"unique"`
	Role                string               `json:"role"`
	NotificationSetting NotificationSettings `json:"notification_setting" gorm:"foreignKey:UserID"`
}
