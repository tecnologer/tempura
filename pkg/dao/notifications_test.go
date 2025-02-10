package dao //nolint: testpackage

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/tecnologer/tempura/pkg/models"
)

func TestNotification_buildTemperatureNotification(t *testing.T) { //nolint:funlen
	t.Parallel()

	now := time.Now()

	tests := []struct {
		name                 string
		record               *models.Record
		lastNotificationTime time.Time
		want                 string
	}{
		{
			name: "normal",
			record: &models.Record{
				Temperature: 25.0,
			},
			want: "",
		},
		{
			name: "min",
			record: &models.Record{
				Temperature: 5.0,
			},
			want: "🌡 ❄️❄️❄️ Temperature is: 5.00°C\n",
		},
		{
			name: "max",
			record: &models.Record{
				Temperature: 40.0,
			},
			want: "🌡 🔥🔥🔥 Temperature is: 40.00°C\n",
		},
		{
			name: "resend_less_than_resend_time",
			record: &models.Record{
				Temperature: 5.0,
			},
			lastNotificationTime: time.Date(
				now.Year(),
				now.Month(),
				now.Day(),
				now.Hour(),
				now.Minute()-1,
				now.Second(),
				now.Nanosecond(),
				now.Location(),
			),
		},
		{
			name: "resend_greater_than_resend_time",
			record: &models.Record{
				Temperature: 5.0,
			},
			lastNotificationTime: time.Date(
				now.Year(),
				now.Month(),
				now.Day(),
				now.Hour()-2,
				now.Minute(),
				now.Second(),
				now.Nanosecond(),
				now.Location(),
			),
			want: "🌡 ❄️❄️❄️ Temperature is: 5.00°C\n",
		},
		{
			name: "normal_back_to_normal",
			record: &models.Record{
				Temperature: 25.0,
			},
			lastNotificationTime: time.Now(),
			want:                 "🌡 Temperature is back to normal: 25.00°C\n",
		},
		{
			name: "normal_after_back_to_normal",
			record: &models.Record{
				Temperature: 25.0,
			},
			want: "",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			notification := &Notification{
				manager: &notificationManager{
					lastTemperatureNotified: test.lastNotificationTime,
				},
			}

			got := notification.buildTemperatureNotification(test.record)
			assert.Equal(t, test.want, got)
		})
	}
}

func TestNotification_buildHumidityNotification(t *testing.T) { //nolint:funlen
	t.Parallel()

	now := time.Now()

	tests := []struct {
		name                 string
		record               *models.Record
		lastNotificationTime time.Time
		want                 string
	}{
		{
			name: "normal",
			record: &models.Record{
				Humidity: 75.0,
			},
			want: "",
		},
		{
			name: "min",
			record: &models.Record{
				Humidity: 60.0,
			},
			want: "💧 🏜️🏜️🏜️ Humidity is: 60.00%\n",
		},
		{
			name: "max",
			record: &models.Record{
				Humidity: 90.0,
			},
			want: "💧 💦💦💦 Humidity is: 90.00%\n",
		},
		{
			name: "resend_less_than_resend_time",
			record: &models.Record{
				Humidity: 10.0,
			},
			lastNotificationTime: time.Date(
				now.Year(),
				now.Month(),
				now.Day(),
				now.Hour(),
				now.Minute()-1,
				now.Second(),
				now.Nanosecond(),
				now.Location(),
			),
		},
		{
			name: "resend_greater_than_resend_time",
			record: &models.Record{
				Humidity: 10.0,
			},
			lastNotificationTime: time.Date(
				now.Year(),
				now.Month(),
				now.Day(),
				now.Hour()-6,
				now.Minute(),
				now.Second(),
				now.Nanosecond(),
				now.Location(),
			),
			want: "💧 🏜️🏜️🏜️ Humidity is: 10.00%\n",
		},
		{
			name: "normal_back_to_normal",
			record: &models.Record{
				Humidity: 65.0,
			},
			lastNotificationTime: time.Now(),
			want:                 "💧 Humidity is back to normal: 65.00%\n",
		},
		{
			name: "normal_after_back_to_normal",
			record: &models.Record{
				Humidity: 65.0,
			},
			want: "",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			notification := &Notification{
				manager: &notificationManager{
					lastHumidityNotified: test.lastNotificationTime,
				},
			}

			got := notification.buildHumidityNotification(test.record)
			assert.Equal(t, test.want, got)
		})
	}
}

func TestNotification_buildFluidLevelNotification(t *testing.T) { //nolint:funlen
	t.Parallel()

	now := time.Now()

	tests := []struct {
		name                 string
		record               *models.Record
		lastNotificationTime time.Time
		want                 string
	}{
		{
			name: "normal",
			record: &models.Record{
				FluidLevel: 0.0,
			},
			want: "",
		},
		{
			name: "full",
			record: &models.Record{
				FluidLevel: 100.0,
			},
			want: "🛢️ Fluid tank is full\n",
		},
		{
			name: "resend_less_than_resend_time",
			record: &models.Record{
				FluidLevel: 100.0,
			},
			lastNotificationTime: time.Date(
				now.Year(),
				now.Month(),
				now.Day(),
				now.Hour(),
				now.Minute()-1,
				now.Second(),
				now.Nanosecond(),
				now.Location(),
			),
		},
		{
			name: "resend_greater_than_resend_time",
			record: &models.Record{
				FluidLevel: 100.0,
			},
			lastNotificationTime: time.Date(
				now.Year(),
				now.Month(),
				now.Day()-1,
				now.Hour()-1,
				now.Minute(),
				now.Second(),
				now.Nanosecond(),
				now.Location(),
			),
			want: "🛢️ Fluid tank is full\n",
		},
		{
			name: "normal_back_to_normal",
			record: &models.Record{
				FluidLevel: 0.0,
			},
			lastNotificationTime: time.Now(),
			want:                 "🛢️ Fluid tank is empty\n",
		},
		{
			name: "normal_after_back_to_normal",
			record: &models.Record{
				FluidLevel: 0.0,
			},
			want: "",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			notification := &Notification{
				manager: &notificationManager{
					lastFluidLevelNotified: test.lastNotificationTime,
				},
			}

			got := notification.buildFluidLevelNotification(test.record)
			assert.Equal(t, test.want, got)
		})
	}
}
