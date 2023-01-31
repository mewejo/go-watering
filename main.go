package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/mewejo/go-watering/arduino"
	"github.com/mewejo/go-watering/config"
)

func main() {

	app := config.GetApplication()

	ard := arduino.GetArduino()

	fmt.Println("Waiting until Arduino is ready")

	ard.WaitUntilReady()

	fmt.Println("The Arduino is ready!")

	setupCloseHandler(ard)

	go maintainMoistureLevels(ard, &app)
	go readMoistureLevels(ard, &app)
	go enforceZoneWaterOutletStates(ard, &app)

	for {
	}
}

func setupCloseHandler(ard arduino.Arduino) {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		fmt.Println("Exiting... turning water off")
		ard.SendCommand(arduino.WATER_OFF)
		ard.Port.Close()
		os.Exit(0)
	}()
}

func enforceZoneWaterOutletStates(ard arduino.Arduino, app *config.Application) {
	ticker := time.NewTicker(60 * time.Second)

	quit := make(chan struct{})

	go func() {
		for {
			select {
			case <-ticker.C:
				for _, zone := range app.Zones {
					zone.EnforceWateringState(ard)
				}

			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()
}

func maintainMoistureLevels(ard arduino.Arduino, app *config.Application) {
	ticker := time.NewTicker(10 * time.Second)

	quit := make(chan struct{})

	go func() {
		for {
			select {
			case <-ticker.C:
				for _, zone := range app.Zones {
					if len(zone.MoistureSensors) < 1 {
						continue
					}

					shouldNotBeWatering, err := zone.ShouldStopWatering()

					if err != nil {
						continue
					}

					if shouldNotBeWatering {
						zone.SetWaterState(ard, false)
						continue
					}

					shouldStartWatering, err := zone.ShouldStartWatering()

					if err != nil {
						continue
					}

					if shouldStartWatering {
						zone.SetWaterState(ard, true)
						continue
					}
				}

			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()
}

func readMoistureLevels(ard arduino.Arduino, app *config.Application) {
	ticker := time.NewTicker(2 * time.Second)

	quit := make(chan struct{})

	go func() {
		for {
			select {
			case <-ticker.C:
				readings, err := ard.GetReadings()

				if err != nil {
					log.Fatal("Could not get readings from Arduino: " + err.Error())
				}

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
