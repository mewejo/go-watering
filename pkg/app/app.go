package app

import "fmt"

type App struct {
	//
}

func (app *App) Run() {
	fmt.Println("Bonsoir, Elliot")
}

func Make() *App {
	return &App{}
}
