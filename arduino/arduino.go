package arduino

import (
	"errors"
	"log"
	"strings"

	"go.bug.st/serial"
)

type Arduino struct {
	Port serial.Port
}

func (a Arduino) SendCommand(command Command) error {
	_, err := a.Port.Write([]byte(command))
	return err
}

func (a Arduino) SetAllWaterState(state bool) error {
	if state {
		return a.SendCommand(WATER_ON)
	} else {
		return a.SendCommand(WATER_OFF)
	}
}

func (a Arduino) SetWaterState(outlet WaterOutlet, state bool) error {

	var err error
	var command Command

	if state {
		command, err = outlet.OnCommand()
	} else {
		command, err = outlet.OffCommand()
	}

	if nil != err {
		return err
	}

	return a.SendCommand(command)
}

func findArduinoPort() (string, error) {
	ports, err := serial.GetPortsList()

	if err != nil {
		log.Fatal(err)
	}

	if len(ports) == 0 {
		log.Fatal("no serial ports found!")
	}

	for _, port := range ports {
		if !strings.Contains(port, "ttyACM") {
			continue
		}

		return port, nil
	}

	return "", errors.New("no devices found which look like an Arduino")
}

func GetArduino() Arduino {

	arduinoPort, err := findArduinoPort()

	if err != nil {
		log.Fatal("could not find Arduino port! " + err.Error())
	}

	mode := &serial.Mode{
		BaudRate: 9600,
	}

	port, err := serial.Open(arduinoPort, mode)

	if err != nil {
		log.Fatal("could not open Arduino port! " + err.Error())
	}

	arduino := Arduino{
		Port: port,
	}

	return arduino
}
