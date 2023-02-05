package model

type Zone struct {
	Id              string
	Name            string
	Mode            *ZoneMode
	MoistureSensors []*MoistureSensor
	WaterOutlets    []*WaterOutlet
}

func NewZone(id string, name string, sensors []*MoistureSensor, waterOutlets []*WaterOutlet) *Zone {
	return &Zone{
		Id:              id,
		Name:            name,
		Mode:            getDefaultZoneMode(),
		MoistureSensors: sensors,
		WaterOutlets:    waterOutlets,
	}
}

func (zone Zone) MqttTopic(device *HassDevice) string {
	return "humidifier/" + device.Namespace + "/zone-" + zone.Id
}

func (zone Zone) AutoDiscoveryPayload(device *HassDevice) HassAutoDiscoverPayload {
	return makeZoneHassConfiguration(zone, device)
}
