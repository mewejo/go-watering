package model

type moistureSensorHassConfiguration struct {
	Name                string      `json:"name"`
	DeviceClass         string      `json:"device_class"`
	ObjectId            string      `json:"object_id"`
	UniqueId            string      `json:"unique_id"`
	StateTopic          string      `json:"state_topic"`
	StateValueTemplate  string      `json:"value_template"`
	AvailabilityTopic   string      `json:"availability_topic"` // TODO this can be set on construct..?
	UnitOfMeasurement   string      `json:"unit_of_measurement"`
	HassDevice          *HassDevice `json:"device"`
	PayloadAvailable    string      `json:"payload_available"`
	PayloadNotAvailable string      `json:"payload_not_available"`
}

func makeMoistureSensorHassConfiguration(sensor MoistureSensor, device *HassDevice) moistureSensorHassConfiguration {
	c := moistureSensorHassConfiguration{}
	c.Name = sensor.Name
	c.ObjectId = device.EntityPrefix + "sensor-" + sensor.IdAsString()
	c.UniqueId = c.ObjectId
	c.StateTopic = "state"
	c.AvailabilityTopic = device.GetFqAvailabilityTopic()
	c.DeviceClass = "moisture"
	c.StateValueTemplate = "{{ value_json.moisture.percentage }}"
	c.UnitOfMeasurement = "%"
	c.PayloadAvailable = "online"
	c.PayloadNotAvailable = "offline"
	c.HassDevice = device

	return c
}
