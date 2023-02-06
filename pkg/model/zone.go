package model

type Zone struct {
	Id                    string
	Name                  string
	Mode                  *ZoneMode
	TargetMoisture        MoistureLevel
	MoistureSensors       []*MoistureSensor
	WaterOutlets          []*WaterOutlet
	Enabled               bool
	AverageMoistureSensor *ZoneAverageMoistureSensor
}

func NewZone(id string, name string, sensors []*MoistureSensor, waterOutlets []*WaterOutlet) *Zone {
	zone := &Zone{
		Id:              id,
		Name:            name,
		Mode:            getDefaultZoneMode(),
		MoistureSensors: sensors,
		WaterOutlets:    waterOutlets,
		Enabled:         true,
		TargetMoisture:  MakeMoistureLevel(50), // TODO set to zero on boot?
	}

	zone.AverageMoistureSensor = newZoneAverageMoistureSensor(zone)

	return zone
}

func (zone Zone) MqttTopic(device *HassDevice) string {
	return "humidifier/" + device.Namespace + "/zone-" + zone.Id
}

func (zone Zone) OverriddenMqttStateTopic(device *HassDevice) string {
	return ""
}

func (zone Zone) MqttCommandTopic(device *HassDevice) string {
	return zone.MqttTopic(device) + "/" + makeZoneHassConfiguration(zone, device).CommandTopic
}

func (zone Zone) MqttStateTopic(device *HassDevice) string {
	return zone.MqttTopic(device) + "/" + makeZoneHassConfiguration(zone, device).StateTopic
}

func (zone Zone) AutoDiscoveryPayload(device *HassDevice) HassAutoDiscoverPayload {
	return makeZoneHassConfiguration(zone, device)
}
