package api

import (
	"fmt"
	"net/http"

	"github.com/mewejo/go-watering/config"
)

type zoneController struct {
	app *config.Application
}

func (c zoneController) handle(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hi there, from controller I love %s!", r.URL.Path[1:])
}
