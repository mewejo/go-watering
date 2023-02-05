package model

type Zone struct {
	Id   uint
	Name string
	Mode *ZoneMode
}

func NewZone(id uint, name string) *Zone {
	return &Zone{
		Id:   id,
		Name: name,
		Mode: getDefaultZoneMode(),
	}
}
