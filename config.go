package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

const (
	appConfigDirName  = "mado-tray"
	appConfigFileName = "config.json"
)

type Script struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Path     string `json:"path"`
	IsActive bool   `json:"is_active"`
}

type ConfigStore struct {
	mu   sync.Mutex
	path string
}

func NewConfigStore() (*ConfigStore, error) {
	path, err := defaultConfigPath()
	if err != nil {
		return nil, err
	}

	return &ConfigStore{path: path}, nil
}

func (s *ConfigStore) Path() string {
	return s.path
}

func (s *ConfigStore) LoadConfig() ([]Script, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.loadConfigLocked()
}

func (s *ConfigStore) SaveConfig(scripts []Script) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.saveConfigLocked(scripts)
}

func (s *ConfigStore) ToggleScript(id string, isActive bool) ([]Script, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	scripts, err := s.loadConfigLocked()
	if err != nil {
		return nil, err
	}

	for index := range scripts {
		if scripts[index].ID == id {
			scripts[index].IsActive = isActive
			return scripts, s.saveConfigLocked(scripts)
		}
	}

	return nil, fmt.Errorf("no existe un script con id %q", id)
}

func (s *ConfigStore) FindScript(id string) (Script, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	scripts, err := s.loadConfigLocked()
	if err != nil {
		return Script{}, err
	}

	for _, script := range scripts {
		if script.ID == id {
			return script, nil
		}
	}

	return Script{}, fmt.Errorf("no existe un script con id %q", id)
}

func (s *ConfigStore) loadConfigLocked() ([]Script, error) {
	if err := ensureConfigFile(s.path); err != nil {
		return nil, err
	}

	raw, err := os.ReadFile(s.path)
	if err != nil {
		return nil, fmt.Errorf("leyendo configuración: %w", err)
	}

	var scripts []Script
	if err := json.Unmarshal(raw, &scripts); err != nil {
		return nil, fmt.Errorf("el archivo %s no contiene una lista JSON válida: %w", s.path, err)
	}

	return scripts, nil
}

func (s *ConfigStore) saveConfigLocked(scripts []Script) error {
	raw, err := json.MarshalIndent(scripts, "", "  ")
	if err != nil {
		return fmt.Errorf("serializando configuración: %w", err)
	}

	if err := os.MkdirAll(filepath.Dir(s.path), 0o755); err != nil {
		return fmt.Errorf("creando carpeta de configuración: %w", err)
	}

	if err := os.WriteFile(s.path, append(raw, '\n'), 0o644); err != nil {
		return fmt.Errorf("guardando configuración: %w", err)
	}

	return nil
}

func ensureConfigFile(path string) error {
	if _, err := os.Stat(path); err == nil {
		return nil
	} else if !errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("revisando configuración: %w", err)
	}

	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return fmt.Errorf("creando carpeta de configuración: %w", err)
	}

	raw, err := json.MarshalIndent(defaultScripts(), "", "  ")
	if err != nil {
		return fmt.Errorf("creando configuración inicial: %w", err)
	}

	if err := os.WriteFile(path, append(raw, '\n'), 0o644); err != nil {
		return fmt.Errorf("guardando configuración inicial: %w", err)
	}

	return nil
}

func defaultConfigPath() (string, error) {
	configRoot, err := os.UserConfigDir()
	if err != nil {
		return "", fmt.Errorf("detectando carpeta de configuración: %w", err)
	}

	return filepath.Join(configRoot, appConfigDirName, appConfigFileName), nil
}

func defaultScripts() []Script {
	return []Script{
		{
			ID:       "example",
			Name:     "Script de ejemplo",
			Path:     "/ruta/a/tu/script.sh",
			IsActive: false,
		},
	}
}
