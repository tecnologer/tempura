package dao

import (
	"fmt"
	"strings"
	"time"

	"github.com/tecnologer/tempura/pkg/casbin"
	"github.com/tecnologer/tempura/pkg/dao/db"
	"github.com/tecnologer/tempura/pkg/dao/notifications/battery"
	"github.com/tecnologer/tempura/pkg/dao/notifications/fluid"
	"github.com/tecnologer/tempura/pkg/dao/notifications/humidity"
	"github.com/tecnologer/tempura/pkg/dao/notifications/temperature"
	"github.com/tecnologer/tempura/pkg/models"
	"github.com/tecnologer/tempura/pkg/telegram"
)

const (
	MinBatLevelNotification        = 30.0
	DangerBatLevelNotification     = 20.0
	BatResendNotificationTime      = 30 * time.Minute
	BatLevelDangerNotificationTime = 5 * time.Minute

	MinHumidityNotification       = 60.0
	MaxHumidityNotification       = 90.0
	MinHumidityResendNotification = 3 * time.Hour
	MaxHumidityResendNotification = 3 * time.Hour

	MinTemperatureNotification       = 5.0
	MaxTemperatureNotification       = 40.0
	MinTemperatureResendNotification = 1 * time.Hour
	MaxTemperatureResendNotification = 30 * time.Minute

	FluidLevelNotification  = 100.0
	FloodResendNotification = 24 * time.Hour
)

type notificationManager struct {
	battery     *battery.Notification
	temperature *temperature.Notification
	humidity    *humidity.Notification
	fluid       *fluid.Notification
}

type Notification struct {
	cnn     *db.Connection
	bot     *telegram.Bot
	manager *notificationManager
}

func NewNotification(cnn *db.Connection, token string) (*Notification, error) {
	bot, err := telegram.NewBot(token, false)
	if err != nil {
		return nil, fmt.Errorf("creating bot: %w", err)
	}

	enforcer, err := casbin.NewEnforcerWithDb(cnn.DB)
	if err != nil {
		return nil, fmt.Errorf("creating enforcer: %w", err)
	}

	return &Notification{
		cnn: cnn,
		bot: bot,
		manager: &notificationManager{
			battery:     battery.NewNotification(enforcer),
			temperature: temperature.NewNotification(enforcer),
			humidity:    humidity.NewNotification(enforcer),
			fluid:       fluid.NewNotification(enforcer),
		},
	}, nil
}

func (n *Notification) NotifyNewRecord(record *models.Record) error {
	var users []*models.User

	tx := n.cnn.DB.Find(&users)
	if tx.Error != nil {
		return fmt.Errorf("getting users: %w", tx.Error)
	}

	errs := make([]error, 0, 1)

	for _, user := range users {
		if user.Role != "admin" || user.TelegramChatID == 0 {
			continue
		}

		message := n.buildMessage(record, user.TelegramChatID)

		message = fmt.Sprintf("Hi, %s, here is the latest update from your device:\n\n%s", user.Username, message)

		if err := n.bot.SendMessage(user.TelegramChatID, message); err != nil {
			errs = append(errs, fmt.Errorf("sending message to %s: %w", user.Username, err))
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("sending message errors: %v", errs)
	}

	return nil
}

func (n *Notification) buildMessage(record *models.Record, chatID int64) string {
	var notification strings.Builder

	notification.WriteString(n.manager.battery.Get(record.BatLevel, chatID))
	notification.WriteString(n.manager.humidity.Get(record.Humidity, chatID))
	notification.WriteString(n.manager.temperature.Get(record.Temperature, chatID))
	notification.WriteString(n.manager.fluid.Get(record.FluidLevel, chatID))

	return notification.String()
}
