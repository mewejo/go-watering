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

	/*
		processMoistureReadings(&app, []arduino.MoistureReading{
			{}, {},
		})

		processMoistureReadings(&app, []arduino.MoistureReading{
			{}, {},
		})

		fmt.Println(app.Zones[0].MoisureReadings)

		os.Exit(0)
	*/

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

func maintainMoistureLevels(ard arduino.Arduino, app *config.Application) {
	ticker := time.NewTicker(10 * time.Millisecond)

	quit := make(chan struct{})

	go func() {
		for {
			select {
			case <-ticker.C:
				fmt.Println("maintaining moisture levels")

				for _, zone := range app.Zones {
					requiresWatering, err := zone.RequiresWatering()

					if err != nil {
						requiresWatering = false
					}

					zone.SetWaterState(ard, requiresWatering)
				}

			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()
}

func readMoistureLevels(ard arduino.Arduino, app *config.Application) {
	ticker := time.NewTicker(10 * time.Second)

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
