package homeassistant

import "github.com/mewejo/go-watering/world"

type HumidifierConfiguration struct {
	StateTopic               string        `json:"state_topic"`
	DeviceClass              string        `json:"device_class"`
	Name                     string        `json:"name"`
	ObjectId                 string        `json:"object_id"`
	UniqueId                 string        `json:"unique_id"`
	CommandTopic             string        `json:"command_topic"`
	TargetHumidityTopic      string        `json:"target_humidity_command_topic"`
	TargetHumidityStateTopic string        `json:"target_humidity_state_topic"`
	AvailabilityTopic        string        `json:"availability_topic"`
	Device                   DeviceDetails `json:"device"`
	PayloadAvailable         string        `json:"payload_available"`
	PayloadNotAvailable      string        `json:"payload_not_available"`
	PayloadOn                string        `json:"payload_on"`
	PayloadOff               string        `json:"payload_off"`
	Optimistic               bool          `json:"optimistic"`
	StateValueTemplate       string        `json:"state_value_template"`
	ModeCommandTopic         string        `json:"mode_command_topic"`
	ModeStateTopic           string        `json:"mode_state_topic"`
	Modes                    []string      `json:"modes"`
}

type MoistureSensorConfiguration struct {
	Name                string        `json:"name"`
	DeviceClass         string        `json:"device_class"`
	ObjectId            string        `json:"object_id"`
	UniqueId            string        `json:"unique_id"`
	StateTopic          string        `json:"state_topic"`
	StateValueTemplate  string        `json:"value_template"`
	AvailabilityTopic   string        `json:"availability_topic"`
	UnitOfMeasurement   string        `json:"unit_of_measurement"`
	Device              DeviceDetails `json:"device"`
	PayloadAvailable    string        `json:"payload_available"`
	PayloadNotAvailable string        `json:"payload_not_available"`
}

type DeviceDetails struct {
	Identifier   string `json:"identifiers"`
	Name         string `json:"name"`
	Model        string `json:"model"`
	Manufacturer string `json:"manufacturer"`
}

type ZoneState struct {
	MoistureLevel world.MoistureLevel `json:"moisture"`
	State         string              `json:"state"`
}

func NewDeviceDetails() DeviceDetails {
	d := DeviceDetails{}
	d.Manufacturer = "Josh Bonfield"
	d.Model = "Go Watering"
	d.Identifier = "vegtable-soaker"
	d.Name = "Vegtable Soaker"

	return d
}

func NewMoistureSensorConfiguration() MoistureSensorConfiguration {
	c := MoistureSensorConfiguration{}
	c.DeviceClass = "moisture"
	c.StateValueTemplate = "{{ value_json.moisture.percentage }}"
	c.UnitOfMeasurement = "%"
	c.PayloadAvailable = "online"
	c.PayloadNotAvailable = "offline"

	return c
}

func NewZoneHumidifierConfiguration() HumidifierConfiguration {
	c := HumidifierConfiguration{}
	c.DeviceClass = "humidifier"
	c.PayloadAvailable = "online"
	c.PayloadNotAvailable = "offline"
	c.PayloadOn = "on"
	c.PayloadOff = "off"
	c.StateValueTemplate = "{{ value_json.state }}"
	c.Modes = []string{"normal"}
	return c
}
