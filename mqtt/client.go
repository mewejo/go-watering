package mqtt

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/mewejo/go-watering/config"
)

var f mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	fmt.Printf("TOPIC: %s\n", msg.Topic())
	fmt.Printf("MSG: %s\n", msg.Payload())
}

func PublishHomeAssistantAvailability(c mqtt.Client, zone config.Zone) {
	c.Publish(
		zone.GetHomeAssistantAvailabilityTopic(),
		0,
		false,
		"online",
	)
}

func PublishHomeAssistantState(c mqtt.Client, zone config.Zone) {
	c.Publish(
		zone.GetHomeAssistantStateTopic(),
		0,
		false,
		"{\"humidity\": 33}", // TODO
	)
}

func PublishHomeAsssitantAutoDiscovery(c mqtt.Client, zone config.Zone) {

	topic := fmt.Sprintf(
		"%v/humidifier/%v/config",
		os.Getenv("HOME_ASSISTANT_DISCOVERY_PREFIX"),
		zone.GetHomeAssistantObjectId(),
	)

	zoneConfigJson, err := json.Marshal(
		zone.GetHomeAssistantConfiguration(),
	)

	if err != nil {
		log.Fatal("Could not create Home Assistant config for zone: " + zone.Id)
	}

	token := c.Publish(topic, 1, true, zoneConfigJson)
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
