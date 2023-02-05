package model

import "strconv"

type WaterOutlet struct {
	Id          uint
	Name        string
	TargetState bool
	ActualState bool
}

func NewWaterOutlet(id uint, name string) *WaterOutlet {
	return &WaterOutlet{
		Id:   id,
		Name: name,
	}
}

func (wo WaterOutlet) IdAsString() string {
	return strconv.FormatUint(uint64(wo.Id), 10)
}

func (wo WaterOutlet) MqttTopic(device *HassDevice) string {
	return "switch/" + device.Namespace + "/outlet-" + wo.IdAsString()
}

func (wo WaterOutlet) AutoDiscoveryPayload(device *HassDevice) HassAutoDiscoverPayload {
	return makeWaterOutletHassConfiguration(wo, device)
}
