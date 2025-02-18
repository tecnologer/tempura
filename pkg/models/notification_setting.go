package models

import (
	"time"

	"github.com/tecnologer/tempura/pkg/models/notification"
	"gorm.io/gorm"
)

type NotificationSettings []*NotificationSetting

func (n NotificationSettings) SettingsByType(notifType notificationtype.Type) *NotificationSetting {
	for _, setting := range n {
		if setting.Type == notifType || setting.Type == notificationtype.TypeAll {
			return setting
		}
	}

	return nil
}

type NotificationSetting struct {
	*gorm.Model
	UserID uint
	Type   notificationtype.Type `gorm:"type:varchar(255)"`
	Delay  time.Duration
	User   *User `gorm:"foreignKey:UserID"`
}

func (n *NotificationSetting) IsEnabled() bool {
	return n != nil && n.Delay > 0
}

func (n *NotificationSetting) IsTimeToSend(last time.Time) bool {
	return n.IsEnabled() && time.Since(last) > n.Delay
}
