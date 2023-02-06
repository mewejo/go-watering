package app

import (
	"time"

	"github.com/mewejo/go-watering/pkg/model"
	"github.com/mewejo/go-watering/pkg/persistence"
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
		outlet.TargetState = zone.WaterOutletsState
	}
}

func (app *App) regulateZone(zone *model.Zone) error {
	if !zone.Enabled {
		zone.SetWaterOutletsState(false)
		return nil
	}

	if zone.Mode.Key == "normal" {
		averageMoisture, err := persistence.GetAverageReadingForSensorsSince(zone.MoistureSensors, 2*time.Minute)

		if err != nil {
			return err
		}

		var hysteresis uint = 5

		if averageMoisture.Percentage < (zone.TargetMoisture.Percentage - hysteresis) {
			zone.SetWaterOutletsState(true)
			return nil
		} else if averageMoisture.Percentage > (zone.TargetMoisture.Percentage + hysteresis) {
			zone.SetWaterOutletsState(false)
			return nil
		}
	}

	return nil
}
