package model

import "time"

type MoistureReading struct {
	Time          time.Time
	MoistureLevel MoistureLevel
}
