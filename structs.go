package main

import (
	"time"

	"github.com/mewejo/go-watering/arduino"
)

type MoistureSensor struct {
	Id   int8
	Name string
}

type MoistureReading struct {
	Time       time.Time
	Raw        int16
	Percentage int8
}

type Zone struct {
	Name            string
	MoistureSensors []MoistureSensor
	WaterOutlets    []arduino.WaterOutlet
}
