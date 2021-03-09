package main

import (
	"log"
	"os"

	"github.com/irth/go-simple"
	ui "github.com/irth/go-simple"
)

type ExitScene struct{}

type exitWidget struct{}

func (e exitWidget) Update(out ui.Output) ([]ui.BoundEventHandler, error) {
	// Update is called after empty draws its stuff and exits, so in our case
	// immediately after cleaning the screen
	log.Println("Goodbye...")
	os.Exit(0)
	return nil, nil
}

func (e exitWidget) Render() (string, error) {
	// Display empty label to clear the screen
	return ui.Label(ui.Pos(ui.Abs(0), ui.Abs(0), ui.Abs(10), ui.Abs(10)), " ").Render()

}

func (e ExitScene) Render() (simple.Widget, error) {
	return exitWidget{}, nil
}
