//go:build darwin

package main

/*
#cgo darwin CFLAGS: -x objective-c -fblocks
#cgo darwin LDFLAGS: -framework Cocoa -framework UniformTypeIdentifiers
#include <stdlib.h>
void MadoTrayCreate(void);
void MadoTrayShow(void);
void MadoTrayHide(void);
void MadoTraySetLocale(char* locale);
*/
import "C"

import (
	"sync"
	"unsafe"
)

var (
	trayOnce sync.Once
	trayApp  *App
)

func initTray(app *App) {
	trayOnce.Do(func() {
		trayApp = app
		C.MadoTrayCreate()
	})
	C.MadoTrayHide()
	C.MadoTrayShow()
}

func showTrayIcon() {
	C.MadoTrayShow()
}

func hideTrayIcon() {
}

func setTrayLocale(locale string) {
	cLocale := C.CString(locale)
	defer C.free(unsafe.Pointer(cLocale))
	C.MadoTraySetLocale(cLocale)
}

//export madoTrayShow
func madoTrayShow() {
	if trayApp != nil {
		trayApp.ShowWindow()
	}
}

//export madoTrayExit
func madoTrayExit() {
	if trayApp != nil {
		trayApp.Quit()
	}
}
