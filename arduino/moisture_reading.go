package arduino

import (
	"time"

	"github.com/mewejo/go-watering/world"
)

type MoistureReading struct {
	Time     time.Time
	Raw      int16
	Original world.MoistureLevel
	Sensor   MoistureSensor
}
