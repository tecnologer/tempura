package models

import (
	"time"

	"github.com/tecnologer/tempura/pkg/models/notification"
	"gorm.io/gorm"
)

type NotificationSetting struct {
	*gorm.Model
	UserID         uint
	TelegramChatID int64
	Type           notification.Type `gorm:"type:varchar(255)"`
	Delay          time.Duration
	User           *User `gorm:"foreignKey:UserID"`
}
