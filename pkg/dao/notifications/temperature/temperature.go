package temperature

import (
	"fmt"

	"github.com/tecnologer/tempura/pkg/casbin"
	"github.com/tecnologer/tempura/pkg/dao/notifications/helper"
	"github.com/tecnologer/tempura/pkg/dao/notifications/notification"
	"github.com/tecnologer/tempura/pkg/models/notificationtype"
)

const (
	prefixEmoji     = "🌡"
	suffixEmojiLow  = " ❄️️"
	suffixEmojiMuch = " 🔥"
	noSuffixEmoji   = ""
)

type Notification struct {
	*notification.Notification
	normalRange *helper.NotificationState
	lowRange    *helper.NotificationState
	highRange   *helper.NotificationState
}

func NewNotification(enforcer *casbin.Enforcer) *Notification {
	return &Notification{
		Notification: notification.NewNotification(enforcer),
		normalRange:  helper.NewNotificationState(notificationtype.TemperatureNormal),
		lowRange:     helper.NewNotificationState(notificationtype.TemperatureLow),
		highRange:    helper.NewNotificationState(notificationtype.TemperatureHigh),
	}
}

func (n *Notification) Get(batteryLevel float64, chatID int64) string {
	n.Init(batteryLevel, chatID)

	if n.HasNotification(n.normalRange) {
		return n.buildMessage(noSuffixEmoji)
	}

	if n.HasNotification(n.lowRange) {
		return n.buildMessage(suffixEmojiLow)
	}

	if n.HasNotification(n.highRange) {
		return n.buildMessage(suffixEmojiMuch)
	}

	return ""
}

func (n *Notification) buildMessage(suffix string) string {
	return fmt.Sprintf("%s Temperature is %.2f°C%s\n", prefixEmoji, n.Value, suffix)
}
