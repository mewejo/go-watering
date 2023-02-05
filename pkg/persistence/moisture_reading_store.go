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
	limitReadings(&s.readings, 100)
}

func limitReadings(s *[]model.MoistureReading, length int) {
	if len(*s) <= length {
		return
	}

	*s = (*s)[len(*s)-length:]
}

var stores []moistureReadingStore

func RecordReading(sensor model.MoistureSensor, reading model.MoistureReading) {
	getOrMakeStore(sensor).recordReading(reading)
}

func getOrMakeStore(sensor model.MoistureSensor) *moistureReadingStore {
	for _, store := range stores {
		if store.sensor == sensor {
			return &store
		}
	}

	store := moistureReadingStore{
		sensor: sensor,
	}

	stores = append(stores, store)

	return &store
}
