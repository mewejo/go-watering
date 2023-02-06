package app

import (
	"log"
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

	select {
	case <-osExit:
		doExit(0)
	case <-arduinoHeartbeatStoppedChan:
		log.Println("did not receive heartbeat from Arduino in time")
		doExit(1)
	}
}

func NewApp() *App {
	return &App{}
}
