package model

import (
	"errors"
	"log"
)

type ZoneMode struct {
	Name string
	Key  string
}

var zoneModes = []*ZoneMode{
	{
		Name: "Normal",
		Key:  "normal",
	},

	{
		Name: "Boost",
		Key:  "boost",
	},
}

func getDefaultZoneMode() *ZoneMode {
	mode, err := getZoneModeFromKey("normal")

	if err != nil {
		log.Fatal(err)
	}

	return mode
}

func getZoneModeFromKey(key string) (*ZoneMode, error) {
	for _, mode := range zoneModes {
		if mode.Key == key {
			return mode, nil
		}
	}

	return &ZoneMode{}, errors.New("invalid zone mode key")
}
