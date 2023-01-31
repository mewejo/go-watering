package api

import (
	"log"
	"net/http"

	"github.com/mewejo/go-watering/config"
)

func StartApi(app *config.Application) {

	zoneController := zoneController{
		app: app,
	}

	http.HandleFunc("/api/zones/", zoneController.handle)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
