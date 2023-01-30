package config

import (
	"github.com/mewejo/go-watering/arduino"
	"github.com/mewejo/go-watering/world"
)

type Zone struct {
	Id              string // This will be user defined, used for API calls
	Name            string
	TargetMoisture  world.MoistureLevel
	MoistureSensors []arduino.MoistureSensor
	WaterOutlets    []arduino.WaterOutlet
	MoisureReadings []arduino.MoistureReading
}
