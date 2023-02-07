package model

import (
	"errors"
	"time"
)

type ArduinoHeartbeat struct {
	Time time.Time
}

func MakeArduinoHeartbeatFromString(line string) (ArduinoHeartbeat, error) {
	if line != "HEARTBEAT" {
		return ArduinoHeartbeat{}, errors.New("line was not an Arduino heartbeat")
	}

	return ArduinoHeartbeat{
		Time: time.Now(),
	}, nil
}
