package persistence

import (
	"errors"
	"time"

	"github.com/mewejo/go-watering/pkg/model"
)

type moistureReadingStore struct {
	sensorId uint
	readings []model.MoistureReading
}

func (s *moistureReadingStore) getAverageSince(since time.Duration) (model.MoistureLevel, error) {
	readings := []model.MoistureLevel{}

	cutOffTime := time.Now().Add(-since)

	for _, reading := range s.readings {
		if reading.Time.Before(cutOffTime) {
			continue
		}

		readings = append(readings, reading.MoistureLevel)
	}

	if len(readings) < 1 {
		return model.MoistureLevel{}, errors.New("no readings to calculate average from")
	}

	var totalPercentage uint

	for _, reading := range readings {
		totalPercentage += reading.Percentage
	}

	return model.MakeMoistureLevel(
		uint(totalPercentage / uint(len(readings))),
	), nil
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

var moistureReadingStores []*moistureReadingStore

func GetLatestReadingForMoistureSensorId(sensorId uint) (*model.MoistureReading, error) {
	return getOrMakeStore(sensorId).getLatest()
}

func GetAverageReadingForSince(sensorId uint, since time.Duration) (model.MoistureLevel, error) {
	return getOrMakeStore(sensorId).getAverageSince(since)
}

func RecordMoistureReading(sensorId uint, reading model.MoistureReading) {
	getOrMakeStore(sensorId).recordReading(reading)
}

func getOrMakeStore(sensorId uint) *moistureReadingStore {
	for _, store := range moistureReadingStores {
		if store.sensorId == sensorId {
			return store
		}
	}

	store := moistureReadingStore{
		sensorId: sensorId,
	}

	moistureReadingStores = append(moistureReadingStores, &store)

	return &store
}
