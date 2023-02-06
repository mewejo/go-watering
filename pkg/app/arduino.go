package app

import (
	"fmt"

	"github.com/mewejo/go-watering/pkg/arduino"
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
			dataChan <- app.arduino.ReadLine()
		}
	}()

	return closeChan, dataChan
}

func (app *App) handleArduinoDataInput(dataChan <-chan string) {
	for line := range dataChan {
		fmt.Println(line)
	}
}
