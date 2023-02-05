package app

import (
	"github.com/mewejo/go-watering/pkg/model"
)

type App struct {
	Zones           []*model.Zone
	WaterOutlets    []*model.WaterOutlet
	MoistureSensors []*model.MoistureSensor
}

func (app *App) configure() {

	waterOutlet1 := model.NewWaterOutlet(1, "Soaker hose #1")
	waterOutlet2 := model.NewWaterOutlet(2, "Soaker hose #2")
	waterOutlet3 := model.NewWaterOutlet(3, "Soaker hose #3")
	waterOutlet4 := model.NewWaterOutlet(4, "Soaker hose #4")

	// The only outlet which isn't tied to a zone.
	app.WaterOutlets = append(app.WaterOutlets, waterOutlet4)

	moistureSensor1 := model.MakeMoistureSensor(1, "Sensor #1")
	moistureSensor2 := model.MakeMoistureSensor(2, "Sensor #2")
	moistureSensor3 := model.MakeMoistureSensor(3, "Sensor #3")
	moistureSensor4 := model.MakeMoistureSensor(4, "Sensor #4")
	moistureSensor5 := model.MakeMoistureSensor(5, "Sensor #5")
	moistureSensor6 := model.MakeMoistureSensor(6, "Sensor #6")

	app.MoistureSensors = append(app.MoistureSensors, &moistureSensor1)
	app.MoistureSensors = append(app.MoistureSensors, &moistureSensor2)
	app.MoistureSensors = append(app.MoistureSensors, &moistureSensor3)
	app.MoistureSensors = append(app.MoistureSensors, &moistureSensor4)
	app.MoistureSensors = append(app.MoistureSensors, &moistureSensor5)
	app.MoistureSensors = append(app.MoistureSensors, &moistureSensor6)

	app.Zones = append(app.Zones, model.NewZone(
		"raised-bed-1",
		"Raised Bed #1",
		[]*model.MoistureSensor{&moistureSensor1, &moistureSensor2},
		[]*model.WaterOutlet{waterOutlet1},
	))

	app.Zones = append(app.Zones, model.NewZone(
		"raised-bed-2",
		"Raised Bed #2",
		[]*model.MoistureSensor{&moistureSensor3, &moistureSensor4},
		[]*model.WaterOutlet{waterOutlet2},
	))

	app.Zones = append(app.Zones, model.NewZone(
		"raised-bed-3",
		"Raised Bed #3",
		[]*model.MoistureSensor{&moistureSensor5, &moistureSensor6},
		[]*model.WaterOutlet{waterOutlet3},
	))
}

func (app *App) Run() {
	app.configure()

	/*
		Make zone configurations
		Connect to MQTT
		Set LWT for availability topic (shared by all entities). Support modes normal/boost
		Publish MQTT auto discovery for Zones (climate), moisture sensors, and outlets not attached to zones (or zones with no sensors?)
		Find Arduino port
		Open serial connection
		Wait for heartbeat
		Wait for above to be ready
		Loop
			Read & process sensors and water states
			Publish zone states (freq as below)
			Publish sensor states (every 5 min in prod, 2 sec in testing)
			Publish outlets without zones states
			Check for heartbeat - if none after X:
				Close Arduino serial connection
				Attempt to establish a new connection
				Wait for a heartbeat
				Loop as above
		Program exit
			Send water off command
			Send unavailable topic to MQTT
			Exit


	*/

}

func NewApp() *App {
	return &App{}
}
