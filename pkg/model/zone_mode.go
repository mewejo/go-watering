package model

import (
	"errors"
	"log"
	"time"
)

type ZoneMode struct {
	Key            string        `json:"key"`
	CutOffDuration time.Duration `json:"-"`
}

const normalMode = "normal"

var zoneModes = []*ZoneMode{
	{
		Key:            normalMode,
		CutOffDuration: time.Minute * 30,
	},

	{
		Key:            "15 Minute Boost",
		CutOffDuration: time.Second * 5,
	},

	{
		Key:            "30 Minute Boost",
		CutOffDuration: time.Second * 10,
	},

	{
		Key:            "1 hour Boost",
		CutOffDuration: time.Second * 15,
	},
}

func GetDefaultZoneMode() *ZoneMode {
	mode, err := GetZoneModeFromKey(normalMode)

	if err != nil {
		log.Fatal(err)
	}

	return mode
}

func GetZoneModeFromKey(key string) (*ZoneMode, error) {
	for _, mode := range zoneModes {
		if mode.Key == key {
			return mode, nil
		}
	}

	return &ZoneMode{}, errors.New("invalid zone mode key")
}
