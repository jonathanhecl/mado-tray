package main

import (
	"context"
	"log"
	"os"
	"sync"

	wailsruntime "github.com/wailsapp/wails/v2/pkg/runtime"
)

type App struct {
	ctx     context.Context
	store   *ConfigStore
	mu      sync.Mutex
	visible bool
}

func NewApp() (*App, error) {
	store, err := NewConfigStore()
	if err != nil {
		return nil, err
	}

	return &App{
		store:   store,
		visible: true,
	}, nil
}

func (a *App) startup(ctx context.Context) {
	a.mu.Lock()
	a.ctx = ctx
	a.mu.Unlock()

	initTray(a)

	scripts, err := a.store.LoadConfig()
	if err != nil {
		log.Printf("no se pudo cargar la configuración: %v", err)
		return
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
		wailsruntime.WindowHide(ctx)
	}
}

func (a *App) ShowWindow() {
	a.mu.Lock()
	ctx := a.ctx
	a.visible = true
	a.mu.Unlock()

	if ctx != nil {
		wailsruntime.WindowShow(ctx)
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

func (a *App) beforeClose(ctx context.Context) bool {
	status, err := GetStartupStatus()
	if err != nil || !status.Enabled {
		return false
	}

	a.mu.Lock()
	a.ctx = ctx
	a.visible = false
	a.mu.Unlock()

	wailsruntime.WindowHide(ctx)
	return true
}
