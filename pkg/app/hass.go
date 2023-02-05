package app

import (
	"os"

	"github.com/mewejo/go-watering/pkg/hass"
	"github.com/mewejo/go-watering/pkg/model"
)

func (app *App) setupHass() error {

	hassDevice := model.NewHassDevice()

	app.hass = hass.NewClient(
		os.Getenv("HOME_ASSISTANT_DISCOVERY_PREFIX"),
		hassDevice,
	)

	return app.hass.Connect(
		hass.MakeMqttMessage(hassDevice.Namespace+"/availability", "unavailable"), // TODO veg soaker + unavail should be constants/generated/configurable
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
