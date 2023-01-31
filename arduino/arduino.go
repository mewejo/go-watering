package arduino

import (
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"go.bug.st/serial"
)

type Arduino struct {
	Port serial.Port
}

func (a Arduino) SendCommand(command Command) error {
	_, err := a.Port.Write([]byte(command))
	return err
}

func (a Arduino) ReadData(buffer []byte) (int, error) {
	return a.Port.Read(buffer)
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

		fmt.Print(data)

		if strings.Contains(data, "\n") {
			fmt.Println("got new line")
			break
		}

		fmt.Println("last data was: ")
		fmt.Print(data)
	}

	return data
}

func (a Arduino) ReadLines(until string) []string {
	lines := []string{}

	for {
		fmt.Println(lines)

		fmt.Println("about to ReadLine()")
		line := a.ReadLine()
		fmt.Println("done with ReadLine(), got: ")
		fmt.Print(line)

		lines = append(lines, line)

		if line == until {
			break
		}
	}

	return lines
}

func (a Arduino) GetReadings() ([]MoistureReading, error) {
	err := a.SendCommand(REQUEST_READINGS)

	if err != nil {
		return nil, errors.New("could not request readings")
	}

	time.Sleep(time.Millisecond * 250)

	readings := []MoistureReading{}

	fmt.Println("About to read lines")

	for _, line := range a.ReadLines("READINGS_END") {
		reading, err := MakeMoistureReadingFromString(line)

		if err != nil {
			fmt.Println(err)
			continue
		}

		readings = append(
			readings,
			reading,
		)
	}

	if len(readings) < 1 {
		return nil, errors.New("no readings returned from Arduino")
	}

	return readings, nil
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
