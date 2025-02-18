package battery

import (
	"fmt"
	"math"

	"github.com/tecnologer/tempura/pkg/casbin"
	"github.com/tecnologer/tempura/pkg/dao/notifications/helper"
	"github.com/tecnologer/tempura/pkg/dao/notifications/notification"
	"github.com/tecnologer/tempura/pkg/models/notificationtype"
)

const (
	prefixEmoji         = "🔋"
	suffixEmojiWarning  = " ⚠️"
	suffixEmojiDanger   = " ☠️"
	suffixEmojiCharged  = " ✅"
	suffixEmojiCharging = " ⚡️"
	chargingCountMin    = 4 // Minimum of states changes to consider charging
	noSuffixEmoji       = ""
)

type Notification struct {
	*notification.Notification
	normalRange   *helper.NotificationState
	lowRange      *helper.NotificationState
	chargingRange *helper.NotificationState
	chargedRange  *helper.NotificationState
	dangerRange   *helper.NotificationState
}

func NewNotification(enforcer *casbin.Enforcer) *Notification {
	return &Notification{
		Notification:  notification.NewNotification(enforcer),
		normalRange:   helper.NewNotificationState(notificationtype.BatteryNormal),
		lowRange:      helper.NewNotificationState(notificationtype.BatteryLow),
		chargingRange: helper.NewNotificationState(notificationtype.BatteryCharging),
		chargedRange:  helper.NewNotificationState(notificationtype.BatteryCharged),
		dangerRange:   helper.NewNotificationState(notificationtype.BatteryDanger),
	}
}

func (n *Notification) Get(batteryLevel float64, chatID int64) string {
	n.Init(batteryLevel, chatID, n.calculateIsCharging)

	if n.HasNotification(n.dangerRange) {
		return n.buildMessage(batteryLevel, suffixEmojiDanger)
	}

	if n.HasNotification(n.normalRange) {
		return n.buildMessage(batteryLevel, noSuffixEmoji)
	}

	if n.HasNotification(n.lowRange) {
		return n.buildMessage(batteryLevel, suffixEmojiWarning)
	}

	if n.chargingRange.InRange && n.HasNotification(n.chargingRange) {
		return n.buildMessage(batteryLevel, suffixEmojiCharging)
	}

	if n.HasNotification(n.chargedRange) {
		return n.buildMessage(batteryLevel, suffixEmojiCharged)
	}

	return ""
}

func (n *Notification) buildMessage(value float64, suffix string) string {
	return fmt.Sprintf("%s Battery level is at %.0f%%%s\n", prefixEmoji, value, suffix)
}

func (n *Notification) calculateIsCharging(batteryLevel float64, _ int64) {
	if n.Value < batteryLevel {
		n.chargingRange.Counter++
	} else {
		n.chargingRange.Counter = int(math.Max(0, float64(n.chargingRange.Counter-1)))
	}

	n.chargingRange.InRange = batteryLevel > 0 && n.chargingRange.Counter > chargingCountMin
}
