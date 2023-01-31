package config

import (
	"github.com/mewejo/go-watering/arduino"
	"github.com/mewejo/go-watering/helpers"
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

func (z Zone) AverageMoistureLevel() (world.MoistureLevel, error) {
	// Loop over the readings until we have one from each sensor
	readingsReversed := make([]arduino.MoistureReading, len(z.MoisureReadings))
	copy(readingsReversed, z.MoisureReadings)
	helpers.ReverseSlice(readingsReversed)

	sensorsFound := []arduino.MoistureSensor{}

	for _, reading := range readingsReversed {
		if moistureSensorInSlice(reading.Sensor, sensorsFound) {
			continue
		}

		sensorsFound = append(sensorsFound, reading.Sensor)

		if len(sensorsFound) == len(z.MoistureSensors) {
			break
		}
	}
}

func (z *Zone) RecordMoistureReading(r arduino.MoistureReading) {
	z.MoisureReadings = append(z.MoisureReadings, r)
	limitMoistureReadings(&z.MoisureReadings, 100)
}

func limitMoistureReadings(s *[]arduino.MoistureReading, length int) {
	*s = (*s)[len(*s)-length:]
}

func moistureSensorInSlice(s arduino.MoistureSensor, sensors []arduino.MoistureSensor) bool {
	for _, v := range sensors {
		if v == s {
			return true
		}
	}

	return false
}
