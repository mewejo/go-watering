package model

type HassAutoDiscoverable interface {
	AutoDiscoveryTopic() string
	AutoDiscoveryPayload(device *HassDevice) interface{}
}
