package persistence

import (
	"errors"

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

func (s *moistureReadingStore) getLatest() (*model.MoistureReading, error) {
	if len(s.readings) < 1 {
		return &model.MoistureReading{}, errors.New("no readings available")
	}

	return &s.readings[len(s.readings)-1], nil
}

func limitReadings(s *[]model.MoistureReading, length int) {
	if len(*s) <= length {
		return
	}

	*s = (*s)[len(*s)-length:]
}

var moistureReadingStores []moistureReadingStore

func GetLatestReadingForMoistureSensorId(sensorId uint) (*model.MoistureReading, error) {
	return getOrMakeStore(sensorId).getLatest()
}

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
