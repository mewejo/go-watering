package app

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	mqtt "github.com/eclipse/paho.mqtt.golang"
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

	go app.handleArduinoDataInput(arduinoInputChan)

	app.hass.Subscribe("switch/vegetable-soaker/outlet-4/command", func(m mqtt.Message) {
		fmt.Println(string(m.Payload()))
	})

	{
		<-osExit
		app.markHassNotAvailable()
		close(closeArduinoChan)
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
