package arduino

import (
	"errors"
	"log"
	"strings"

	"go.bug.st/serial"
)

type Arduino struct {
	port serial.Port
}

func (a Arduino) SendCommand(command Command) (int, error) {
	return a.port.Write([]byte(command))
}

func (a Arduino) ReadData(buffer []byte) (int, error) {
	return a.port.Read(buffer)
}

func findArduinoPort() (string, error) {
	ports, err := serial.GetPortsList()

	if err != nil {
		return "", err
	}

	if len(ports) == 0 {
		return "", errors.New("no serial ports found!")
	}

	for _, port := range ports {
		if !strings.Contains(port, "ttyACM") {
			continue
		}

		return port, nil
	}

	return "", errors.New("no devices found which look like an Arduino")
}

func (a Arduino) ReadLine() string {
	buff := make([]byte, 1)
	data := ""

	for {
		n, err := a.ReadData(buff)

		if err != nil {
			log.Fatal(err)
		}

		if n == 0 {
			break
		}

		data += string(buff[:n])

		if strings.Contains(data, "\n") {
			break
		}
	}

	data = strings.TrimSuffix(data, "\n")
	data = strings.TrimSuffix(data, "\r")

	return data
}

func (a *Arduino) ClosePort() error {
	return a.port.Close()
}

func (a *Arduino) FindAndOpenPort() error {
	arduinoPort, err := findArduinoPort()

	if err != nil {
		return err
	}

	mode := &serial.Mode{
		BaudRate: 9600,
	}

	port, err := serial.Open(arduinoPort, mode)

	if err != nil {
		return err
	}

	a.port = port

	return nil
}

func NewArduino() *Arduino {

	return &Arduino{}
}
