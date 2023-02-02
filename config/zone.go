package config

import (
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/mewejo/go-watering/arduino"
	"github.com/mewejo/go-watering/helpers"
	"github.com/mewejo/go-watering/homeassistant"
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
	ForcedWatering    bool
}

func (z Zone) GetHomeAssistantHumidifierConfiguration() homeassistant.HumidifierConfiguration {
	c := homeassistant.NewZoneHumidifierConfiguration()
	c.ObjectId = z.GetHomeAssistantObjectId()
	c.UniqueId = z.GetHomeAssistantObjectId()
	c.Name = z.Name + " Hygrostat"
	c.StateTopic = z.GetHomeAssistantStateTopic()
	c.CommandTopic = z.GetHomeAssistantCommandTopic()
	c.TargetHumidityTopic = z.GetHomeAssistantTargetHumidityTopic()
	c.TargetHumidityStateTopic = z.GetHomeAssistantTargetStateHumidityTopic()
	c.AvailabilityTopic = z.GetHomeAssistantAvailabilityTopic()
	c.ModeStateTopic = z.GetHomeAssistantModeStateTopic()
	c.ModeCommandTopic = z.GetHomeAssistantModeCommandTopic()
	c.Device = homeassistant.NewDeviceDetails()

	return c
}

func (z Zone) GetHomeAssistantMoistureSensorConfiguration() homeassistant.MoistureSensorConfiguration {
	c := homeassistant.NewMoistureSensorConfiguration()
	c.Name = z.Name + " Moisture"
	c.ObjectId = z.GetHomeAssistantObjectId()
	c.UniqueId = z.GetHomeAssistantObjectId()
	c.StateTopic = z.GetHomeAssistantStateTopic()
	c.AvailabilityTopic = z.GetHomeAssistantAvailabilityTopic()
	c.Device = homeassistant.NewDeviceDetails()

	return c
}

func (z Zone) GetHomeAssistantHumidifierBaseTopic() string {
	return fmt.Sprintf(
		"%v/humidifier/%v",
		os.Getenv("HOME_ASSISTANT_DISCOVERY_PREFIX"),
		z.GetHomeAssistantObjectId(),
	)
}

func (z Zone) GetHomeAssistantMoistureSensorBaseTopic() string {
	return fmt.Sprintf(
		"%v/sensor/%v/average-moisture",
		os.Getenv("HOME_ASSISTANT_DISCOVERY_PREFIX"),
		z.GetHomeAssistantObjectId(),
	)
}

func (z Zone) GetHomeAssistantModeStateTopic() string {
	return fmt.Sprintf("%v/mode_state", z.GetHomeAssistantHumidifierBaseTopic())
}

func (z Zone) GetHomeAssistantModeCommandTopic() string {
	return fmt.Sprintf("%v/mode_command", z.GetHomeAssistantHumidifierBaseTopic())
}

func (z Zone) GetHomeAssistantAvailabilityTopic() string {
	return fmt.Sprintf("%v/availability", z.GetHomeAssistantHumidifierBaseTopic())
}

func (z Zone) GetHomeAssistantCommandTopic() string {
	return fmt.Sprintf("%v/command", z.GetHomeAssistantHumidifierBaseTopic())
}

func (z Zone) GetHomeAssistantTargetHumidityTopic() string {
	return fmt.Sprintf("%v/target", z.GetHomeAssistantHumidifierBaseTopic())
}

func (z Zone) GetHomeAssistantTargetStateHumidityTopic() string {
	return fmt.Sprintf("%v/humidity_state", z.GetHomeAssistantHumidifierBaseTopic())
}

func (z Zone) GetHomeAssistantStateTopic() string {
	return fmt.Sprintf("%v/state", z.GetHomeAssistantHumidifierBaseTopic())
}

func (z Zone) GetHomeAssistantObjectId() string {
	return os.Getenv("HOME_ASSISTANT_OBJECT_ID_PREFIX") + z.Id
}

func (z *Zone) SetForcedWateringState(ard arduino.Arduino, state bool) {
	if z.ForcedWatering == state {
		return
	}

	z.ForcedWatering = state
	z.WateringChangedAt = time.Now()

	z.EnforceWateringState(ard)
}

func (z *Zone) SetWaterState(ard arduino.Arduino, state bool) {
	if z.Watering == state {
		return
	}

	z.Watering = state
	z.WateringChangedAt = time.Now()

	z.EnforceWateringState(ard)
}

func (z Zone) EnforceWateringState(ard arduino.Arduino) {
	for _, outlet := range z.WaterOutlets {
		err := ard.SetWaterState(outlet, z.Watering || z.ForcedWatering)

		if err != nil {
			log.Fatal("could not set water state for zone")
		}
	}
}

func (z Zone) ShouldStartWatering() (bool, error) {
	moistureLevel, err := z.AverageMoistureLevel()

	if err != nil {
		return false, errors.New("could not get average moisture level for zone")
	}

	// Already watering... don't need to start
	if z.Watering {
		return false, nil
	}

	if moistureLevel.Percentage < z.TargetMoisture.HysteresisOnLevel().Percentage {
		return true, nil
	}

	return false, nil
}

func (z Zone) ShouldStopWatering() (bool, error) {
	moistureLevel, err := z.AverageMoistureLevel()

	if err != nil {
		return false, errors.New("could not get average moisture level for zone")
	}

	// Not watering... don't need to stop
	if !z.Watering {
		return false, nil
	}

	if moistureLevel.Percentage > z.TargetMoisture.HysteresisOffLevel().Percentage {
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

	if len(readings) < 1 {
		return world.MoistureLevel{}, errors.New("no readings to make average from")
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
