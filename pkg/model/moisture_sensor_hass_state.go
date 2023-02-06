package model

type MoistureSensorHassState struct {
	Sensor  *MoistureSensor  `json:"sensor"`
	Reading *MoistureReading `json:"reading"`
}
