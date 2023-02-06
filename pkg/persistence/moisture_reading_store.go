package persistence

import (
	"github.com/mewejo/go-watering/pkg/model"
)

type moistureReadingStore struct {
	sensor   model.MoistureSensor
	readings []model.MoistureReading
}

func (s *moistureReadingStore) recordReading(r model.MoistureReading) {
	s.readings = append(s.readings, r)
	limitReadings(&s.readings, 1000)
}

func limitReadings(s *[]model.MoistureReading, length int) {
	if len(*s) <= length {
		return
	}

	*s = (*s)[len(*s)-length:]
}

var moistureReadingStores []moistureReadingStore

func RecordMoistureReading(sensor model.MoistureSensor, reading model.MoistureReading) {
	getOrMakeStore(sensor).recordReading(reading)
}

func getOrMakeStore(sensor model.MoistureSensor) *moistureReadingStore {
	for _, store := range moistureReadingStores {
		if store.sensor.Id == sensor.Id {
			return &store
		}
	}

	store := moistureReadingStore{
		sensor: sensor,
	}

	moistureReadingStores = append(moistureReadingStores, store)

	return &store
}
