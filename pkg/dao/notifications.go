package dao

import (
	"fmt"
	"strings"
	"time"

	"github.com/tecnologer/tempura/pkg/dao/db"
	"github.com/tecnologer/tempura/pkg/models"
	"github.com/tecnologer/tempura/pkg/telegram"
	"github.com/tecnologer/tempura/pkg/utils/log"
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
	lastBatteryLevelNotified time.Time
	lastHumidityNotified     time.Time
	lastTemperatureNotified  time.Time
	lastFluidLevelNotified   time.Time
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

	return &Notification{
		cnn:     cnn,
		bot:     bot,
		manager: &notificationManager{},
	}, nil
}

func (n *Notification) NotifyNewRecord(record *models.Record) error {
	message := n.buildMessage(record)
	if message == "" {
		return nil
	}

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

func (n *Notification) buildMessage(record *models.Record) string {
	var notification strings.Builder

	notification.WriteString(n.buildHumidityNotification(record))
	notification.WriteString(n.buildTemperatureNotification(record))
	notification.WriteString(n.buildFluidLevelNotification(record))
	notification.WriteString(n.buildBatteryLevelNotification(record))

	return notification.String()
}

func (n *Notification) buildBatteryLevelNotification(record *models.Record) string {
	isBatteryLevelInRange := record.BatLevel > MinBatLevelNotification
	isBelowMinBatteryLevel := record.BatLevel <= MinBatLevelNotification
	isBatteryLevelDanger := record.BatLevel <= DangerBatLevelNotification

	// Check if no notification is needed when battery level is in range and no previous notifications
	if n.isNotificationBatteryLevelEmpty(isBatteryLevelInRange, isBelowMinBatteryLevel, isBatteryLevelDanger) {
		return ""
	}

	prefixEmoji, suffixEmoji := n.batteryLevelEmoji(isBelowMinBatteryLevel, isBatteryLevelDanger)

	// Handle returning to normal range
	if isBatteryLevelInRange {
		n.manager.lastBatteryLevelNotified = time.Time{}

		return fmt.Sprintf("%s Battery level back to normal: %.2f%%\n", prefixEmoji, record.BatLevel)
	}

	// Update last notified values
	n.manager.lastBatteryLevelNotified = time.Now()

	// Return the formatted message
	return fmt.Sprintf("%s Battery level is: %.2f%%%s\n", prefixEmoji, record.BatLevel, suffixEmoji)
}

func (n *Notification) buildHumidityNotification(record *models.Record) string {
	isHumidityInRange := record.Humidity > MinHumidityNotification && record.Humidity < MaxHumidityNotification
	isBelowMinHumidity := record.Humidity <= MinHumidityNotification
	isAboveMaxHumidity := record.Humidity >= MaxHumidityNotification

	// Check if no notification is needed when humidity is in range and no previous notifications
	if n.isNotificationHumidityEmpty(isHumidityInRange, isBelowMinHumidity, isAboveMaxHumidity) {
		return ""
	}

	prefixEmoji, suffixEmoji := n.humiditySuffixEmoji(isBelowMinHumidity, isAboveMaxHumidity)

	// Handle returning to normal range
	if isHumidityInRange {
		n.manager.lastHumidityNotified = time.Time{}

		return fmt.Sprintf("%s Humidity is back to normal: %.2f%%\n", prefixEmoji, record.Humidity)
	}

	// Update last notified values
	n.manager.lastHumidityNotified = time.Now()

	// Return the formatted message
	return fmt.Sprintf("%s Humidity is: %.2f%%%s\n", prefixEmoji, record.Humidity, suffixEmoji)
}

func (n *Notification) buildTemperatureNotification(record *models.Record) string {
	isTemperatureInRange := record.Temperature > MinTemperatureNotification && record.Temperature < MaxTemperatureNotification
	isBelowMinTemperature := record.Temperature <= MinTemperatureNotification
	isAboveMaxTemperature := record.Temperature >= MaxTemperatureNotification

	// if the temperature is in range, and there is no previous notifications, do not send a notification
	if n.isNotificationTemperatureEmpty(isTemperatureInRange, isBelowMinTemperature, isAboveMaxTemperature) {
		return ""
	}

	prefixEmoji, suffixEmoji := n.temperatureEmoji(isBelowMinTemperature, isAboveMaxTemperature)

	// Handle returning to normal range
	if isTemperatureInRange {
		n.manager.lastTemperatureNotified = time.Time{}

		return fmt.Sprintf("%s Temperature is back to normal: %.2f°C\n", prefixEmoji, record.Temperature)
	}

	n.manager.lastTemperatureNotified = time.Now()

	return fmt.Sprintf("%s Temperature is: %.2f°C%s\n", prefixEmoji, record.Temperature, suffixEmoji)
}

func (n *Notification) buildFluidLevelNotification(record *models.Record) string {
	isFluidTankFull := record.FluidLevel >= FluidLevelNotification

	if !isFluidTankFull && n.manager.lastFluidLevelNotified.IsZero() {
		return ""
	}

	// Handle returning to normal range
	if !isFluidTankFull {
		n.manager.lastFluidLevelNotified = time.Time{}

		return n.fluidLevelEmoji() + " Fluid tank is empty\n"
	}

	if time.Since(n.manager.lastFluidLevelNotified) <= FloodResendNotification {
		return ""
	}

	n.manager.lastFluidLevelNotified = time.Now()

	return n.fluidLevelEmoji() + " Fluid tank is full\n"
}

func (n *Notification) batteryLevelEmoji(isBelowMinBatLevel, isBatteryLevelDanger bool) (string, string) {
	var suffixEmoji string

	const prefixEmoji = "🔋"

	// Adjust emoji based on battery level
	if isBatteryLevelDanger {
		suffixEmoji = " ⚡️☠️☠️"
	} else if isBelowMinBatLevel {
		suffixEmoji = " ⚡️⚠️⚠️"
	}

	return prefixEmoji, suffixEmoji
}

func (n *Notification) humiditySuffixEmoji(isBelowMinHumidity, isAboveMaxHumidity bool) (string, string) {
	var suffixEmoji string

	const prefixEmoji = "💧"

	// Adjust emoji based on humidity level
	if isBelowMinHumidity {
		suffixEmoji = " 🏜️️"
	} else if isAboveMaxHumidity {
		suffixEmoji = " 💦"
	}

	return prefixEmoji, suffixEmoji
}

func (n *Notification) temperatureEmoji(isBelowMinTemperature, isAboveMaxTemperature bool) (string, string) {
	var suffixEmoji string

	const prefixEmoji = "🌡"

	// Adjust emoji based on temperature level
	if isBelowMinTemperature {
		suffixEmoji = "❄️️"
	} else if isAboveMaxTemperature {
		suffixEmoji = "🔥"
	}

	return prefixEmoji, suffixEmoji
}

func (n *Notification) fluidLevelEmoji() string {
	return "🛢️"
}

// isNotificationBatteryLevelEmpty return true if no notification is needed
// if the battery level is in range, and there is no previous notifications, do not send a notification
// if it's too soon to resend notification for out of range
func (n *Notification) isNotificationBatteryLevelEmpty(isBatteryLevelInRange, isBelowMinBatLevel, isBatteryLevelDanger bool) bool {
	log.Debugf(
		"isBatteryLevelInRange: %v, isBelowMinBatLevel: %v, isBatteryLevelDanger: %v, lastBatteryLevelNotified: %v (want: %v)",
		isBatteryLevelInRange,
		isBelowMinBatLevel,
		isBatteryLevelDanger,
		time.Since(n.manager.lastBatteryLevelNotified),
		BatResendNotificationTime,
	)

	// Check if no notification is needed when battery level is in range and no previous notifications
	return (isBatteryLevelInRange && n.manager.lastBatteryLevelNotified.IsZero()) ||
		// Check if it's too soon to resend notification for out of range
		(time.Since(n.manager.lastBatteryLevelNotified) <= BatResendNotificationTime && isBelowMinBatLevel) ||
		(time.Since(n.manager.lastBatteryLevelNotified) <= BatLevelDangerNotificationTime && isBatteryLevelDanger)
}

// isNotificationHumidityEmpty return true if no notification is needed
// if the humidity is in range, and there is no previous notifications, do not send a notification
// if it's too soon to resend notification for out of range
func (n *Notification) isNotificationHumidityEmpty(isHumidityInRange, isBelowMinHumidity, isAboveMaxHumidity bool) bool {
	log.Debugf(
		"isHumidityInRange: %v, isBelowMinHumidity: %v, isAboveMaxHumidity: %v, minHumidityNotified: %v (want: %v), maxHumidityNotified: %v (want: %v)",
		isHumidityInRange,
		isBelowMinHumidity,
		isAboveMaxHumidity,
		time.Since(n.manager.lastHumidityNotified),
		MinHumidityResendNotification,
		time.Since(n.manager.lastHumidityNotified),
		MaxHumidityResendNotification,
	)

	// if the humidity is in range, and there is no previous notifications, do not send a notification
	return (isHumidityInRange && n.manager.lastHumidityNotified.IsZero()) ||
		// if it's too soon to resend notification for out of range
		(time.Since(n.manager.lastHumidityNotified) <= MinHumidityResendNotification && isBelowMinHumidity) ||
		(time.Since(n.manager.lastHumidityNotified) <= MaxHumidityResendNotification && isAboveMaxHumidity)
}

// isNotificationTemperatureEmpty return true if no notification is needed
// if the temperature is in range, and there is no previous notifications, do not send a notification
// if it's too soon to resend notification for out of range
func (n *Notification) isNotificationTemperatureEmpty(isTemperatureInRange, isBelowMinTemperature, isAboveMaxTemperature bool) bool {
	//nolint: lll
	log.Debugf(
		"isTemperatureInRange: %v, isBelowMinTemperature: %v, isAboveMaxTemperature: %v, minTemperatureNotified: %v (want: %v), maxTemperatureNotified: %v (want: %v)",
		isTemperatureInRange,
		isBelowMinTemperature,
		isAboveMaxTemperature,
		time.Since(n.manager.lastTemperatureNotified),
		MinTemperatureResendNotification,
		time.Since(n.manager.lastTemperatureNotified),
		MaxTemperatureResendNotification,
	)

	// if the temperature is in range, and there is no previous notifications, do not send a notification
	return (isTemperatureInRange && n.manager.lastTemperatureNotified.IsZero()) ||
		// if it's too soon to resend notification for out of range
		(time.Since(n.manager.lastTemperatureNotified) <= MinTemperatureResendNotification && isBelowMinTemperature) ||
		(time.Since(n.manager.lastTemperatureNotified) <= MaxTemperatureResendNotification && isAboveMaxTemperature)
}
