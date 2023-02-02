package mqtt

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/mewejo/go-watering/config"
	"github.com/mewejo/go-watering/homeassistant"
	"github.com/mewejo/go-watering/world"
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
		"33", // TODO
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
	state := homeassistant.ZoneState{}
	state.State = "on" // TODO
	state.MoistureLevel = world.MoistureLevel{
		Percentage: 66, // TODO
	}

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

func PublishHomeAsssitantAutoDiscovery(c mqtt.Client, zone config.Zone) {

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
