package homeassistant

type ZoneConfiguration struct {
	StateTopic          string        `json:"state_topic"`
	DeviceClass         string        `json:"device_class"`
	Name                string        `json:"name"`
	ObjectId            string        `json:"object_id"`
	UniqueId            string        `json:"unique_id"`
	CommandTopic        string        `json:"command_topic"`
	TargetHumidityTopic string        `json:"target_humidity_command_topic"`
	AvailabilityTopic   string        `json:"availability_topic"`
	Device              DeviceDetails `json:"device"`
	PayloadAvailable    string        `json:"payload_available"`
	PayloadNotAvailable string        `json:"payload_not_available"`
	Optimistic          bool          `json:"optimistic"`
}

type DeviceDetails struct {
	Identifier   string `json:"identifiers"`
	Name         string `json:"name"`
	Model        string `json:"model"`
	Manufacturer string `json:"manufacturer"`
}

func NewDeviceDetails() DeviceDetails {
	d := DeviceDetails{}
	d.Manufacturer = "Josh Bonfield"
	d.Model = "Go Watering"
	return d
}

func NewZoneConfiguration() ZoneConfiguration {
	c := ZoneConfiguration{}
	c.DeviceClass = "humidifier"
	c.PayloadAvailable = "online"
	c.PayloadNotAvailable = "offline"
	return c
}
