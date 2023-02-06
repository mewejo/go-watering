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

func (m MqttMessage) Retained() MqttMessage {
	m.retain = true
	return m
}

func (m MqttMessage) Qos(qos byte) MqttMessage {
	m.qos = qos
	return m
}
