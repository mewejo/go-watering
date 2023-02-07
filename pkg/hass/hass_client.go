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
	Namespace string // prefix all MQTT topics with this NS
	Device    *model.HassDevice
}

func NewClient(namespace string, device *model.HassDevice) *HassClient {
	return &HassClient{
		Namespace: namespace,
		Device:    device,
	}
}

func (c *HassClient) defaultMessageHandler(msg mqtt.Message) {
	// TODO
	fmt.Printf("TOPIC: %s\n", msg.Topic())
	fmt.Printf("MSG: %s\n", msg.Payload())
}

func (c *HassClient) PublishAutoDiscovery(entity model.HassAutoDiscoverable) (mqtt.Token, error) {

	payload := entity.AutoDiscoveryPayload(c.Device).WithGlobalTopicPrefix(c.Namespace, c.Device, entity)

	json, err := json.Marshal(
		payload,
	)

	if err != nil {
		return nil, err
	}

	return c.Publish(MakeMqttMessage(
		entity.MqttTopic(c.Device)+"/config",
		string(json),
	)), nil
}

func (c *HassClient) Publish(message MqttMessage) mqtt.Token {
	return c.client.Publish(
		c.Namespace+"/"+message.topic,
		message.qos,
		message.retain,
		message.payload,
	)
}

func (c *HassClient) Subscribe(topic string, handler MessageHandler) {

	var passHandler = func(client mqtt.Client, message mqtt.Message) {
		handler(message)
	}

	c.client.Subscribe(
		c.Namespace+"/"+topic,
		2,
		passHandler,
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

	opts := mqtt.NewClientOptions()
	opts.AddBroker(connectionString)
	opts.SetWill(c.Namespace+"/"+lwt.topic, lwt.payload, lwt.qos, lwt.retain)
	opts.SetClientID(os.Getenv("MQTT_CLIENT_ID"))
	opts.SetUsername(os.Getenv("MQTT_USERNAME"))
	opts.SetPassword(os.Getenv("MQTT_PASSWORD"))
	opts.SetKeepAlive(2 * time.Second)
	opts.SetPingTimeout(1 * time.Second)

	mqttClient := mqtt.NewClient(opts)

	if token := mqttClient.Connect(); token.Wait() && token.Error() != nil {
		return token.Error()
	}

	c.client = mqttClient

	return nil
}
