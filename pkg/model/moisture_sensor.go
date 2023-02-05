package model

type MoistureSensor struct {
	Id   uint
	Name string
	// TODO there will be props to translate the raw readings from Arduino into a percentage
}
