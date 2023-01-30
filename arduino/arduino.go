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

func (a Arduino) SetWaterState(outlet int, state bool) error {
	var command Command

	waterStateCommand := func(state bool, on Command, off Command) Command {
		if state {
			return on
		} else {
			return off
		}
	}

	if outlet == 1 {
		command = waterStateCommand(state, WATER_1_ON, WATER_1_OFF)
	} else if outlet == 2 {
		command = waterStateCommand(state, WATER_2_ON, WATER_2_OFF)
	} else if outlet == 3 {
		command = waterStateCommand(state, WATER_3_ON, WATER_3_OFF)
	} else if outlet == 4 {
		command = waterStateCommand(state, WATER_4_ON, WATER_4_OFF)
	} else {
		return errors.New("Invalid outlet")
	}

	return a.SendCommand(command)
}

func findArduinoPort() (string, error) {
	ports, err := serial.GetPortsList()

	if err != nil {
		log.Fatal(err)
	}

	if len(ports) == 0 {
		log.Fatal("No serial ports found!")
	}

	for _, port := range ports {
		if !strings.Contains(port, "ttyACM") {
			continue
		}

		return port, nil
	}

	return "", errors.New("No devices found which look like an Arduino")
}

func GetArduino() Arduino {

	arduinoPort, err := findArduinoPort()

	if err != nil {
		log.Fatal("Could not find Arduino port! " + err.Error())
	}

	mode := &serial.Mode{
		BaudRate: 9600,
	}

	port, err := serial.Open(arduinoPort, mode)

	if err != nil {
		log.Fatal("Could not open Arduino port! " + err.Error())
	}

	arduino := Arduino{
		Port: port,
	}

	return arduino
}
