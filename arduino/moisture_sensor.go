package arduino

import (
	"errors"
	"fmt"
	"os"

	"github.com/mewejo/go-watering/homeassistant"
)

type MoistureSensor int

const (
	MOISTURE_SENSOR_1 MoistureSensor = iota
	MOISTURE_SENSOR_2
	MOISTURE_SENSOR_3
	MOISTURE_SENSOR_4
	MOISTURE_SENSOR_5
	MOISTURE_SENSOR_6
)

func (ms MoistureSensor) GetHomeAssistantMoistureSensorConfiguration() homeassistant.MoistureSensorConfiguration {
	c := homeassistant.NewMoistureSensorConfiguration()
	c.Name = fmt.Sprintf("Moisture Sensor #%v", ms.GetId())
	c.ObjectId = ms.GetHomeAssistantObjectId()
	c.UniqueId = ms.GetHomeAssistantObjectId()
	c.StateTopic = ms.GetHomeAssistantStateTopic()
	c.AvailabilityTopic = ms.GetHomeAssistantAvailabilityTopic()
	c.Device = homeassistant.NewDeviceDetails()

	return c
}

func (ms MoistureSensor) GetHomeAssistantStateTopic() string {
	return fmt.Sprintf("%v/state", ms.GetHomeAssistantBaseTopic())
}

func (ms MoistureSensor) GetHomeAssistantAvailabilityTopic() string {
	return fmt.Sprintf("%v/availability", ms.GetHomeAssistantBaseTopic())
}

func (ms MoistureSensor) GetHomeAssistantBaseTopic() string {
	return fmt.Sprintf(
		"%v/sensor/vegetable-soaker/%v",
		os.Getenv("HOME_ASSISTANT_DISCOVERY_PREFIX"),
		ms.GetHomeAssistantObjectId(),
	)
}

func (ms MoistureSensor) GetHomeAssistantObjectId() string {
	return fmt.Sprintf(
		"moisture-sensor-%v",
		ms.GetId(),
	)
}

func (ms MoistureSensor) GetId() int {
	if ms == MOISTURE_SENSOR_1 {
		return 1
	} else if ms == MOISTURE_SENSOR_2 {
		return 2
	} else if ms == MOISTURE_SENSOR_3 {
		return 3
	} else if ms == MOISTURE_SENSOR_4 {
		return 4
	} else if ms == MOISTURE_SENSOR_5 {
		return 5
	} else if ms == MOISTURE_SENSOR_6 {
		return 6
	}

	return 0
}

func MoistureSensorFromId(id int) (MoistureSensor, error) {
	if id == 1 {
		return MOISTURE_SENSOR_1, nil
	} else if id == 2 {
		return MOISTURE_SENSOR_2, nil
	} else if id == 3 {
		return MOISTURE_SENSOR_3, nil
	} else if id == 4 {
		return MOISTURE_SENSOR_4, nil
	} else if id == 5 {
		return MOISTURE_SENSOR_5, nil
	} else if id == 6 {
		return MOISTURE_SENSOR_6, nil
	}

	return MOISTURE_SENSOR_1, errors.New("invalid moisture sensor ID")
}
