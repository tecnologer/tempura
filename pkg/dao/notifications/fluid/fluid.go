package fluid

import (
	"github.com/tecnologer/tempura/pkg/casbin"
	"github.com/tecnologer/tempura/pkg/dao/notifications/helper"
	"github.com/tecnologer/tempura/pkg/dao/notifications/notification"
	"github.com/tecnologer/tempura/pkg/models/notificationtype"
)

const (
	prefixEmoji = "🛢️"
)

type Notification struct {
	*notification.Notification
	fullRange  *helper.NotificationState
	emptyRange *helper.NotificationState
}

func NewNotification(enforcer *casbin.Enforcer) *Notification {
	return &Notification{
		Notification: notification.NewNotification(enforcer),
		fullRange:    helper.NewNotificationState(notificationtype.FluidTankFull),
		emptyRange:   helper.NewNotificationState(notificationtype.FluidTankEmpty),
	}
}

func (n *Notification) Get(humidityLevel float64, chatID int64) string {
	n.Init(humidityLevel, chatID)

	if n.HasNotification(n.emptyRange) {
		return prefixEmoji + " Fluid tank is empty\n"
	}

	if n.HasNotification(n.fullRange) {
		return prefixEmoji + " Fluid tank is full\n"
	}

	return ""
}
