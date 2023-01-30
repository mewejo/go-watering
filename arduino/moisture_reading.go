package arduino

import "time"

type MoistureReading struct {
	Time       time.Time
	Raw        int16
	Percentage int8
	Sensor     MoistureSensor
}
