package model

type MoistureLevel struct {
	Percentage uint `json:"percentage"`
}

func NewMoistureLevel(percentage uint) MoistureLevel {
	return MoistureLevel{
		Percentage: percentage,
	}
}
