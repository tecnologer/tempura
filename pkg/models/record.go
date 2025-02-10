package models

import (
	"fmt"

	"gorm.io/gorm"
)

type Record struct {
	gorm.Model
	Temperature float64 `json:"temperature"`
	Humidity    float64 `json:"humidity"`
	FluidLevel  float64 `json:"fluid_level"`
	BatLevel    float64 `json:"bat_level"`
}

func (r Record) String() string {
	return fmt.Sprintf("Record{ID: %d, Temperature: %f, Humidity: %f, FluidLevel: %f, BatLevel: %f}",
		r.ID, r.Temperature, r.Humidity, r.FluidLevel, r.BatLevel,
	)
}
