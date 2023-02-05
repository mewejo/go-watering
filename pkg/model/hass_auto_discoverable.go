package model

type HassAutoDiscoverable interface {
	AutoDiscoveryTopic(device *HassDevice) string
	AutoDiscoveryPayload(device *HassDevice) interface{}
}
