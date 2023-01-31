package arduino

import "errors"

type MoistureSensor int

const (
	MOISTURE_SENSOR_1 MoistureSensor = iota
	MOISTURE_SENSOR_2
	MOISTURE_SENSOR_3
	MOISTURE_SENSOR_4
	MOISTURE_SENSOR_5
	MOISTURE_SENSOR_6
)

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
