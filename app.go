package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"

	wailsruntime "github.com/wailsapp/wails/v2/pkg/runtime"
)

type App struct {
	ctx     context.Context
	store   *ConfigStore
	mu      sync.Mutex
	visible bool
	locale  string
}

func NewApp() (*App, error) {
	store, err := NewConfigStore()
	if err != nil {
		return nil, err
	}

	return &App{
		store:   store,
		visible: true,
		locale:  "en",
	}, nil
}

func (a *App) startup(ctx context.Context) {
	a.mu.Lock()
	a.ctx = ctx
	a.mu.Unlock()

	initTray(a)
	showTrayIcon()
	a.SetLocale(preferredLocale())

	scripts, err := a.store.LoadConfig()
	if err != nil {
		log.Printf("no se pudo cargar la configuración: %v", err)
		return
	}

	if err := CloseInactiveMadoTerminals(); err != nil {
		log.Printf("no se pudieron cerrar terminales inactivas de Mado-Tray: %v", err)
	}

	for _, script := range scripts {
		if !script.IsActive {
			continue
		}

		current := script
		go func() {
			if err := RunInVisibleTerminal(current.Path); err != nil {
				log.Printf("no se pudo ejecutar %s: %v", current.Name, err)
			}
		}()
	}
}

func (a *App) domReady(ctx context.Context) {
	a.mu.Lock()
	a.ctx = ctx
	visible := a.visible
	a.mu.Unlock()
	showTrayIcon()
	a.updateTrayLocale()
	if !visible {
		a.HideWindow()
	}
}

func preferredLocale() string {
	lang := strings.ToLower(os.Getenv("LANG"))
	if strings.HasPrefix(lang, "es") {
		return "es"
	}
	return "en"
}

func (a *App) GetScripts() ([]Script, error) {
	return a.store.LoadConfig()
}

func (a *App) ToggleScript(id string, isActive bool) ([]Script, error) {
	return a.store.ToggleScript(id, isActive)
}

func (a *App) AddScript(input ScriptInput) ([]Script, error) {
	return a.store.AddScript(input)
}

func (a *App) UpdateScript(id string, input ScriptInput) ([]Script, error) {
	return a.store.UpdateScript(id, input)
}

func (a *App) DeleteScript(id string) ([]Script, error) {
	return a.store.DeleteScript(id)
}

func (a *App) RunScript(id string) error {
	script, err := a.store.FindScript(id)
	if err != nil {
		return err
	}

	return RunInVisibleTerminal(script.Path)
}

func (a *App) GetStartupStatus() (StartupStatus, error) {
	return GetStartupStatus()
}

func (a *App) EnableStartup() (StartupStatus, error) {
	return EnableStartup()
}

func (a *App) DisableStartup() (StartupStatus, error) {
	return DisableStartup()
}

func (a *App) ToggleWindow() {
	a.mu.Lock()
	ctx := a.ctx
	visible := a.visible
	a.visible = !visible
	a.mu.Unlock()

	if ctx == nil {
		return
	}

	if visible {
		wailsruntime.WindowHide(ctx)
		return
	}

	wailsruntime.WindowShow(ctx)
}

func (a *App) HideWindow() {
	a.mu.Lock()
	ctx := a.ctx
	a.visible = false
	a.mu.Unlock()
	if ctx != nil {
		// Ocultamos la aplicación completa para mantener coherencia con
		// "Show All" / "Mostrar todas las ventanas" del menú de macOS.
		wailsruntime.Hide(ctx)
	}
}

func (a *App) ShowWindow() {
	a.mu.Lock()
	ctx := a.ctx
	a.visible = true
	a.mu.Unlock()

	if ctx != nil {
		// En macOS, tras ocultar la app puede quedar "hidden" a nivel aplicación.
		// runtime.Show la vuelve al frente y luego mostramos la ventana.
		wailsruntime.Show(ctx)
		wailsruntime.WindowUnminimise(ctx)
		wailsruntime.WindowShow(ctx)
		wailsruntime.WindowSetAlwaysOnTop(ctx, true)
		wailsruntime.WindowSetAlwaysOnTop(ctx, false)
	}
}

func (a *App) Quit() {
	a.mu.Lock()
	ctx := a.ctx
	a.mu.Unlock()

	if ctx != nil {
		wailsruntime.Quit(ctx)
		return
	}

	os.Exit(0)
}

func (a *App) SetLocale(locale string) {
	a.mu.Lock()
	if locale == "es" {
		a.locale = "es"
	} else {
		a.locale = "en"
	}
	a.mu.Unlock()

	a.updateTrayLocale()
}

func (a *App) PickScriptPath() (string, error) {
	a.mu.Lock()
	ctx := a.ctx
	locale := a.locale
	a.mu.Unlock()

	if ctx == nil {
		return "", fmt.Errorf("la aplicación todavía no está lista")
	}

	title := "Select script or executable"
	if locale == "es" {
		title = "Seleccionar script o ejecutable"
	}

	return wailsruntime.OpenFileDialog(ctx, wailsruntime.OpenDialogOptions{
		Title: title,
		Filters: []wailsruntime.FileFilter{
			{DisplayName: "Scripts (*.sh)", Pattern: "*.sh"},
			{DisplayName: "All files (*.*)", Pattern: "*.*"},
		},
	})
}

func (a *App) updateTrayLocale() {
	a.mu.Lock()
	locale := a.locale
	a.mu.Unlock()
	setTrayLocale(locale)
}
