package helper

import (
	"time"

	"github.com/tecnologer/tempura/pkg/models/notificationtype"
)

type NotificationState struct {
	NotifType        notificationtype.Type
	InRange          bool
	LastNotification time.Time
	Counter          int
}

func NewNotificationState(notifType notificationtype.Type) *NotificationState {
	return &NotificationState{
		NotifType: notifType,
	}
}
