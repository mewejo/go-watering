package main

import (
	"time"

	"github.com/mewejo/go-watering/arduino"
)

func main() {

	a := arduino.GetArduino()

	for {
		a.SendCommand(arduino.WATER_1_ON)
		time.Sleep(time.Second)
		a.SendCommand(arduino.WATER_1_OFF)
		time.Sleep(time.Second)

		a.SendCommand(arduino.WATER_2_ON)
		time.Sleep(time.Second)
		a.SendCommand(arduino.WATER_2_OFF)
		time.Sleep(time.Second)

		a.SendCommand(arduino.WATER_3_ON)
		time.Sleep(time.Second)
		a.SendCommand(arduino.WATER_3_OFF)
		time.Sleep(time.Second)

		a.SendCommand(arduino.WATER_4_ON)
		time.Sleep(time.Second)
		a.SendCommand(arduino.WATER_4_OFF)
		time.Sleep(time.Second)

		time.Sleep(time.Second * 2)
	}
}
