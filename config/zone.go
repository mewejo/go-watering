package config

import (
	"errors"
	"log"
	"time"

	"github.com/mewejo/go-watering/arduino"
	"github.com/mewejo/go-watering/helpers"
	"github.com/mewejo/go-watering/world"
)

type Zone struct {
	Id                string // This will be user defined, used for API calls
	Name              string
	TargetMoisture    world.MoistureLevel
	MoistureSensors   []arduino.MoistureSensor
	WaterOutlets      []arduino.WaterOutlet
	MoisureReadings   []arduino.MoistureReading
	Watering          bool
	WateringChangedAt time.Time
}

func (z *Zone) SetWaterState(ard arduino.Arduino, state bool) {
	z.Watering = state
	z.WateringChangedAt = time.Now()

	z.enforceWateringState(ard)
}

func (z Zone) enforceWateringState(ard arduino.Arduino) {
	for _, outlet := range z.WaterOutlets {
		err := ard.SetWaterState(outlet, z.Watering)

		if err != nil {
			log.Fatal("could not set water state for zone")
		}
	}
}

func (z Zone) RequiresWatering() (bool, error) {
	moistureLevel, err := z.AverageMoistureLevel()

	if err != nil {
		return false, errors.New("could not get average moisture level for zone")
	}

	if moistureLevel.Percentage < z.TargetMoisture.Percentage {
		return true, nil
	}

	return false, nil
}

func (z Zone) AverageMoistureLevel() (world.MoistureLevel, error) {
	// Loop over the readings until we have one from each sensor
	readingsReversed := make([]arduino.MoistureReading, len(z.MoisureReadings))
	copy(readingsReversed, z.MoisureReadings)
	helpers.ReverseSlice(readingsReversed)

	sensorsFound := []arduino.MoistureSensor{}
	readings := []world.MoistureLevel{}

	for _, reading := range readingsReversed {
		if moistureSensorInSlice(reading.Sensor, sensorsFound) {
			continue
		}

		sensorsFound = append(sensorsFound, reading.Sensor)
		readings = append(readings, reading.Original)

		if len(sensorsFound) == len(z.MoistureSensors) {
			break
		}
	}

	if len(sensorsFound) != len(z.MoistureSensors) {
		return world.MoistureLevel{}, errors.New("incomplete data (sensors), cannot calculate moisture level")
	}

	if len(readings) != len(z.MoistureSensors) {
		return world.MoistureLevel{}, errors.New("incomplete data (readings), cannot calculate moisture level")
	}

	var totalPercentage uint

	for _, reading := range readings {
		totalPercentage += reading.Percentage
	}

	return world.MoistureLevel{
		Percentage: uint(totalPercentage / uint(len(readings))),
	}, nil
}

func (z *Zone) RecordMoistureReading(r arduino.MoistureReading) {
	z.MoisureReadings = append(z.MoisureReadings, r)
	limitMoistureReadings(&z.MoisureReadings, 100)
}

func limitMoistureReadings(s *[]arduino.MoistureReading, length int) {
	if len(*s) <= length {
		return
	}

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
