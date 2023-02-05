package model

import (
	"errors"
	"strconv"
	"strings"
	"time"
)

type MoistureReading struct {
	Time          time.Time
	raw           uint
	MoistureLevel MoistureLevel
}

func (r *MoistureReading) CalculateMoistureLevelForSensor(sensor MoistureSensor) {
	r.MoistureLevel = NewMoistureLevel(
		sensor.mapRawReadingToPercentage(r.raw),
	)
}

// Returns the reading, the sensor ID and an error
func NewMoistureReadingFromString(line string) (MoistureReading, uint, error) {
	// MS:1:700:44 # MS:ID:RAW:PERCENTAGE
	parts := strings.Split(line, ":")

	if parts[0] != "MS" {
		return MoistureReading{}, 0, errors.New("line was not a moisture reading")
	}

	sensorId, err := strconv.Atoi(parts[1])

	if err != nil {
		return MoistureReading{}, 0, err
	}

	rawValue, err := strconv.Atoi(parts[2])

	if err != nil {
		return MoistureReading{}, 0, err
	}

	return MoistureReading{
		Time: time.Now(),
		raw:  uint(rawValue),
	}, uint(sensorId), nil
}
