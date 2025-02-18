package humidity

import (
	"fmt"
	"time"

	"github.com/tecnologer/tempura/pkg/casbin"
	"github.com/tecnologer/tempura/pkg/dao/notifications/helper"
	"github.com/tecnologer/tempura/pkg/dao/notifications/notification"
	"github.com/tecnologer/tempura/pkg/models/notificationtype"
)

const (
	prefixEmoji     = "💧"
	suffixEmojiLow  = " 🏜️️"
	suffixEmojiMuch = " 💦"
	noSuffixEmoji   = ""
)

type Notification struct {
	*notification.Notification
	normalRange *helper.NotificationState
	lowRange    *helper.NotificationState
	muchRange   *helper.NotificationState
}

type notificationState struct {
	notifType        notificationtype.Type
	inRange          bool
	lastNotification time.Time
	counter          int
}

func NewNotification(enforcer *casbin.Enforcer) *Notification {
	return &Notification{
		Notification: notification.NewNotification(enforcer),
		normalRange:  helper.NewNotificationState(notificationtype.HumidityNormal),
		lowRange:     helper.NewNotificationState(notificationtype.HumidityLow),
		muchRange:    helper.NewNotificationState(notificationtype.HumidityAbove),
	}
}

func (n *Notification) Get(humidityLevel float64, chatID int64) string {
	n.Init(humidityLevel, chatID)

	if n.HasNotification(n.normalRange) {
		return n.buildMessage(noSuffixEmoji)
	}

	if n.HasNotification(n.lowRange) {
		return n.buildMessage(suffixEmojiLow)
	}

	if n.HasNotification(n.muchRange) {
		return n.buildMessage(suffixEmojiMuch)
	}

	return ""
}

func (n *Notification) buildMessage(suffix string) string {
	return fmt.Sprintf("%s Humidity is: %.2f%%%s\n", prefixEmoji, n.Value, suffix)
}
