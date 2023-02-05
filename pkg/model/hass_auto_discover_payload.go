package model

type HassAutoDiscoverPayload interface {
	WithGlobalTopicPrefix(prefix string) HassAutoDiscoverPayload
}
