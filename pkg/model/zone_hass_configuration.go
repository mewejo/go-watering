package model

import "github.com/mewejo/go-watering/pkg/constants"

type zoneHassConfiguration struct {
	Name                             string      `json:"name"`
	ObjectId                         string      `json:"object_id"`
	UniqueId                         string      `json:"unique_id"`
	StateTopic                       string      `json:"state_topic"`
	StateValueTemplate               string      `json:"state_value_template"`
	AvailabilityTopic                string      `json:"availability_topic"`
	HassDevice                       *HassDevice `json:"device"`
	PayloadAvailable                 string      `json:"payload_available"`
	PayloadNotAvailable              string      `json:"payload_not_available"`
	CommandTopic                     string      `json:"command_topic"`
	StateOn                          string      `json:"payload_on"`
	StateOff                         string      `json:"payload_off"`
	TargetMoistureTopic              string      `json:"target_humidity_command_topic"`
	TargetMoistureStateTopic         string      `json:"target_humidity_state_topic"`
	TargetMoistureStateValueTemplate string      `json:"target_humidity_state_template"`
	ModeStateTemplate                string      `json:"mode_state_template"`
	Modes                            []string    `json:"modes"`
	ModeCommandTopic                 string      `json:"mode_command_topic"`
	ModeStateTopic                   string      `json:"mode_state_topic"`
}

func (c zoneHassConfiguration) WithGlobalTopicPrefix(prefix string, device *HassDevice, entity HassAutoDiscoverable) HassAutoDiscoverPayload {
	c.AvailabilityTopic = prefix + "/" + c.HassDevice.GetFqAvailabilityTopic()
	c.CommandTopic = prefix + "/" + entity.MqttTopic(device) + "/" + c.CommandTopic
	c.ModeCommandTopic = prefix + "/" + entity.MqttTopic(device) + "/" + c.ModeCommandTopic
	c.ModeStateTopic = prefix + "/" + entity.MqttTopic(device) + "/" + c.ModeStateTopic
	c.StateTopic = prefix + "/" + entity.MqttTopic(device) + "/" + c.StateTopic
	c.TargetMoistureStateTopic = prefix + "/" + entity.MqttTopic(device) + "/" + c.TargetMoistureStateTopic
	return c
}

func makeZoneHassConfiguration(zone Zone, device *HassDevice) zoneHassConfiguration {
	c := zoneHassConfiguration{}
	c.Name = zone.Name
	c.StateOn = constants.HASS_STATE_ON
	c.StateOff = constants.HASS_STATE_OFF
	c.ObjectId = device.EntityPrefix + "zone-" + zone.Id
	c.UniqueId = c.ObjectId
	c.StateTopic = "state"
	c.TargetMoistureStateTopic = "state"
	c.TargetMoistureTopic = "target_moisture"
	c.AvailabilityTopic = device.GetFqAvailabilityTopic()
	c.ModeStateTemplate = "{{ value_json.mode.key }}"
	c.TargetMoistureStateValueTemplate = "{{ value_json.target_moisture.percentage }}"
	c.StateValueTemplate = "{% if value_json.enabled -%}" + c.StateOn + "{%- else -%}" + c.StateOff + "{%- endif %}"
	c.PayloadAvailable = device.PayloadAvailable
	c.PayloadNotAvailable = device.PayloadNotAvailable
	c.HassDevice = device
	c.CommandTopic = "command"
	c.ModeCommandTopic = "mode_command"
	c.ModeStateTopic = "mode"

	for _, mode := range zoneModes {
		c.Modes = append(c.Modes, mode.Key)
	}

	return c
}
