package app

import (
	"github.com/mewejo/go-watering/pkg/arduino"
	"github.com/mewejo/go-watering/pkg/model"
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

	}

	handleWaterOutletState := func(outletId uint, realState bool, setState bool) {

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

		outletId, realState, setState, err := model.DecodeWaterOutletStateFromString(line)

		if err == nil {
			go handleWaterOutletState(outletId, realState, setState)
			continue
		}
	}
}
