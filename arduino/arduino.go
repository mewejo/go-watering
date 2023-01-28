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
	_, err := a.Port.Write(make([]byte, command))
	return err
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
