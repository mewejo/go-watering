package model

type ZoneHassState struct {
	Mode            *ZoneMode     `json:"mode"`
	AverageMoisture MoistureLevel `json:"average_moisture"`
	TargetMoisture  MoistureLevel `json:"target_moisture"`
	Enabled         bool          `json:"enabled"`
}

func MakeZoneHassState(zone *Zone, averageMoisture MoistureLevel) ZoneHassState {
	return ZoneHassState{
		Mode:            zone.Mode,
		AverageMoisture: averageMoisture,
		Enabled:         zone.Enabled,
		TargetMoisture:  zone.TargetMoisture,
	}
}
