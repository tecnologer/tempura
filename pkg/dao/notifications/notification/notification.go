package notification

import (
	"time"

	"github.com/tecnologer/tempura/pkg/casbin"
	"github.com/tecnologer/tempura/pkg/dao/notifications/helper"
	"github.com/tecnologer/tempura/pkg/utils/log"
)

type InitFunc func(value float64, chatID int64)

type Notification struct {
	ChatID   int64
	Value    float64
	Enforcer *casbin.Enforcer
}

func NewNotification(enforcer *casbin.Enforcer) *Notification {
	return &Notification{
		Enforcer: enforcer,
	}
}

func (n *Notification) HasNotification(notifState *helper.NotificationState) bool {
	requestCharging := &casbin.Request{
		ChatID:   n.ChatID,
		Type:     notifState.NotifType,
		Value:    n.Value,
		LastSent: notifState.LastNotification,
	}

	hasNotif, err := n.Enforcer.Enforce(requestCharging)
	if err != nil {
		log.Errorf("enforcing %s notification: %v", notifState.NotifType, err)
		return false
	}

	if !hasNotif {
		return false
	}

	notifState.Counter = 0
	notifState.LastNotification = time.Now()

	return true
}

func (n *Notification) Init(value float64, chatID int64, callbacks ...InitFunc) {
	n.Value = value
	n.ChatID = chatID

	for _, callback := range callbacks {
		callback(value, chatID)
	}
}
