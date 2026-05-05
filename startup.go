package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

const loginItemName = "Mado-Tray"

type StartupStatus struct {
	Enabled   bool   `json:"enabled"`
	AppPath   string `json:"app_path"`
	Available bool   `json:"available"`
	Message   string `json:"message"`
}

func GetStartupStatus() (StartupStatus, error) {
	appPath, err := currentAppBundlePath()
	if err != nil {
		return StartupStatus{
			Enabled:   false,
			AppPath:   "",
			Available: false,
			Message:   err.Error(),
		}, nil
	}

	enabled, err := loginItemExists()
	if err != nil {
		return StartupStatus{}, err
	}

	return StartupStatus{
		Enabled:   enabled,
		AppPath:   appPath,
		Available: true,
		Message:   "Mado-Tray puede configurarse como ítem de inicio.",
	}, nil
}

func EnableStartup() (StartupStatus, error) {
	appPath, err := currentAppBundlePath()
	if err != nil {
		return StartupStatus{}, err
	}

	enabled, err := loginItemExists()
	if err != nil {
		return StartupStatus{}, err
	}
	if !enabled {
		script := fmt.Sprintf(`tell application "System Events"
	make login item at end with properties {name:"%s", path:"%s", hidden:true}
end tell`, escapeAppleScriptString(loginItemName), escapeAppleScriptString(appPath))
		if err := runAppleScript(script); err != nil {
			return StartupStatus{}, err
		}
	}

	return GetStartupStatus()
}

func DisableStartup() (StartupStatus, error) {
	enabled, err := loginItemExists()
	if err != nil {
		return StartupStatus{}, err
	}
	if enabled {
		script := fmt.Sprintf(`tell application "System Events"
	delete login item "%s"
end tell`, escapeAppleScriptString(loginItemName))
		if err := runAppleScript(script); err != nil {
			return StartupStatus{}, err
		}
	}

	return GetStartupStatus()
}

func loginItemExists() (bool, error) {
	script := fmt.Sprintf(`tell application "System Events"
	exists login item "%s"
end tell`, escapeAppleScriptString(loginItemName))

	output, err := runAppleScriptWithOutput(script)
	if err != nil {
		return false, err
	}

	return strings.TrimSpace(output) == "true", nil
}

func currentAppBundlePath() (string, error) {
	if runtime.GOOS != "darwin" {
		return "", fmt.Errorf("el arranque automático integrado solo está disponible en macOS")
	}

	executable, err := os.Executable()
	if err != nil {
		return "", fmt.Errorf("detectando ejecutable actual: %w", err)
	}

	resolved, err := filepath.EvalSymlinks(executable)
	if err != nil {
		resolved = executable
	}

	for current := filepath.Clean(resolved); current != string(filepath.Separator); current = filepath.Dir(current) {
		if strings.HasSuffix(current, ".app") {
			return current, nil
		}
	}

	return "", fmt.Errorf("para activar el arranque automático primero compila e instala Mado-Tray.app; en modo desarrollo no hay una .app final")
}

func runAppleScript(script string) error {
	_, err := runAppleScriptWithOutput(script)
	return err
}

func runAppleScriptWithOutput(script string) (string, error) {
	cmd := exec.Command("osascript", "-e", script)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("ejecutando AppleScript: %w: %s", err, strings.TrimSpace(string(output)))
	}

	return string(output), nil
}
