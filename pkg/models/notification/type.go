package notification

//go:generate enumer -type=Type -json -sql -transform=snake -trimprefix=Type
type Type byte //nolint:recvcheck

const (
	TypeAll Type = iota
	TypeBatteryLow
	TypeBatteryCharging
	TypeBatteryCharged
	TypeHumidityBelow
	TypeHumidityAbove
	TypeTemperatureBelow
	TypeTemperatureAbove
	TypeFluidTankEmpty
	TypeFluidTankFull
)

func (Type) TableName() string {
	return "notification_types"
}
