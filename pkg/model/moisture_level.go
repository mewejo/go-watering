package model

type MoistureLevel struct {
	Percentage uint `json:"percentage"`
}

func MakeMoistureLevel(percentage uint) MoistureLevel {
	return MoistureLevel{
		Percentage: percentage,
	}
}
