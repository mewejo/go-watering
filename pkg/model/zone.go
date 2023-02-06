package model

import (
	"time"
)

type Zone struct {
	Id                         string
	Name                       string
	Mode                       *ZoneMode
	TargetMoisture             MoistureLevel
	MoistureSensors            []*MoistureSensor
	WaterOutlets               []*WaterOutlet
	Enabled                    bool
	AverageMoistureSensor      *ZoneAverageMoistureSensor
	WaterOutletsState          bool
	WaterOutletsStateChangedAt time.Time
}

func NewZone(id string, name string, sensors []*MoistureSensor, waterOutlets []*WaterOutlet) *Zone {
	zone := &Zone{
		Id:                         id,
		Name:                       name,
		Mode:                       GetDefaultZoneMode(),
		MoistureSensors:            sensors,
		WaterOutlets:               waterOutlets,
		Enabled:                    false,
		TargetMoisture:             MakeMoistureLevel(0),
		WaterOutletsState:          false,
		WaterOutletsStateChangedAt: time.Now(),
	}

	zone.AverageMoistureSensor = newZoneAverageMoistureSensor(zone)

	return zone
}

func (zone *Zone) SetWaterOutletsState(state bool) {

	changing := zone.WaterOutletsState != state

	zone.WaterOutletsState = state

	if changing {
		zone.WaterOutletsStateChangedAt = time.Now()
	}
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

func (zone Zone) MqttTargetMoistureCommandTopic(device *HassDevice) string {
	return zone.MqttTopic(device) + "/" + makeZoneHassConfiguration(zone, device).TargetMoistureCommandTopic
}

func (zone Zone) MqttModeCommandTopic(device *HassDevice) string {
	return zone.MqttTopic(device) + "/" + makeZoneHassConfiguration(zone, device).ModeCommandTopic
}

func (zone Zone) MqttStateTopic(device *HassDevice) string {
	return zone.MqttTopic(device) + "/" + makeZoneHassConfiguration(zone, device).StateTopic
}

func (zone Zone) AutoDiscoveryPayload(device *HassDevice) HassAutoDiscoverPayload {
	return makeZoneHassConfiguration(zone, device)
}
