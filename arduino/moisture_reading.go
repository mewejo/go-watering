package arduino

import (
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/mewejo/go-watering/world"
)

type MoistureReading struct {
	Time     time.Time
	Raw      int16
	Original world.MoistureLevel
	Sensor   MoistureSensor
}

func MakeMoistureReadingFromString(line string) (MoistureReading, error) {
	parts := strings.Split(line, ":")

	if parts[0] != "MS" {
		return MoistureReading{}, errors.New("line was not a moisture reading")

	}

	sensorId, err := strconv.Atoi(parts[1])

	if err != nil {
		return MoistureReading{}, err
	}

	rawValue, err := strconv.Atoi(parts[2])

	if err != nil {
		return MoistureReading{}, err
	}

	percentageValue, err := strconv.Atoi(parts[3])

	if err != nil {
		return MoistureReading{}, err
	}

	sensor, err := MoistureSensorFromId(sensorId)

	if err != nil {
		return MoistureReading{}, err
	}

	return MoistureReading{
		Time:   time.Now(),
		Sensor: sensor,
		Raw:    int16(rawValue),
		Original: world.MoistureLevel{
			Percentage: uint(percentageValue),
		},
	}, nil
}
