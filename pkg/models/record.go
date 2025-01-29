package models

import "gorm.io/gorm"

type Record struct {
	gorm.Model
	Temperature float64 `json:"temperature"`
	Humidity    float64 `json:"humidity"`
	FluidLevel  float64 `json:"fluid_level"`
}
