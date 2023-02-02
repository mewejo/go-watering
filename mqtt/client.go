package mqtt

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/mewejo/go-watering/arduino"
	"github.com/mewejo/go-watering/config"
	"github.com/mewejo/go-watering/homeassistant"
)

var f mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	fmt.Printf("TOPIC: %s\n", msg.Topic())
	fmt.Printf("MSG: %s\n", msg.Payload())
}

func PublishHomeAssistantModeState(c mqtt.Client, zone config.Zone) mqtt.Token {
	return c.Publish(
		zone.GetHomeAssistantModeStateTopic(),
		0,
		false,
		"normal", // TODO
	)
}

func PublishHomeAssistantTargetHumidity(c mqtt.Client, zone config.Zone) mqtt.Token {
	return c.Publish(
		zone.GetHomeAssistantTargetStateHumidityTopic(),
		0,
		false,
		fmt.Sprintf("%v", zone.TargetMoisture.Percentage),
	)
}

func PublishHomeAssistantAvailability(c mqtt.Client, zone config.Zone) mqtt.Token {
	return c.Publish(
		zone.GetHomeAssistantAvailabilityTopic(),
		0,
		false,
		"online",
	)
}

func PublishHomeAssistantState(c mqtt.Client, zone config.Zone) (mqtt.Token, error) {

	moistureLevel, err := zone.AverageMoistureLevel()

	if err != nil {
		return nil, err
	}

	state := homeassistant.ZoneState{}
	state.State = "on" // TODO
	state.MoistureLevel = moistureLevel

	json, err := json.Marshal(state)

	if err != nil {
		return nil, err
	}

	return c.Publish(
		zone.GetHomeAssistantStateTopic(),
		0,
		true,
		json,
	), nil
}

func PublishHomeAsssitantAutoDiscovery(c mqtt.Client, zone config.Zone, moistureSensors []arduino.MoistureSensor) {

	var topic string
	var err error
	var zoneConfigJson []byte
	var token mqtt.Token

	// Main device
	topic = fmt.Sprintf(
		"%v/config",
		zone.GetHomeAssistantHumidifierBaseTopic(),
	)

	zoneConfigJson, err = json.Marshal(
		zone.GetHomeAssistantHumidifierConfiguration(),
	)

	if err != nil {
		log.Fatal("Could not create Home Assistant config for zone humidifier: " + zone.Id)
	}

	token = c.Publish(topic, 1, true, zoneConfigJson)
	token.Wait()

	// No sensors, no average!
	if len(zone.MoistureSensors) < 1 {
		return
	}

	// Now the average sensor
	topic = fmt.Sprintf(
		"%v/config",
		zone.GetHomeAssistantMoistureSensorBaseTopic(),
	)

	zoneConfigJson, err = json.Marshal(
		zone.GetHomeAssistantMoistureSensorConfiguration(),
	)

	if err != nil {
		log.Fatal("Could not create Home Assistant config for zone average moisture sensor: " + zone.Id)
	}

	token = c.Publish(topic, 1, true, zoneConfigJson)
	token.Wait()

	// Now the sensors on their own
	for _, sensor := range moistureSensors {
		topic = fmt.Sprintf(
			"%v/config",
			sensor.GetHomeAssistantBaseTopic(),
		)

		sensorConfigJson, err := json.Marshal(
			sensor.GetHomeAssistantMoistureSensorConfiguration(),
		)

		if err != nil {
			log.Fatal(
				fmt.Sprintf("Could not create Home Assistant config for zone moisture sensor: %v", sensor.GetId()),
			)
		}

		token = c.Publish(topic, 1, true, sensorConfigJson)
		token.Wait()
	}
}

func GetClient() mqtt.Client {
	connectionString := fmt.Sprintf(
		"tcp://%v:%v",
		os.Getenv("MQTT_HOST"),
		os.Getenv("MQTT_PORT"),
	)

	opts := mqtt.NewClientOptions()
	opts.AddBroker(connectionString)
	opts.SetClientID(os.Getenv("MQTT_CLIENT_ID"))
	opts.SetUsername(os.Getenv("MQTT_USERNAME"))
	opts.SetPassword(os.Getenv("MQTT_PASSWORD"))
	opts.SetKeepAlive(2 * time.Second)
	opts.SetDefaultPublishHandler(f)
	opts.SetPingTimeout(1 * time.Second)

	c := mqtt.NewClient(opts)

	if token := c.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	return c
}
