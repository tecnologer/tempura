package casbin_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	tcasbin "github.com/tecnologer/tempura/pkg/casbin"
	"github.com/tecnologer/tempura/pkg/models/notification"
)

func TestEnforcer_Enforce(t *testing.T) { //nolint:funlen
	t.Parallel()

	enforcer, err := tcasbin.NewEnforcer(nil)
	require.NoError(t, err)
	require.NotNil(t, enforcer)

	_, err = enforcer.AddPolicies([][]string{
		{"12345", notificationtype.TypeBatteryLow.String(), "10", "20", "100"},
	})
	require.NoError(t, err)

	now := time.Now()

	tests := []struct {
		name     string
		requests []*tcasbin.Request
		want     bool
		wantErr  bool
	}{
		{
			name: "12345_is_authorized",
			requests: []*tcasbin.Request{
				{
					ChatID: 12345,
					Type:   notificationtype.TypeBatteryLow,
					Value:  15,
					LastSent: time.Date(
						now.Year(),
						now.Month(),
						now.Day(),
						now.Hour(),
						now.Minute(),
						now.Second(),
						now.Nanosecond()-200,
						now.Location(),
					),
				},
			},
			want: true,
		},
		{
			name: "654321_is_not_authorized",
			requests: []*tcasbin.Request{
				{
					ChatID: 654321,
					Type:   notificationtype.TypeBatteryLow,
					Value:  15,
					LastSent: time.Date(
						now.Year(),
						now.Month(),
						now.Day(),
						now.Hour(),
						now.Minute(),
						now.Second(),
						now.Nanosecond()-200,
						now.Location(),
					),
				},
			},
			want: false,
		},
		{
			name: "alice_is_not_authorized_by_value",
			requests: []*tcasbin.Request{
				{
					ChatID: 12345,
					Type:   notificationtype.TypeBatteryLow,
					Value:  25,
					LastSent: time.Date(
						now.Year(),
						now.Month(),
						now.Day(),
						now.Hour(),
						now.Minute(),
						now.Second(),
						now.Nanosecond()-200,
						now.Location(),
					),
				},
			},
		},
		{
			name: "alice_is_not_authorized_by_delay",
			requests: []*tcasbin.Request{
				{
					ChatID: 12345,
					Type:   notificationtype.TypeBatteryLow,
					Value:  25,
					LastSent: time.Date(
						now.Year(),
						now.Month(),
						now.Day(),
						now.Hour(),
						now.Minute(),
						now.Second(),
						now.Nanosecond()-200,
						now.Location(),
					),
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			isAuthorized, err := enforcer.Enforce(test.requests...)
			require.NoError(t, err)
			assert.Equal(t, test.want, isAuthorized)
		})
	}
}
