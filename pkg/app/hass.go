package app

import (
	"os"
	"time"

	"github.com/mewejo/go-watering/pkg/hass"
	"github.com/mewejo/go-watering/pkg/model"
)

func (app *App) markHassNotAvailable() {
	app.hass.Publish(
		hass.MakeMqttMessage(
			app.hassDevice.GetFqAvailabilityTopic(),
			app.hassDevice.PayloadNotAvailable,
		),
	).Wait()
}

func (app *App) startAvailabilityTimer() chan bool {
	ticker := time.NewTicker(5 * time.Second)

	quit := make(chan bool)

	go func() {
		for {
			select {
			case <-ticker.C:
				app.hass.Publish(
					hass.MakeMqttMessage(
						app.hassDevice.GetFqAvailabilityTopic(),
						app.hassDevice.PayloadAvailable,
					),
				)

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

	return nil
}
