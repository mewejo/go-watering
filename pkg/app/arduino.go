package app

import (
	"github.com/mewejo/go-watering/pkg/arduino"
	"github.com/mewejo/go-watering/pkg/model"
	"github.com/mewejo/go-watering/pkg/persistence"
)

func (app *App) initialiseArduino() (chan bool, <-chan string) {
	app.arduino = arduino.NewArduino()

	if app.arduino.FindAndOpenPort() != nil {
		panic("could not find or open Arduino port")
	}

	closeChan := make(chan bool)
	dataChan := make(chan string, 500)

	go func() {
		{
			<-closeChan
			app.arduino.ClosePort()
			return
		}
	}()

	go func() {
		for {
			select {
			case <-closeChan:
				return
			default:
				line, err := app.arduino.ReadLine()

				if err != nil {
					continue
				}

				dataChan <- line
			}
		}
	}()

	return closeChan, dataChan
}

func (app *App) handleArduinoDataInput(dataChan <-chan string) {

	handleHeartbeat := func(hb model.ArduinoHeartbeat) {
		app.arduino.LastHeartbeat = &hb
	}

	handleMoistureReading := func(reading model.MoistureReading, sensorId uint) {
		persistence.RecordMoistureReading(sensorId, reading)
	}

	handleWaterOutletState := func(outletId uint, actualState bool, targetState bool) {
		for _, outlet := range app.waterOutlets {
			if outlet.Id != outletId {
				continue
			}

			outlet.ActualState = actualState
		}
	}

	for line := range dataChan {
		heartbeat, err := model.MakeArduinoHeartbeatFromString(line)

		if err == nil {
			go handleHeartbeat(heartbeat)
			continue
		}

		moistureReading, sensorId, err := model.MakeMoistureReadingFromString(line)

		if err == nil {
			go handleMoistureReading(moistureReading, sensorId)
			continue
		}

		outletId, actualState, targetState, err := model.DecodeWaterOutletStateFromString(line)

		if err == nil {
			go handleWaterOutletState(outletId, actualState, targetState)
			continue
		}
	}
}
