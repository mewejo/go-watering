package app

import (
	"time"

	"github.com/mewejo/go-watering/pkg/model"
)

func (app *App) regulateZones() chan bool {
	ticker := time.NewTicker(1 * time.Second)

	quit := make(chan bool)

	go func() {
		for {
			select {
			case <-ticker.C:
				for _, zone := range app.zones {
					go app.regulateZone(zone)
				}
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()

	return quit
}

func (app *App) regulateZone(zone *model.Zone) error {
	if !zone.Enabled {
		zone.SetWaterOutletState(false)
		return nil
	}

	zone.SetWaterOutletState(true)

	return nil
}
