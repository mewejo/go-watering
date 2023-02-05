package model

type zoneHassConfiguration struct {
	Name                string      `json:"name"`
	ObjectId            string      `json:"object_id"`
	UniqueId            string      `json:"unique_id"`
	StateTopic          string      `json:"state_topic"`
	StateValueTemplate  string      `json:"value_template"`
	AvailabilityTopic   string      `json:"availability_topic"` // TODO this can be set on construct..?
	HassDevice          *HassDevice `json:"device"`
	PayloadAvailable    string      `json:"payload_available"`
	PayloadNotAvailable string      `json:"payload_not_available"`
	CommandTopic        string      `json:"command_topic"`
	StateOn             string      `json:"state_on"`
	StateOff            string      `json:"state_off"`
	TargetMoistureTopic string      `json:"target_humidity_command_topic"`
}

func makeZoneHassConfiguration(zone Zone, device *HassDevice) zoneHassConfiguration {
	c := zoneHassConfiguration{}
	c.Name = zone.Name
	c.ObjectId = device.EntityPrefix + "zone-" + zone.Id
	c.UniqueId = c.ObjectId
	c.StateTopic = "humidifier"
	c.TargetMoistureTopic = "target_moisture"
	c.AvailabilityTopic = device.GetFqAvailabilityTopic()
	c.StateValueTemplate = "{{ value_json.target }}"
	c.PayloadAvailable = "online"
	c.PayloadNotAvailable = "offline"
	c.HassDevice = device
	c.CommandTopic = "command"
	c.StateOn = "on"
	c.StateOff = "off"

	return c
}
