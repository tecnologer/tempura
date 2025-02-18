package notificationtype

//go:generate enumer -type=Type -json -sql -transform=snake
type Type byte //nolint:recvcheck

const (
	All Type = iota
	BatteryNormal
	BatteryLow
	BatteryCharging
	BatteryCharged
	BatteryDanger
	HumidityNormal
	HumidityLow
	HumidityAbove
	TemperatureNormal
	TemperatureLow
	TemperatureHigh
	FluidTankEmpty
	FluidTankFull
)

func (Type) TableName() string {
	return "notification_types"
}
