package model

type WaterOutlet struct {
	Id          uint
	Name        string
	TargetState bool
	ActualState bool
}

func NewWaterOutlet(id uint, name string) *WaterOutlet {
	return &WaterOutlet{
		Id:   id,
		Name: name,
	}
}
