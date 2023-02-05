package model

import "github.com/mewejo/go-watering/pkg/number"

type MoistureSensor struct {
	Id           uint
	Name         string
	DryThreshold uint
	WetThreshold uint
}

func (ms MoistureSensor) mapRawReadingToPercentage(raw uint) uint {
	return uint(number.ChangeRange(
		float64(raw),
		float64(ms.DryThreshold),
		float64(ms.WetThreshold),
		0,
		100,
	))
}

func MakeMoistureSensor(id uint, name string) MoistureSensor {
	return MoistureSensor{
		Id:           id,
		Name:         name,
		DryThreshold: 500,
		WetThreshold: 240,
	}
}