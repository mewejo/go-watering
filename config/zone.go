package config

import (
	"github.com/mewejo/go-watering/arduino"
	"github.com/mewejo/go-watering/world"
)

type Zone struct {
	Id              string // This will be user defined, used for API calls
	Name            string
	MoistureSensors []arduino.MoistureSensor
	MoisureReadings []arduino.MoistureReading
	TargetMoisture  world.MoistureLevel
}
