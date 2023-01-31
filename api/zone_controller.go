package api

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/mewejo/go-watering/config"
	"github.com/mewejo/go-watering/world"
)

type zoneController struct {
	app *config.Application
}

func (c zoneController) handle(w http.ResponseWriter, r *http.Request) {
	percentage, err := strconv.Atoi(r.URL.Query().Get("target"))

	if err != nil {
		fmt.Fprintf(w, "Bad percentage!")
		return
	}

	zoneIndex, err := strconv.Atoi(r.URL.Query().Get("zone"))

	if err != nil {
		fmt.Fprintf(w, "Bad zone!")
		return
	}

	c.app.Zones[zoneIndex].TargetMoisture = world.MoistureLevel{
		Percentage: uint(percentage),
	}

	fmt.Fprintf(w, "Set zone %v to %v pc", zoneIndex, percentage)
}
