package model

import (
	"strconv"

	"github.com/mewejo/go-watering/pkg/number"
)

type MoistureSensor struct {
	Id           uint
	Name         string
	DryThreshold uint `json:"dry_threshold"`
	WetThreshold uint `json:"wet_threshold"`
}

func (ms MoistureSensor) MqttTopic(device *HassDevice) string {
	return "sensor/" + device.Namespace + "/sensor-" + ms.IdAsString()
}

func (ms MoistureSensor) MqttStateTopic(device *HassDevice) string {
	return ms.MqttTopic(device) + "/" + makeMoistureSensorHassConfiguration(ms, device).StateTopic
}

func (ms MoistureSensor) AutoDiscoveryPayload(device *HassDevice) HassAutoDiscoverPayload {
	return makeMoistureSensorHassConfiguration(ms, device)
}

func (ms MoistureSensor) IdAsString() string {
	return strconv.FormatUint(uint64(ms.Id), 10)
}

func (ms MoistureSensor) mapRawReadingToPercentage(raw uint) uint {
	return uint(number.ChangeRange(
		float64(raw),
		float64(ms.DryThreshold),
		float64(ms.WetThreshold),
		0,
		100,
	))
}

func MakeMoistureSensor(id uint, name string) MoistureSensor {
	return MoistureSensor{
		Id:           id,
		Name:         name,
		DryThreshold: 500,
		WetThreshold: 240,
	}
}
