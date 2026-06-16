package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"unicode"
)

const (
	appConfigDirName  = "mado-tray"
	appConfigFileName = "config.json"
)

type Script struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Path     string `json:"path"`
	Args     string `json:"args,omitempty"`
	IsActive bool   `json:"is_active"`
}

type ScriptInput struct {
	Name     string `json:"name"`
	Path     string `json:"path"`
	Args     string `json:"args"`
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

func (s *ConfigStore) AddScript(input ScriptInput) ([]Script, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	script, err := validateScriptInput(input)
	if err != nil {
		return nil, err
	}

	scripts, err := s.loadConfigLocked()
	if err != nil {
		return nil, err
	}

	script.ID = nextScriptID(scripts, script.Name)
	scripts = append(scripts, script)

	return scripts, s.saveConfigLocked(scripts)
}

func (s *ConfigStore) UpdateScript(id string, input ScriptInput) ([]Script, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	updated, err := validateScriptInput(input)
	if err != nil {
		return nil, err
	}

	scripts, err := s.loadConfigLocked()
	if err != nil {
		return nil, err
	}

	for index := range scripts {
		if scripts[index].ID == id {
			updated.ID = id
			scripts[index] = updated
			return scripts, s.saveConfigLocked(scripts)
		}
	}

	return nil, fmt.Errorf("no existe un script con id %q", id)
}

func (s *ConfigStore) DeleteScript(id string) ([]Script, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	scripts, err := s.loadConfigLocked()
	if err != nil {
		return nil, err
	}

	for index := range scripts {
		if scripts[index].ID == id {
			scripts = append(scripts[:index], scripts[index+1:]...)
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

	for i := range scripts {
		scripts[i] = normalizeScript(scripts[i])
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
			Args:     "",
			IsActive: false,
		},
	}
}

func ScriptCommand(path, args string) string {
	path = strings.TrimSpace(path)
	args = strings.TrimSpace(args)
	if args == "" {
		return path
	}
	return path + " " + args
}

func normalizeScript(script Script) Script {
	script.Name = strings.TrimSpace(script.Name)
	script.Path = strings.TrimSpace(script.Path)
	script.Args = strings.TrimSpace(script.Args)

	if script.Args != "" {
		return script
	}

	parts, err := splitShellWords(script.Path)
	if err != nil || len(parts) <= 1 {
		return script
	}

	script.Path = parts[0]
	script.Args = strings.Join(parts[1:], " ")
	return script
}

func validateScriptInput(input ScriptInput) (Script, error) {
	name := strings.TrimSpace(input.Name)
	if name == "" {
		return Script{}, fmt.Errorf("el nombre del proceso es obligatorio")
	}

	path := strings.TrimSpace(input.Path)
	if path == "" {
		return Script{}, fmt.Errorf("la ruta del proceso es obligatoria")
	}

	args := strings.TrimSpace(input.Args)

	return Script{
		Name:     name,
		Path:     path,
		Args:     args,
		IsActive: input.IsActive,
	}, nil
}

func nextScriptID(scripts []Script, name string) string {
	base := slugify(name)
	if base == "" {
		base = "process"
	}

	used := make(map[string]bool, len(scripts))
	for _, script := range scripts {
		used[script.ID] = true
	}

	if !used[base] {
		return base
	}

	for suffix := 2; ; suffix++ {
		candidate := fmt.Sprintf("%s-%d", base, suffix)
		if !used[candidate] {
			return candidate
		}
	}
}

func slugify(value string) string {
	var builder strings.Builder
	lastDash := false

	for _, char := range strings.ToLower(value) {
		switch {
		case unicode.IsLetter(char) || unicode.IsDigit(char):
			builder.WriteRune(char)
			lastDash = false
		case !lastDash:
			builder.WriteRune('-')
			lastDash = true
		}
	}

	return strings.Trim(builder.String(), "-")
}
