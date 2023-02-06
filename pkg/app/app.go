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
	debug           bool
}

func (app *App) setupCloseHandler() <-chan os.Signal {
	sigChan := make(chan os.Signal)

	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	return sigChan
}

func (app *App) Run() {

	app.debug = os.Getenv("APP_DEBUG") == "true"

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
	app.sendInitialZoneModesToHass()

	osExit := app.setupCloseHandler()

	closeArduinoChan, arduinoInputChan := app.initialiseArduino()
	arduinoHeartbeatStoppedChan, stopMonitoringArduinoHeartbeatChan := app.monitorArduinoHeartbeat()
	stopRequestingOutletStatesChan := app.startRequestingWaterOutletStates()
	stopRequestingMoistureSensorReadingsChan := app.startRequestingMoistureSensorReadings()
	stopSendingOutletStatesToArduinoChan := app.startSendingWaterStatesToArduino()
	stopSendingMoistureSensorReadingsToHassChan := app.startSendingMoistureSensorReadingsToHass()
	stopSendingZoneStatesToHassChan := app.startSendingZoneStateToHass()
	stopRegulatingZonesChan := app.regulateZones()

	app.listenForWaterOutletCommands()
	app.listenForZoneCommands()

	go app.handleArduinoDataInput(arduinoInputChan)

	doExit := func(code int) {
		close(stopMonitoringArduinoHeartbeatChan)
		close(stopRequestingOutletStatesChan)
		app.forceSetAllWaterOutletStates(false)
		close(stopSendingOutletStatesToArduinoChan)
		close(stopRequestingMoistureSensorReadingsChan)
		close(stopSendingMoistureSensorReadingsToHassChan)
		close(closeArduinoChan)
		close(stopSendingZoneStatesToHassChan)
		close(stopRegulatingZonesChan)
		app.markHassNotAvailable()
		app.hass.Disconnect()
		os.Exit(code)
	}

	{
		<-osExit
		doExit(0)
		<-arduinoHeartbeatStoppedChan
		doExit(1)
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
