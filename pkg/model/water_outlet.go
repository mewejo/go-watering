package model

import (
	"errors"
	"strconv"
	"strings"
)

type WaterOutlet struct {
	Id                      uint
	Name                    string
	TargetState             bool
	ActualState             bool
	IndependentlyControlled bool
}

func NewWaterOutlet(id uint, name string, independentlyControlled bool) *WaterOutlet {
	return &WaterOutlet{
		Id:                      id,
		Name:                    name,
		IndependentlyControlled: independentlyControlled,
	}
}

func (wo WaterOutlet) IdAsString() string {
	return strconv.FormatUint(uint64(wo.Id), 10)
}

func (wo WaterOutlet) MqttTopic(device *HassDevice) string {
	return "switch/" + device.Namespace + "/outlet-" + wo.IdAsString()
}

func (wo WaterOutlet) AutoDiscoveryPayload(device *HassDevice) HassAutoDiscoverPayload {
	return makeWaterOutletHassConfiguration(wo, device)
}

func DecodeWaterOutletStateFromString(line string) (uint, bool, bool, error) {
	// WO:1:0:0 # WO:ID:REAL_STATE:SET_STATE (1 = on, 0 = off)
	parts := strings.Split(line, ":")

	if parts[0] != "WO" {
		return 0, false, false, errors.New("line was not a water outlet state")
	}

	outletId, err := strconv.Atoi(parts[1])

	if err != nil {
		return 0, false, false, err
	}

	actualState, err := strconv.Atoi(parts[2])

	if err != nil {
		return 0, false, false, err
	}

	targetState, err := strconv.Atoi(parts[3])

	if err != nil {
		return 0, false, false, err
	}

	return uint(outletId),
		actualState == 1,
		targetState == 1,
		nil
}
