package model

import (
	"errors"
	"log"
)

type ZoneMode struct {
	Name string `json:"name"`
	Key  string `json:"key"`
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

func GetDefaultZoneMode() *ZoneMode {
	mode, err := GetZoneModeFromKey("normal")

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
