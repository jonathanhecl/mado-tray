package main

import (
	"embed"
	"log"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	app, err := NewApp()
	if err != nil {
		log.Fatal(err)
	}

	err = wails.Run(&options.App{
		Title:            "Mado-Tray",
		Width:            440,
		Height:           640,
		MinWidth:         360,
		MinHeight:        520,
		Frameless:        true,
		DisableResize:    true,
		Assets:           assets,
		BackgroundColour: &options.RGBA{R: 18, G: 22, B: 30, A: 1},
		OnStartup:        app.startup,
		OnBeforeClose:    app.beforeClose,
		Bind: []interface{}{
			app,
		},
	})
	if err != nil {
		log.Fatal(err)
	}
}
