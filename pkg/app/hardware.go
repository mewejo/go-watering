package app

import "github.com/mewejo/go-watering/pkg/model"

func (app *App) configureHardware() {

	waterOutlet1 := model.NewWaterOutlet(1, "Soaker hose #1")
	waterOutlet2 := model.NewWaterOutlet(2, "Soaker hose #2")
	waterOutlet3 := model.NewWaterOutlet(3, "Soaker hose #3")
	waterOutlet4 := model.NewWaterOutlet(4, "Soaker hose #4")

	// The only outlet which isn't tied to a zone.
	app.waterOutlets = append(app.waterOutlets, waterOutlet4)

	moistureSensor1 := model.MakeMoistureSensor(1, "Sensor #1")
	moistureSensor2 := model.MakeMoistureSensor(2, "Sensor #2")
	moistureSensor3 := model.MakeMoistureSensor(3, "Sensor #3")
	moistureSensor4 := model.MakeMoistureSensor(4, "Sensor #4")
	moistureSensor5 := model.MakeMoistureSensor(5, "Sensor #5")
	moistureSensor6 := model.MakeMoistureSensor(6, "Sensor #6")

	app.moistureSensors = append(app.moistureSensors, &moistureSensor1)
	app.moistureSensors = append(app.moistureSensors, &moistureSensor2)
	app.moistureSensors = append(app.moistureSensors, &moistureSensor3)
	app.moistureSensors = append(app.moistureSensors, &moistureSensor4)
	app.moistureSensors = append(app.moistureSensors, &moistureSensor5)
	app.moistureSensors = append(app.moistureSensors, &moistureSensor6)

	app.zones = append(app.zones, model.NewZone(
		"raised-bed-1",
		"Raised Bed #1",
		[]*model.MoistureSensor{&moistureSensor1, &moistureSensor2},
		[]*model.WaterOutlet{waterOutlet1},
	))

	app.zones = append(app.zones, model.NewZone(
		"raised-bed-2",
		"Raised Bed #2",
		[]*model.MoistureSensor{&moistureSensor3, &moistureSensor4},
		[]*model.WaterOutlet{waterOutlet2},
	))

	app.zones = append(app.zones, model.NewZone(
		"raised-bed-3",
		"Raised Bed #3",
		[]*model.MoistureSensor{&moistureSensor5, &moistureSensor6},
		[]*model.WaterOutlet{waterOutlet3},
	))
}