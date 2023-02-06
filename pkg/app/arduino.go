package app

import (
	"errors"
	"log"
	"time"

	"github.com/mewejo/go-watering/pkg/arduino"
	"github.com/mewejo/go-watering/pkg/model"
	"github.com/mewejo/go-watering/pkg/persistence"
)

func (app *App) monitorArduinoHeartbeat() (<-chan bool, chan bool) {
	ticker := time.NewTicker(1 * time.Second)
	deadArduino := make(chan bool)
	closeTimer := make(chan bool)

	go func() {
		for {
			select {
			case <-ticker.C:
				if app.arduino.LastHeartbeat.Time.IsZero() {
					log.Println("no hb available..")
					continue
				}

				cutOff := time.Now().Add(-time.Millisecond)

				if app.arduino.LastHeartbeat.Time.Before(cutOff) {
					log.Println("dead")
					deadArduino <- true
					return
				}

			case <-closeTimer:
				ticker.Stop()
				return
			}
		}
	}()

	return deadArduino, closeTimer
}

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

func (app *App) startSendingWaterStatesToArduino() chan bool {
	ticker := time.NewTicker(1 * time.Second)

	quit := make(chan bool)

	go func() {
		for {
			select {
			case <-ticker.C:
				for _, outlet := range app.waterOutlets {
					go app.arduino.SetWaterOutletState(outlet)
				}
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()

	return quit
}

func (app *App) startRequestingWaterOutletStates() chan bool {
	ticker := time.NewTicker(1 * time.Second)

	quit := make(chan bool)

	go func() {
		for {
			select {
			case <-ticker.C:
				app.arduino.SendCommand(arduino.REQUEST_OUTLETS)
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()

	return quit
}

func (app *App) startRequestingMoistureSensorReadings() chan bool {
	ticker := time.NewTicker(1 * time.Second)

	quit := make(chan bool)

	go func() {
		for {
			select {
			case <-ticker.C:
				app.arduino.SendCommand(arduino.REQUEST_READINGS)
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()

	return quit
}

func (app *App) findMoistureSensorById(id uint) (*model.MoistureSensor, error) {
	for _, sensor := range app.moistureSensors {
		if sensor.Id != id {
			continue
		}

		return sensor, nil
	}

	return &model.MoistureSensor{}, errors.New("could not find sensor by ID")
}

func (app *App) handleArduinoDataInput(dataChan <-chan string) {

	handleHeartbeat := func(hb model.ArduinoHeartbeat) {
		app.arduino.LastHeartbeat = hb
	}

	handleMoistureReading := func(reading model.MoistureReading, sensorId uint) {
		if app.debug {
			log.Printf(
				"Sensor ID %d - Raw value: %d",
				sensorId,
				reading.Raw,
			)
		}

		persistence.RecordMoistureReading(sensorId, reading)
	}

	handleWaterOutletState := func(outletId uint, actualState bool, targetState bool) {
		for _, outlet := range app.waterOutlets {
			if outlet.Id != outletId {
				continue
			}

			outlet.ActualState = actualState
			app.sendWaterOutletStateToHass(outlet)
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

			sensor, err := app.findMoistureSensorById(sensorId)

			if err == nil {
				moistureReading.CalculateMoistureLevelForSensor(sensor)
				go handleMoistureReading(moistureReading, sensorId)
			}

			continue
		}

		outletId, actualState, targetState, err := model.DecodeWaterOutletStateFromString(line)

		if err == nil {
			go handleWaterOutletState(outletId, actualState, targetState)
			continue
		}
	}
}
