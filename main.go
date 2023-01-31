package main

import (
	"time"

	"github.com/mewejo/go-watering/arduino"
)

func main() {

	//app := config.GetApplication()

	//app.Zones[0].RecordMoistureReading(arduino.MoistureReading{})
	//app.Zones[0].RecordMoistureReading(arduino.MoistureReading{})

	a := arduino.GetArduino()

	for {
		a.SendCommand(arduino.WATER_1_ON)
		time.Sleep(time.Millisecond * 500)
		a.SendCommand(arduino.WATER_1_OFF)
		time.Sleep(time.Second)

		a.SendCommand(arduino.WATER_2_ON)
		time.Sleep(time.Millisecond * 500)
		a.SendCommand(arduino.WATER_2_OFF)
		time.Sleep(time.Second)

		a.SendCommand(arduino.WATER_3_ON)
		time.Sleep(time.Millisecond * 500)
		a.SendCommand(arduino.WATER_3_OFF)
		time.Sleep(time.Second)

		a.SendCommand(arduino.WATER_4_ON)
		time.Sleep(time.Millisecond * 500)
		a.SendCommand(arduino.WATER_4_OFF)
		time.Sleep(time.Second)

		time.Sleep(time.Second * 2)
	}
}
