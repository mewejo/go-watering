package persistence

import (
	"github.com/mewejo/go-watering/pkg/model"
)

type moistureReadingStore struct {
	sensorId uint
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

func RecordMoistureReading(sensorId uint, reading model.MoistureReading) {
	getOrMakeStore(sensorId).recordReading(reading)
}

func getOrMakeStore(sensorId uint) *moistureReadingStore {
	for _, store := range moistureReadingStores {
		if store.sensorId == sensorId {
			return &store
		}
	}

	store := moistureReadingStore{
		sensorId: sensorId,
	}

	moistureReadingStores = append(moistureReadingStores, store)

	return &store
}
