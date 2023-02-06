package app

import (
	"time"

	"github.com/mewejo/go-watering/pkg/model"
)

func (app *App) regulateZones() chan bool {
	ticker := time.NewTicker(1 * time.Second)

	quit := make(chan bool)

	handleZone := func(zone *model.Zone) {
		app.regulateZone(zone)
		app.ensureZoneWaterOutletState(zone)
	}

	go func() {
		for {
			select {
			case <-ticker.C:
				for _, zone := range app.zones {
					go handleZone(zone)
				}
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()

	return quit
}

func (app *App) ensureZoneWaterOutletState(zone *model.Zone) {
	for _, outlet := range zone.WaterOutlets {
		outlet.TargetState = zone.WaterOutletsOpen
	}
}

func (app *App) regulateZone(zone *model.Zone) error {
	if !zone.Enabled {
		zone.WaterOutletsOpen = false
		return nil
	}

	return nil
}
