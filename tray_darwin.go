//go:build darwin

package main

/*
#cgo darwin CFLAGS: -x objective-c -fblocks
#cgo darwin LDFLAGS: -framework Cocoa -framework UniformTypeIdentifiers
void MadoTrayCreate(void);
void MadoTrayShow(void);
void MadoTrayHide(void);
*/
import "C"

import "sync"

var (
	trayOnce sync.Once
	trayApp  *App
)

func initTray(app *App, visible bool) {
	trayOnce.Do(func() {
		trayApp = app
		C.MadoTrayCreate()
	})

	if visible {
		C.MadoTrayShow()
		return
	}

	C.MadoTrayHide()
}

func showTrayIcon() {
	C.MadoTrayShow()
}

func hideTrayIcon() {
	C.MadoTrayHide()
}

//export madoTrayToggle
func madoTrayToggle() {
	if trayApp != nil {
		trayApp.ToggleWindow()
	}
}
