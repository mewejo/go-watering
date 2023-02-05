package model

type HassAutoDiscoverable interface {
	EntityTopic(device *HassDevice) string
	AutoDiscoveryPayload(device *HassDevice) interface{}
}
