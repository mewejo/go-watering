package hass

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/mewejo/go-watering/pkg/model"
)

type HassClient struct {
	client    mqtt.Client
	namespace string // prefix all MQTT topics with this NS
	device    *model.HassDevice
}

func NewClient(namespace string, device *model.HassDevice) *HassClient {
	return &HassClient{
		namespace: namespace,
		device:    device,
	}
}

func (c *HassClient) defaultMessageHandler(msg mqtt.Message) {
	// TODO
	fmt.Printf("TOPIC: %s\n", msg.Topic())
	fmt.Printf("MSG: %s\n", msg.Payload())
}

func (c *HassClient) PublishAutoDiscovery(entity model.HassAutoDiscoverable) (mqtt.Token, error) {

	json, err := json.Marshal(
		entity.AutoDiscoveryPayload(c.device),
	)

	if err != nil {
		return nil, err
	}

	return c.Publish(MakeMqttMessage(
		entity.EntityTopic(c.device)+"/config",
		string(json),
	)), nil
}

func (c *HassClient) Publish(message MqttMessage) mqtt.Token {
	return c.client.Publish(
		c.namespace+"/"+message.topic,
		message.qos,
		message.retain,
		message.payload,
	)
}

func (c *HassClient) Disconnect() {
	c.client.Disconnect(500)
}

func (c *HassClient) Connect(lwt MqttMessage) error {
	connectionString := fmt.Sprintf(
		"tcp://%v:%v",
		os.Getenv("MQTT_HOST"),
		os.Getenv("MQTT_PORT"),
	)

	var defaultHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
		c.defaultMessageHandler(msg)
	}

	opts := mqtt.NewClientOptions()
	opts.AddBroker(connectionString)
	opts.SetWill(c.namespace+lwt.topic, lwt.payload, lwt.qos, lwt.retain)
	opts.SetClientID(os.Getenv("MQTT_CLIENT_ID"))
	opts.SetUsername(os.Getenv("MQTT_USERNAME"))
	opts.SetPassword(os.Getenv("MQTT_PASSWORD"))
	opts.SetKeepAlive(2 * time.Second)
	opts.SetDefaultPublishHandler(defaultHandler)
	opts.SetPingTimeout(1 * time.Second)

	mqttClient := mqtt.NewClient(opts)

	if token := mqttClient.Connect(); token.Wait() && token.Error() != nil {
		return token.Error()
	}

	c.client = mqttClient

	return nil
}
