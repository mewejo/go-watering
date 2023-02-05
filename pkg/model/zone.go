package model

type Zone struct {
	Id              string
	Name            string
	Mode            *ZoneMode
	MoistureSensors []*MoistureSensor
	WaterOutlets    []*WaterOutlet
}

func NewZone(id string, name string, sensors []*MoistureSensor, waterOutlets []*WaterOutlet) *Zone {
	return &Zone{
		Id:              id,
		Name:            name,
		Mode:            getDefaultZoneMode(),
		MoistureSensors: sensors,
		WaterOutlets:    waterOutlets,
	}
}
