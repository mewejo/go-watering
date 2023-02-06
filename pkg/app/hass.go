package app

import (
	"encoding/json"
	"os"
	"strconv"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/mewejo/go-watering/pkg/constants"
	"github.com/mewejo/go-watering/pkg/hass"
	"github.com/mewejo/go-watering/pkg/model"
	"github.com/mewejo/go-watering/pkg/persistence"
)

func (app *App) listenForWaterOutletCommands() {

	subscribe := func(outlet *model.WaterOutlet) {
		app.hass.Subscribe(
			outlet.MqttCommandTopic(app.hassDevice),
			func(message mqtt.Message) {
				if string(message.Payload()) == constants.HASS_STATE_ON {
					outlet.TargetState = true
				} else if string(message.Payload()) == constants.HASS_STATE_OFF {
					outlet.TargetState = false
				}

				app.arduino.SetWaterOutletState(outlet)
			},
		)
	}

	for _, outlet := range app.waterOutlets {
		if !outlet.IndependentlyControlled {
			continue
		}

		subscribe(outlet)
	}
}

func (app *App) listenForZoneCommands() {

	subscribe := func(zone *model.Zone) {
		app.hass.Subscribe(
			zone.MqttCommandTopic(app.hassDevice),
			func(message mqtt.Message) {
				if string(message.Payload()) == constants.HASS_STATE_ON {
					zone.Enabled = true
				} else if string(message.Payload()) == constants.HASS_STATE_OFF {
					zone.Enabled = false
				}
			},
		)

		app.hass.Subscribe(
			zone.MqttTargetMoistureCommandTopic(app.hassDevice),
			func(message mqtt.Message) {
				moisturePercent, err := strconv.Atoi(string(message.Payload()))

				if err != nil {
					return
				}

				zone.TargetMoisture = model.MakeMoistureLevel(uint(moisturePercent))
			},
		)
	}

	for _, zone := range app.zones {
		subscribe(zone)
	}
}

func (app *App) publishWaterOutletState(outlet *model.WaterOutlet) error {

	payload, err := json.Marshal(outlet)

	if err != nil {
		return err
	}

	app.hass.Publish(
		hass.MakeMqttMessage(
			outlet.MqttStateTopic(app.hassDevice),
			string(payload),
		),
	)

	return nil
}

func (app *App) sendZoneStateToHas(zone *model.Zone) error {

	average, err := persistence.GetAverageReadingForSensorsSince(zone.MoistureSensors, 2*time.Minute)

	if err != nil {
		return err
	}

	state := model.MakeZoneHassState(
		zone,
		average,
	)

	payload, err := json.Marshal(state)

	if err != nil {
		return err
	}

	app.hass.Publish(
		hass.MakeMqttMessage(
			zone.MqttStateTopic(app.hassDevice),
			string(payload),
		),
	)

	return nil

}

func (app *App) startSendingZoneStateToHass() chan bool {
	ticker := time.NewTicker(1 * time.Second)

	quit := make(chan bool)

	go func() {
		for {
			select {
			case <-ticker.C:
				for _, zone := range app.zones {
					go app.sendZoneStateToHas(zone)
				}
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()

	return quit
}

func (app *App) startSendingMoistureSensorReadingsToHass() chan bool {

	ticker := time.NewTicker(5 * time.Second)

	quit := make(chan bool)

	sendSensorStates := func() {
		for _, sensor := range app.moistureSensors {
			go app.publishMoistureSensorStateToHass(sensor)
		}
	}

	go func() {
		for {
			select {
			case <-ticker.C:
				go sendSensorStates()

			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()

	return quit
}

func (app *App) publishMoistureSensorStateToHass(sensor *model.MoistureSensor) error {

	moistureLevel, err := persistence.GetAverageReadingForSensorIdSince(sensor.Id, 2*time.Minute)

	if err != nil {
		return err
	}

	state := model.MoistureSensorHassState{
		Sensor:        sensor,
		MoistureLevel: moistureLevel,
	}

	payload, err := json.Marshal(state)

	if err != nil {
		return err
	}

	app.hass.Publish(
		hass.MakeMqttMessage(
			sensor.MqttStateTopic(app.hassDevice),
			string(payload),
		),
	)

	return nil
}

func (app *App) markHassNotAvailable() {
	app.hass.Publish(
		hass.MakeMqttMessage(
			app.hassDevice.GetFqAvailabilityTopic(),
			app.hassDevice.PayloadNotAvailable,
		),
	).Wait()
}

func (app *App) startHassAvailabilityTimer() chan bool {
	ticker := time.NewTicker(5 * time.Second)

	quit := make(chan bool)

	sendAvailableMessage := func() {
		app.hass.Publish(
			hass.MakeMqttMessage(
				app.hassDevice.GetFqAvailabilityTopic(),
				app.hassDevice.PayloadAvailable,
			),
		)
	}

	sendAvailableMessage()

	go func() {
		for {
			select {
			case <-ticker.C:
				sendAvailableMessage()

			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()

	return quit
}

func (app *App) setupHass() error {

	app.hassDevice = model.NewHassDevice()

	app.hass = hass.NewClient(
		os.Getenv("HOME_ASSISTANT_DISCOVERY_PREFIX"),
		app.hassDevice,
	)

	return app.hass.Connect(
		hass.MakeMqttMessage(app.hassDevice.GetFqAvailabilityTopic(), app.hassDevice.PayloadNotAvailable),
	)
}

func (app *App) publishHassAutoDiscovery() error {
	for _, entity := range app.moistureSensors {
		token, err := app.hass.PublishAutoDiscovery(entity)

		if err != nil {
			return err
		}

		token.Wait()
	}

	for _, entity := range app.waterOutlets {

		if !entity.IndependentlyControlled {
			continue
		}

		token, err := app.hass.PublishAutoDiscovery(entity)

		if err != nil {
			return err
		}

		token.Wait()
	}

	for _, entity := range app.zones {
		token, err := app.hass.PublishAutoDiscovery(entity)

		if err != nil {
			return err
		}

		token.Wait()
	}

	for _, entity := range app.zones {
		token, err := app.hass.PublishAutoDiscovery(entity.AverageMoistureSensor)

		if err != nil {
			return err
		}

		token.Wait()
	}

	return nil
}
