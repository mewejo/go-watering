package model

type moistureSensorHassConfiguration struct {
	Name                string      `json:"name"`
	DeviceClass         string      `json:"device_class"`
	ObjectId            string      `json:"object_id"`
	UniqueId            string      `json:"unique_id"`
	StateTopic          string      `json:"state_topic"`
	StateValueTemplate  string      `json:"value_template"`
	AvailabilityTopic   string      `json:"availability_topic"`
	UnitOfMeasurement   string      `json:"unit_of_measurement"`
	HassDevice          *HassDevice `json:"device"`
	PayloadAvailable    string      `json:"payload_available"`
	PayloadNotAvailable string      `json:"payload_not_available"`
}

func (c moistureSensorHassConfiguration) WithGlobalTopicPrefix(prefix string, device *HassDevice, entity HassAutoDiscoverable) HassAutoDiscoverPayload {
	c.AvailabilityTopic = prefix + "/" + c.HassDevice.GetFqAvailabilityTopic()
	c.StateTopic = prefix + "/" + entity.MqttTopic(device) + "/" + c.StateTopic

	stateOverride := entity.OverriddenMqttStateTopic(device)

	if stateOverride != "" {
		c.StateTopic = prefix + "/" + stateOverride
	}

	return c
}

func makeMoistureSensorForZoneAverageMoistureSensor(averageMoistureSensor ZoneAverageMoistureSensor, device *HassDevice) moistureSensorHassConfiguration {
	c := makeMoistureSensorConfiugurationForDevice(device)
	c.Name = averageMoistureSensor.Name
	c.ObjectId = device.EntityPrefix + "sensor-" + averageMoistureSensor.Id
	c.StateValueTemplate = "{{ value_json.average_moisture.percentage }}"

	return c
}

func makeMoistureSensorConfiugurationForDevice(device *HassDevice) moistureSensorHassConfiguration {
	c := moistureSensorHassConfiguration{}
	c.UniqueId = c.ObjectId
	c.StateTopic = "state"
	c.AvailabilityTopic = device.GetFqAvailabilityTopic()
	c.DeviceClass = "moisture"
	c.UnitOfMeasurement = "%"
	c.PayloadAvailable = device.PayloadAvailable
	c.PayloadNotAvailable = device.PayloadNotAvailable
	c.HassDevice = device

	return c
}

func makeMoistureSensorHassConfiguration(sensor MoistureSensor, device *HassDevice) moistureSensorHassConfiguration {
	c := makeMoistureSensorConfiugurationForDevice(device)
	c.Name = sensor.Name
	c.ObjectId = device.EntityPrefix + "sensor-" + sensor.IdAsString()
	c.StateValueTemplate = "{{ value_json.moisture_level.percentage }}"

	return c
}
