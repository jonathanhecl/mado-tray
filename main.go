package main

import (
	"embed"
	"log"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/menu"
	"github.com/wailsapp/wails/v2/pkg/options"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	app, err := NewApp()
	if err != nil {
		log.Fatal(err)
	}

	startHidden := shouldStartHidden()
	app.visible = !startHidden

	err = wails.Run(&options.App{
		Title:         "Mado-Tray",
		Width:         440,
		Height:        640,
		MinWidth:      360,
		MinHeight:     520,
		Frameless:     true,
		DisableResize: true,
		// Dejamos que la app arranque visible y aplicamos el ocultado
		// en domReady para evitar cierres prematuros en modo dev.
		StartHidden:       false,
		HideWindowOnClose: true,
		Assets:            assets,
		BackgroundColour:  &options.RGBA{R: 18, G: 22, B: 30, A: 1},
		OnStartup:         app.startup,
		OnDomReady:        app.domReady,
		Menu:              buildMenu(app),
		SingleInstanceLock: &options.SingleInstanceLock{
			UniqueId: "com.jonathanhecl.mado-tray",
			OnSecondInstanceLaunch: func(_ options.SecondInstanceData) {
				app.ShowWindow()
			},
		},
		Bind: []interface{}{
			app,
		},
	})
	if err != nil {
		log.Fatal(err)
	}
}

func buildMenu(app *App) *menu.Menu {
	root := menu.NewMenu()
	appMenu := root.AddSubmenu("Mado-Tray")
	appMenu.AddText("Show", nil, func(_ *menu.CallbackData) {
		app.ShowWindow()
	})
	appMenu.AddSeparator()
	appMenu.AddText("Exit Mado-Tray", nil, func(_ *menu.CallbackData) {
		app.Quit()
	})
	root.Append(menu.EditMenu())
	return root
}

func shouldStartHidden() bool {
	status, err := GetStartupStatus()
	if err != nil {
		return false
	}

	return status.Enabled
}
