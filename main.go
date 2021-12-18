package main

import (
	_ "embed"
	"log"
	"os"

	"github.com/mattouille/proman/service/config"
	"github.com/mattouille/proman/service/database"
	"github.com/wailsapp/wails"
)

const (
	DefaultWidth  = 360
	DefaultHeight = 640
)

//go:embed frontend/public/build/bundle.js
var js string

//go:embed frontend/public/build/bundle.css
var css string

func main() {
	app := wails.CreateApp(&wails.AppConfig{
		Width:     DefaultWidth,
		MinWidth:  DefaultWidth,
		Height:    DefaultHeight,
		MinHeight: DefaultHeight,
		Resizable: true,
		Title:     "proman",
		JS:        js,
		CSS:       css,
		Colour:    "#131313",
	})

	// configuration automatically loads, but does need to be read in and checked for errors.
	err := config.ReadInConfig()
	if err != nil {
		log.Printf("Unable to load configuration: %s", err)

		os.Exit(1)
	}

	// start a new database connection
	err = database.New()
	if err != nil {
		log.Printf("Unable to start database service: %s", err)

		os.Exit(1)
	}

	app.Bind(NewConfig())
	app.Bind(NewValidator())
	app.Bind(NewProjects())
	app.Bind(NewEditorConfig())

	err = app.Run()
	if err != nil {
		log.Printf("Exited with errors: %s", err)

		os.Exit(1)
	}
}
