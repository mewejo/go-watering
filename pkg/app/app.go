package app

import (
	"os"

	"github.com/mewejo/go-watering/pkg/hass"
	"github.com/mewejo/go-watering/pkg/model"
)

type App struct {
	zones           []*model.Zone
	waterOutlets    []*model.WaterOutlet
	moistureSensors []*model.MoistureSensor
	hass            *hass.HassClient
}

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

func (app *App) setupHass() error {

	hassDevice := model.NewHassDevice()

	app.hass = hass.NewClient(
		os.Getenv("HOME_ASSISTANT_DISCOVERY_PREFIX")+"/",
		hassDevice,
	)

	return app.hass.Connect(
		hass.MakeMqttMessage("vegetable-soaker/availability", "unavailable"), // TODO veg soaker + unavail should be constants/generated/configurable
	)
}

func (app *App) publishHassAutoDiscovery() error {
	for _, moistureSensor := range app.moistureSensors {
		token, err := app.hass.PublishAutoDiscovery(moistureSensor)

		if err != nil {
			return err
		}

		token.Wait()
	}

	return nil

	/*

		for _, waterOutlet := range app.waterOutlets {
			// TODO
		}

		for _, zone := range app.zones {
			// TODO
		}

	*/
}

func (app *App) Run() {
	app.configureHardware()

	err := app.setupHass()

	if err != nil {
		panic(err)
	}

	err = app.publishHassAutoDiscovery()

	if err != nil {
		panic(err)
	}

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
