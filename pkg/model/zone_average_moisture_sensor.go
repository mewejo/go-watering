package model

type ZoneAverageMoistureSensor struct {
	Id   string
	Name string
	Zone *Zone
}

func newZoneAverageMoistureSensor(zone *Zone) *ZoneAverageMoistureSensor {
	return &ZoneAverageMoistureSensor{
		Id:   "zone-" + zone.Id + "-average-moisture",
		Name: zone.Name + " Average Moisture",
		Zone: zone,
	}
}

func (sensor ZoneAverageMoistureSensor) OverriddenMqttStateTopic(device *HassDevice) string {
	return sensor.Zone.MqttStateTopic(device)
}

func (sensor ZoneAverageMoistureSensor) MqttTopic(device *HassDevice) string {
	return "sensor/" + device.Namespace + "/zone-" + sensor.Zone.Id + "-average-moisture-sensor"
}

func (sensor ZoneAverageMoistureSensor) AutoDiscoveryPayload(device *HassDevice) HassAutoDiscoverPayload {
	return makeMoistureSensorForZoneAverageMoistureSensor(sensor, device)
}
