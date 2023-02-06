package model

type waterOutletHassConfiguration struct {
	Name                string      `json:"name"`
	DeviceClass         string      `json:"device_class"`
	ObjectId            string      `json:"object_id"`
	UniqueId            string      `json:"unique_id"`
	StateTopic          string      `json:"state_topic"`
	StateValueTemplate  string      `json:"value_template"`
	AvailabilityTopic   string      `json:"availability_topic"`
	HassDevice          *HassDevice `json:"device"`
	PayloadAvailable    string      `json:"payload_available"`
	PayloadNotAvailable string      `json:"payload_not_available"`
	CommandTopic        string      `json:"command_topic"`
	StateOn             string      `json:"state_on"`
	StateOff            string      `json:"state_off"`
}

func (c waterOutletHassConfiguration) WithGlobalTopicPrefix(prefix string, device *HassDevice, entity HassAutoDiscoverable) HassAutoDiscoverPayload {
	c.AvailabilityTopic = prefix + "/" + c.HassDevice.GetFqAvailabilityTopic()
	c.CommandTopic = prefix + "/" + entity.MqttTopic(device) + "/" + c.CommandTopic
	return c
}

func makeWaterOutletHassConfiguration(outlet WaterOutlet, device *HassDevice) waterOutletHassConfiguration {
	c := waterOutletHassConfiguration{}
	c.Name = outlet.Name
	c.ObjectId = device.EntityPrefix + "outlet-" + outlet.IdAsString()
	c.UniqueId = c.ObjectId
	c.StateTopic = "state"
	c.AvailabilityTopic = device.GetFqAvailabilityTopic()
	c.DeviceClass = "switch"
	c.StateValueTemplate = "{{ value_json.actual ? 'on' : 'off' }}"
	c.PayloadAvailable = device.PayloadAvailable
	c.PayloadNotAvailable = device.PayloadNotAvailable
	c.HassDevice = device
	c.CommandTopic = "command"
	c.StateOn = "on"
	c.StateOff = "off"

	return c
}
