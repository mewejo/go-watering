package app

import (
	"os"
	"strconv"

	"github.com/mewejo/go-watering/pkg/model"
)

func (app *App) configureHardware() {

	waterOutlet1 := model.NewWaterOutlet(1, "Soaker hose #1", false)
	waterOutlet2 := model.NewWaterOutlet(2, "Soaker hose #2", false)
	waterOutlet3 := model.NewWaterOutlet(3, "Soaker hose #3", false)
	waterOutlet4 := model.NewWaterOutlet(4, "Soaker hose #4", true)

	app.waterOutlets = append(app.waterOutlets, waterOutlet1)
	app.waterOutlets = append(app.waterOutlets, waterOutlet2)
	app.waterOutlets = append(app.waterOutlets, waterOutlet3)
	app.waterOutlets = append(app.waterOutlets, waterOutlet4)

	moistureSensor1 := model.NewMoistureSensor(1, "Sensor #1")
	moistureSensor2 := model.NewMoistureSensor(2, "Sensor #2")
	moistureSensor3 := model.NewMoistureSensor(3, "Sensor #3")
	moistureSensor4 := model.NewMoistureSensor(4, "Sensor #4")
	moistureSensor5 := model.NewMoistureSensor(5, "Sensor #5")
	moistureSensor6 := model.NewMoistureSensor(6, "Sensor #6")

	app.moistureSensors = append(app.moistureSensors, moistureSensor1)
	app.moistureSensors = append(app.moistureSensors, moistureSensor2)
	app.moistureSensors = append(app.moistureSensors, moistureSensor3)
	app.moistureSensors = append(app.moistureSensors, moistureSensor4)
	app.moistureSensors = append(app.moistureSensors, moistureSensor5)
	app.moistureSensors = append(app.moistureSensors, moistureSensor6)

	for _, moistureSensor := range app.moistureSensors {
		dryThreshold := os.Getenv("MOISTURE_SENSOR_" + moistureSensor.IdAsString() + "_DRY")
		wetThreshold := os.Getenv("MOISTURE_SENSOR_" + moistureSensor.IdAsString() + "_WET")

		if dryThreshold != "" {
			threshold, err := strconv.Atoi(dryThreshold)

			if err != nil {
				panic("could not get dry threshold for sensor ID " + moistureSensor.IdAsString() + ": " + err.Error())
			}

			moistureSensor.DryThreshold = uint(threshold)
		}

		if wetThreshold != "" {
			threshold, err := strconv.Atoi(wetThreshold)

			if err != nil {
				panic("could not get wet threshold for sensor ID " + moistureSensor.IdAsString() + ": " + err.Error())
			}

			moistureSensor.WetThreshold = uint(threshold)
		}
	}

	app.zones = append(app.zones, model.NewZone(
		"raised-bed-1",
		"Raised Bed #1",
		[]*model.MoistureSensor{moistureSensor1, moistureSensor2},
		[]*model.WaterOutlet{waterOutlet1},
	))

	app.zones = append(app.zones, model.NewZone(
		"raised-bed-2",
		"Raised Bed #2",
		[]*model.MoistureSensor{moistureSensor3, moistureSensor4},
		[]*model.WaterOutlet{waterOutlet2},
	))

	app.zones = append(app.zones, model.NewZone(
		"raised-bed-3",
		"Raised Bed #3",
		[]*model.MoistureSensor{moistureSensor5, moistureSensor6},
		[]*model.WaterOutlet{waterOutlet3},
	))
}
