package app

import "fmt"

type App struct {
	//
}

func (app *App) Run() {
	fmt.Println("Bonsoir, Elliot")

	/*
		Make zone configurations
		Connect to MQTT
		Set LWT for availability topic (shared by all entities). Support modes normal/boost
		Publish MQTT auto discovery for Zones (climate), moisture sensors, and outlets not attached to zones (or zones with no sensors?)
		Find Arduino port
		Open serial connection
		Wait for heartbeat
		Wait for above to be ready
		Loop
			Read & process sensors and water states
			Publish zone states (freq as below)
			Publish sensor states (every 5 min in prod, 2 sec in testing)
			Publish outlets without zones states
			Check for heartbeat - if none after X:
				Close Arduino serial connection
				Attempt to establish a new connection
				Wait for a heartbeat
				Loop as above
		Program exit
			Send water off command
			Send unavailable topic to MQTT
			Exit


	*/

}

func Make() *App {
	return &App{}
}
