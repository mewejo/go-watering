package model

type ZoneHassState struct {
	Mode            *ZoneMode     `json:"mode"`
	AverageMoisture MoistureLevel `json:"average_moisture"`
}

func MakeZoneHassState(mode *ZoneMode, averageMoisture MoistureLevel) ZoneHassState {
	return ZoneHassState{
		Mode:            mode,
		AverageMoisture: averageMoisture,
	}
}
