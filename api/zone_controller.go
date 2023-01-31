package api

import (
	"fmt"
	"net/http"

	"github.com/mewejo/go-watering/config"
	"github.com/mewejo/go-watering/world"
)

type zoneController struct {
	app *config.Application
}

func (c zoneController) handle(w http.ResponseWriter, r *http.Request) {
	c.app.Zones[0].TargetMoisture = world.MoistureLevel{
		Percentage: 0,
	}

	fmt.Fprintf(w, "Hi there, from controller I love %s!", r.URL.Path[1:])
}
