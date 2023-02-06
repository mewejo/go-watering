package model

import (
	"errors"
	"strconv"
	"strings"
	"time"
)

type MoistureReading struct {
	Time          time.Time `json:"recorded_at"`
	Raw           uint
	MoistureLevel MoistureLevel `json:"percentage"`
}

func (r *MoistureReading) CalculateMoistureLevelForSensor(sensor *MoistureSensor) {
	r.MoistureLevel = MakeMoistureLevel(
		sensor.mapRawReadingToPercentage(r.Raw),
	)
}

// Returns the reading, the sensor ID and an error
func MakeMoistureReadingFromString(line string) (MoistureReading, uint, error) {
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
		Raw:  uint(rawValue),
	}, uint(sensorId), nil
}
