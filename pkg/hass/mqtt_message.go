package hass

type MqttMessage struct {
	topic   string
	payload string
	retain  bool
	qos     byte
}

func MakeMqttMessage(topic string, payload string) MqttMessage {
	return MqttMessage{
		topic:   topic,
		payload: payload,
	}
}
