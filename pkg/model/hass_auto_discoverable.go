package model

type HassAutoDiscoverable interface {
	MqttTopic(device *HassDevice) string
	OverriddenMqttStateTopic(device *HassDevice) string
	AutoDiscoveryPayload(device *HassDevice) HassAutoDiscoverPayload
}
