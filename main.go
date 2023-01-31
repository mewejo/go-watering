package main

import (
	"fmt"
	"log"
	"time"

	"github.com/mewejo/go-watering/arduino"
	"github.com/mewejo/go-watering/config"
)

func main() {

	app := config.GetApplication()

	//app.Zones[0].RecordMoistureReading(arduino.MoistureReading{})
	//app.Zones[0].RecordMoistureReading(arduino.MoistureReading{})

	ard := arduino.GetArduino()

	fmt.Println("Waiting until Arduino is ready")

	ard.WaitUntilReady()

	fmt.Println("The Arduino is ready!")

	go readMoistureLevels(ard, &app)

	for {
		ard.SendCommand(arduino.WATER_1_ON)
		time.Sleep(time.Millisecond * 500)
		ard.SendCommand(arduino.WATER_1_OFF)
		time.Sleep(time.Second)

		ard.SendCommand(arduino.WATER_2_ON)
		time.Sleep(time.Millisecond * 500)
		ard.SendCommand(arduino.WATER_2_OFF)
		time.Sleep(time.Second)

		ard.SendCommand(arduino.WATER_3_ON)
		time.Sleep(time.Millisecond * 500)
		ard.SendCommand(arduino.WATER_3_OFF)
		time.Sleep(time.Second)

		ard.SendCommand(arduino.WATER_4_ON)
		time.Sleep(time.Millisecond * 500)
		ard.SendCommand(arduino.WATER_4_OFF)
		time.Sleep(time.Second)

		time.Sleep(time.Second * 2)
	}
}

func readMoistureLevels(ard arduino.Arduino, app *config.Application) {
	ticker := time.NewTicker(100 * time.Millisecond)

	quit := make(chan struct{})

	go func() {
		for {
			select {
			case <-ticker.C:
				fmt.Println("Asking for moisture readings")

				readings, err := ard.GetReadings()

				if err != nil {
					log.Fatal("Could not get readings from Arduino: " + err.Error())
				}

				fmt.Println("Got moisture readings")

				processMoistureReadings(app, readings)
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()
}

func processMoistureReadings(app *config.Application, readings []arduino.MoistureReading) {
	for _, zone := range app.Zones {
		for _, sensor := range zone.MoistureSensors {
			for _, reading := range readings {
				if sensor != reading.Sensor {
					continue
				}

				// This 3 level nesting feels nasty

				zone.RecordMoistureReading(reading)
			}
		}
	}
}
