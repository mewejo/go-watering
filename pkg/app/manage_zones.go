package app

import (
	"log"
	"time"

	"github.com/mewejo/go-watering/pkg/constants"
	"github.com/mewejo/go-watering/pkg/model"
	"github.com/mewejo/go-watering/pkg/persistence"
)

func (app *App) regulateZones() chan bool {
	ticker := time.NewTicker(1 * time.Second)

	quit := make(chan bool)

	handleZone := func(zone *model.Zone) {
		if err := app.regulateZone(zone); err != nil {
			log.Println("regulating zone " + zone.Name + ": " + err.Error())
		}

		app.preventZoneFlooding(zone)
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

func (app *App) preventZoneFlooding(zone *model.Zone) {
	if !zone.WaterOutletsState {
		zone.Mode = model.GetDefaultZoneMode()
		return
	}

	cutoff := time.Now().Add(-zone.Mode.CutOffDuration)

	if zone.WaterOutletsStateChangedAt.Before(cutoff) {
		zone.Mode = model.GetDefaultZoneMode()
		zone.SetWaterOutletsState(false)
	}
}

func (app *App) regulateZone(zone *model.Zone) error {
	if !zone.Enabled {
		zone.SetWaterOutletsState(false)
		return nil
	}

	if zone.Mode == model.GetDefaultZoneMode() {
		averageMoisture, err := persistence.GetAverageReadingForSensorsSince(zone.MoistureSensors, constants.MOISTURE_READINGS_AVERAGE_DURATION)

		if err != nil {
			return err
		}

		var hysteresis int = 5

		averageMoisturePercent := int(averageMoisture.Percentage)
		targetMoisturePercent := int(zone.TargetMoisture.Percentage)

		if averageMoisturePercent < (targetMoisturePercent - hysteresis) {
			zone.SetWaterOutletsState(true)
			return nil
		} else if averageMoisturePercent > (targetMoisturePercent + hysteresis) {
			zone.SetWaterOutletsState(false)
			return nil
		}
	} else {
		zone.SetWaterOutletsState(true)
	}

	return nil
}
