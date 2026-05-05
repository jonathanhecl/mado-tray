//go:build darwin

package main

/*
#cgo darwin CFLAGS: -x objective-c -fblocks
#cgo darwin LDFLAGS: -framework Cocoa -framework UniformTypeIdentifiers
void MadoTrayCreate(void);
*/
import "C"

import "sync"

var (
	trayOnce sync.Once
	trayApp  *App
)

func initTray(app *App) {
	trayOnce.Do(func() {
		trayApp = app
		C.MadoTrayCreate()
	})
}

//export madoTrayToggle
func madoTrayToggle() {
	if trayApp != nil {
		trayApp.ToggleWindow()
	}
}

//export madoTrayQuit
func madoTrayQuit() {
	if trayApp != nil {
		trayApp.Quit()
	}
}
