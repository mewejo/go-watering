package app

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/mewejo/go-watering/pkg/arduino"
	"github.com/mewejo/go-watering/pkg/hass"
	"github.com/mewejo/go-watering/pkg/model"
)

type App struct {
	zones           []*model.Zone
	waterOutlets    []*model.WaterOutlet
	moistureSensors []*model.MoistureSensor
	hass            *hass.HassClient
	hassDevice      *model.HassDevice
	arduino         *arduino.Arduino
}

func (app *App) setupCloseHandler() <-chan os.Signal {
	sigChan := make(chan os.Signal)

	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	return sigChan
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

	app.startHassAvailabilityTimer()

	osExit := app.setupCloseHandler()

	closeArduinoChan, arduinoInputChan := app.initialiseArduino()
	stopRequestingOutletStatesChan := app.startRequestingWaterOutletStates()
	stopRequestingMoistureSensorReadingsChan := app.startRequestingMoistureSensorReadings()
	stopSendingOutletStatesToArduinoChan := app.startSendingWaterStatesToArduino()
	stopSendingMoistureSensorReadingsToHassChan := app.startSendingMoistureSensorReadingsToHass()
	stopSendingZoneStatesToHassChan := app.startSendingZoneStateToHass()

	app.listenForWaterOutletCommands()

	go app.handleArduinoDataInput(arduinoInputChan)

	{
		<-osExit
		close(stopRequestingOutletStatesChan)
		close(stopSendingOutletStatesToArduinoChan)
		close(stopRequestingMoistureSensorReadingsChan)
		close(stopSendingMoistureSensorReadingsToHassChan)
		close(closeArduinoChan)
		close(stopSendingZoneStatesToHassChan)
		app.markHassNotAvailable()
		app.hass.Disconnect()
		os.Exit(0)
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
