package config

import (
	"github.com/mewejo/go-watering/arduino"
	"github.com/mewejo/go-watering/world"
)

type Application struct {
	Zones []Zone
}

func GetApplication() *Application {
	a := Application{
		Zones: []Zone{
			{
				Id:   "zone-1",
				Name: "Zone 1",
				TargetMoisture: world.MoistureLevel{
					Percentage: 70,
				},
				MoistureSensors: []arduino.MoistureSensor{
					arduino.MOSTURE_SENSOR_1,
					arduino.MOSTURE_SENSOR_2,
				},
				WaterOutlets: []arduino.WaterOutlet{
					arduino.WATER_OUTLET_1,
				},
			},
			{
				Id:   "zone-2",
				Name: "Zone 2",
				TargetMoisture: world.MoistureLevel{
					Percentage: 70,
				},
				MoistureSensors: []arduino.MoistureSensor{
					arduino.MOSTURE_SENSOR_3,
					arduino.MOSTURE_SENSOR_4,
				},
				WaterOutlets: []arduino.WaterOutlet{
					arduino.WATER_OUTLET_2,
				},
			},
			{
				Id:   "zone-3",
				Name: "Zone 3",
				TargetMoisture: world.MoistureLevel{
					Percentage: 70,
				},
				MoistureSensors: []arduino.MoistureSensor{
					arduino.MOSTURE_SENSOR_5,
					arduino.MOSTURE_SENSOR_6,
				},
				WaterOutlets: []arduino.WaterOutlet{
					arduino.WATER_OUTLET_3,
				},
			},
		},
	}

	return &a
}
