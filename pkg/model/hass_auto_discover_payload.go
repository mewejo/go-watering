package model

type HassAutoDiscoverPayload interface {
	WithGlobalTopicPrefix(prefix string, device *HassDevice, entity HassAutoDiscoverable) HassAutoDiscoverPayload
}
