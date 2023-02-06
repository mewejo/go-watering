package model

type MoistureSensorHassState struct {
	Sensor        *MoistureSensor `json:"sensor"`
	MoistureLevel MoistureLevel   `json:"moisture_level"`
}
